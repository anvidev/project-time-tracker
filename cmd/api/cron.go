package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/anvidev/project-time-tracker/internal/mailer"
	"github.com/anvidev/project-time-tracker/internal/store/users"
	"github.com/go-co-op/gocron/v2"
)

func initCronScheduler() (gocron.Scheduler, error) {
	tz, err := loadTimezone("Europe/Berlin")

	cronScheduler, err := gocron.NewScheduler(gocron.WithLocation(tz))
	if err != nil {
		return nil, err
	}

	return cronScheduler, nil
}

func (api *api) createCronJobs() error {
	if err := api.dailyJobAt(gocron.NewAtTime(06, 00, 00), api.notifyOnEmptyDay); err != nil {
		return err
	}
	return nil
}

func loadTimezone(name string) (*time.Location, error) {
	tz, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	return tz, nil
}

func (api *api) dailyJobAt(at gocron.AtTime, fn func()) error {
	if _, err := api.cron.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(at)),
		gocron.NewTask(fn)); err != nil {
		api.logger.Warn("failed to create cron job")
		return err
	}

	return nil

}

func (api *api) notifyOnEmptyDay() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	yesterday := time.Now().AddDate(0, 0, -1)
	if yesterday.Weekday() == 0 || yesterday.Weekday() == 6 {
		// yesterday was saturday or sunday
		return
	}

	isHoliday, err := isDanishHoliday(yesterday)
	if err != nil {
		api.logger.Warn("[CRON JOB] notifyOnEmptyDay - failed to fetch day info", "error", err)
		return
	}

	if isHoliday {
		return
	}

	userList, err := api.store.Users.List(ctx)
	if err != nil {
		api.logger.Warn("[CRON JOB] notifyOnEmptyDay - failed to fetch users", "error", err)
		return
	}

	userList = slices.DeleteFunc(userList, func(u users.User) bool {
		return !u.IsActive
	})

	for _, user := range userList {
		summary, err := api.store.TimeEntries.SummaryDay(ctx, user.Id, yesterday)
		if err != nil {
			api.logger.Warn(fmt.Sprintf("[CRON JOB] notifyOnEmptyDay - failed to get daily summary for %d", user.Id), "error", err)
			continue
		}

		if len(summary.TimeEntries) > 0 {
			continue
		}

		mailData := struct {
			User      users.User
			Yesterday string
		}{
			User:      user,
			Yesterday: yesterday.Format(time.DateOnly),
		}

		err = api.mails.Send([]string{user.Email}, "Ingen tidsregistreringer i går", mailer.NotifyEmptyDay, mailData)
		if err != nil {
			api.logger.Warn(fmt.Sprintf("[CRON JOB] notifyOnEmptyDay - failed to send email to %s", user.Email), "error", err)
		}
	}
}

func isDanishHoliday(date time.Time) (bool, error) {
	url := fmt.Sprintf("https://api.kalendarium.dk/Dayinfo/%s", date.Format("02-01-2006"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("api request to kalendarium failed with status %d: %s", res.StatusCode, string(body))
	}

	type DayInfo struct {
		Holiday bool `json:"holliday"`
		Events  []struct {
			Holiday bool `json:"holliday"`
		} `json:"events"`
	}

	var dayInfo DayInfo
	if err := json.Unmarshal(body, &dayInfo); err != nil {
		return false, err
	}

	isHoliday := dayInfo.Holiday
	if !isHoliday {
		for _, event := range dayInfo.Events {
			if event.Holiday {
				isHoliday = true
				break
			}
		}
	}

	return isHoliday, nil
}

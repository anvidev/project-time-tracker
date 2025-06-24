package main

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func (api *api) initCronScheduler() error {
	tz, err := loadTimezone("Europe/Berlin")

	cronScheduler, err := gocron.NewScheduler(gocron.WithLocation(tz))
	if err != nil {
		return err
	}

	api.cron = cronScheduler

	return err
}

func (api *api) createCronJobs() error {
	if err := api.dailyJobAt(gocron.NewAtTime(13, 40, 00), api.f); err != nil {
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

func (api *api) f() {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	api.logger.Info("hello")
}

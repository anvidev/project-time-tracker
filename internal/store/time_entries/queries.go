package time_entries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/types"
)

var (
	ErrTimeEntryNotDeleted = errors.New("time entry not deleted")
)

func (s *Store) Register(ctx context.Context, userId int64, input RegisterTimeEntryInput) (*TimeEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		insert into time_entries (
			category_id, user_id, date, duration, description
		)
		values (?, ?, ?, ?, ?)
		returning id
	`

	entry := TimeEntry{
		UserId:      userId,
		CategoryId:  input.CategoryId,
		Date:        input.Date,
		Duration:    input.Duration,
		Description: input.Description,
	}

	err := s.db.QueryRowContext(
		ctx,
		stmt,
		entry.CategoryId,
		entry.UserId,
		entry.Date,
		entry.Duration.String(),
		entry.Description,
	).Scan(
		&entry.Id,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *Store) Update(ctx context.Context, userId int64, input UpdateTimeEntryInput) (*TimeEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		update time_entries 
		set duration = ?
		where id = ? and user_id = ?
		returning id, category_id, user_id, date, duration, description
	`

	var entry TimeEntry
	err := s.db.QueryRowContext(
		ctx,
		stmt,
		input.Duration,
		input.EntryId,
		userId,
	).Scan(
		&entry.Id,
		&entry.CategoryId,
		&entry.UserId,
		&entry.Date,
		&entry.Duration,
		&entry.Description,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *Store) SummaryDay(ctx context.Context, userId int64, date time.Time) (*SummaryDay, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	summary, err := database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*SummaryDay, error) {
		weekdayHours, err := s.getWeekdayHours(ctx, tx, userId, date.Weekday())
		if err != nil {
			return nil, err
		}

		day, err := s.getDailySummary(ctx, tx, userId, date)
		if err != nil {
			return nil, err
		}

		day.MaxHours = *weekdayHours
		day.Weekday = strings.ToLower(date.Weekday().String())

		return day, nil
	})

	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (s *Store) SummaryMonth(ctx context.Context, userId int64, month time.Month, year int) (*SummaryMonth, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	type result struct {
		summary *SummaryDay
		err     error
	}

	results := make(chan result, daysInMonth)

	for i := 1; i <= daysInMonth; i++ {
		go func(i int) {
			date := time.Date(year, month, i, 0, 0, 0, 0, time.UTC)
			summaryDay, err := s.SummaryDay(ctx, userId, date)
			if err != nil {
				results <- result{err: fmt.Errorf("could not get summary for date %s: %w", date, err)}
				return
			}
			results <- result{summary: summaryDay}
		}(i)

	}

	summaryMonth := SummaryMonth{
		Month: strings.ToLower(month.String()),
	}

	for range daysInMonth {
		r := <-results
		if r.err != nil {
			cancel()
			return nil, r.err
		}
		summaryMonth.Days = append(summaryMonth.Days, *r.summary)
		summaryMonth.TotalHours.Duration += r.summary.TotalHours.Duration
		summaryMonth.MaxHours.Duration += r.summary.MaxHours.Duration
	}

	slices.SortFunc(summaryMonth.Days, func(a, b SummaryDay) int {
		return strings.Compare(a.Date, b.Date)
	})

	return &summaryMonth, nil
}

func (s *Store) Delete(ctx context.Context, id, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `delete from time_entries where id = ? and user_id = ?`

	result, err := s.db.ExecContext(ctx, stmt, id, userId)
	if err != nil {
		return err
	}

	affacted, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affacted != 1 {
		return ErrTimeEntryNotDeleted
	}

	return nil
}

func (s *Store) getWeekdayHours(ctx context.Context, tx *sql.Tx, userId int64, weekday time.Weekday) (*types.Duration, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `select hours from users_hours where user_id = ? and weekday = ?`

	var hours types.Duration

	if err := tx.QueryRowContext(ctx, stmt, userId, int(weekday)).Scan(&hours); err != nil {
		return nil, err
	}

	return &hours, nil
}

func (s *Store) getDailySummary(ctx context.Context, tx *sql.Tx, userId int64, date time.Time) (*SummaryDay, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		select 
			te.id,
			te.category_id,
			te.user_id,
			te.date,
			te.duration,
			te.description,
			(select title from categories where id = te.category_id) as category
		from time_entries te
		where user_id = ? and date = ?
		order by id desc
	`

	timeEntries := []TimeEntry{}
	dateString := date.Format(time.DateOnly)

	rows, err := tx.QueryContext(ctx, stmt, userId, dateString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e TimeEntry
		rows.Scan(
			&e.Id,
			&e.CategoryId,
			&e.UserId,
			&e.Date,
			&e.Duration,
			&e.Description,
			&e.Category,
		)
		timeEntries = append(timeEntries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var totalHours types.Duration

	for _, entry := range timeEntries {
		totalHours.Duration += entry.Duration.Duration
	}

	summary := &SummaryDay{
		Date:        dateString,
		TotalHours:  totalHours,
		TimeEntries: timeEntries,
	}

	return summary, nil
}

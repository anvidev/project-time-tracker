package live_timers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
	"github.com/anvidev/project-time-tracker/internal/types"
)

var (
	ErrTimerNotFound   = errors.New("no active timer")
	ErrTimerExists     = errors.New("a timer is already active")
	ErrTimerNotRunning = errors.New("timer is not running")
	ErrTimerNotPaused  = errors.New("timer is not paused")
	ErrNothingToSave   = errors.New("no time tracked")
	ErrInvalidDate     = errors.New("invalid date, expected yyyy-MM-dd")
)

type rowQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type timerRow struct {
	Timer
	accumulated time.Duration
	startedAt   time.Time
}

func (r *timerRow) elapsed(now time.Time) time.Duration {
	if r.Status == StatusRunning {
		return r.accumulated + now.Sub(r.startedAt)
	}
	return r.accumulated
}

func (r *timerRow) toTimer(now time.Time) *Timer {
	t := r.Timer
	t.Elapsed = types.Duration{Duration: r.elapsed(now)}
	t.ServerTime = now.Format(time.RFC3339)
	return &t
}

func get(ctx context.Context, q rowQuerier, userId int64) (*timerRow, error) {
	stmt := `
		select category_id, description, status, accumulated, started_at, date
		from live_timers
		where user_id = ?
	`

	var (
		row         timerRow
		accumulated int64
		startedAt   string
	)

	err := q.QueryRowContext(ctx, stmt, userId).Scan(
		&row.CategoryId,
		&row.Description,
		&row.Status,
		&accumulated,
		&startedAt,
		&row.Date,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTimerNotFound
		}
		return nil, err
	}

	row.accumulated = time.Duration(accumulated)
	row.startedAt, err = time.Parse(time.RFC3339Nano, startedAt)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

// Get returns the active timer of the user, or ErrTimerNotFound.
func (s *Store) Get(ctx context.Context, userId int64) (*Timer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	row, err := get(ctx, s.db, userId)
	if err != nil {
		return nil, err
	}

	return row.toTimer(time.Now()), nil
}

func (s *Store) Start(ctx context.Context, userId int64, input StartTimerInput) (*Timer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	now := time.Now()

	date := input.Date
	if date == "" {
		date = now.Format(time.DateOnly)
	} else if _, err := time.Parse(time.DateOnly, date); err != nil {
		return nil, ErrInvalidDate
	}

	stmt := `
		insert into live_timers (user_id, category_id, description, status, accumulated, started_at, date)
		values (?, ?, ?, ?, 0, ?, ?)
		on conflict (user_id) do nothing
	`

	result, err := s.db.ExecContext(
		ctx,
		stmt,
		userId,
		input.CategoryId,
		input.Description,
		StatusRunning,
		now.UTC().Format(time.RFC3339Nano),
		date,
	)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected != 1 {
		return nil, ErrTimerExists
	}

	return s.Get(ctx, userId)
}

func (s *Store) Pause(ctx context.Context, userId int64) (*Timer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	return database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*Timer, error) {
		row, err := get(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		if row.Status != StatusRunning {
			return nil, ErrTimerNotRunning
		}

		now := time.Now()
		row.accumulated = row.elapsed(now)
		row.Status = StatusPaused

		stmt := `update live_timers set status = ?, accumulated = ? where user_id = ?`
		if _, err := tx.ExecContext(ctx, stmt, row.Status, int64(row.accumulated), userId); err != nil {
			return nil, err
		}

		return row.toTimer(now), nil
	})
}

func (s *Store) Resume(ctx context.Context, userId int64) (*Timer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	return database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*Timer, error) {
		row, err := get(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		if row.Status != StatusPaused {
			return nil, ErrTimerNotPaused
		}

		now := time.Now()
		row.startedAt = now
		row.Status = StatusRunning

		stmt := `update live_timers set status = ?, started_at = ? where user_id = ?`
		if _, err := tx.ExecContext(ctx, stmt, row.Status, now.UTC().Format(time.RFC3339Nano), userId); err != nil {
			return nil, err
		}

		return row.toTimer(now), nil
	})
}

func (s *Store) Update(ctx context.Context, userId int64, input UpdateTimerInput) (*Timer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	return database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*Timer, error) {
		row, err := get(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		if input.CategoryId != nil {
			row.CategoryId = *input.CategoryId
		}
		if input.Description != nil {
			row.Description = *input.Description
		}

		stmt := `update live_timers set category_id = ?, description = ? where user_id = ?`
		if _, err := tx.ExecContext(ctx, stmt, row.CategoryId, row.Description, userId); err != nil {
			return nil, err
		}

		return row.toTimer(time.Now()), nil
	})
}

// Save turns the timer into a time entry and removes the timer.
func (s *Store) Save(ctx context.Context, userId int64, input SaveTimerInput) (*time_entries.TimeEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	return database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*time_entries.TimeEntry, error) {
		row, err := get(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		duration := row.elapsed(time.Now()).Round(time.Second)
		if input.Duration != nil {
			duration = input.Duration.Duration
		}
		if duration <= 0 {
			return nil, ErrNothingToSave
		}

		description := row.Description
		if input.Description != nil {
			description = *input.Description
		}

		entry := time_entries.TimeEntry{
			UserId:      userId,
			CategoryId:  row.CategoryId,
			Date:        row.Date,
			Duration:    types.Duration{Duration: duration},
			Description: description,
		}

		stmt := `
			insert into time_entries (category_id, user_id, date, duration, description)
			values (?, ?, ?, ?, ?)
			returning id
		`
		if err := tx.QueryRowContext(
			ctx,
			stmt,
			entry.CategoryId,
			entry.UserId,
			entry.Date,
			entry.Duration.String(),
			entry.Description,
		).Scan(&entry.Id); err != nil {
			return nil, err
		}

		if _, err := tx.ExecContext(ctx, `delete from live_timers where user_id = ?`, userId); err != nil {
			return nil, err
		}

		return &entry, nil
	})
}

// Discard removes the timer without saving it.
func (s *Store) Discard(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	result, err := s.db.ExecContext(ctx, `delete from live_timers where user_id = ?`, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrTimerNotFound
	}

	return nil
}

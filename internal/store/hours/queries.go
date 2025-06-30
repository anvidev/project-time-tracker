package hours

import (
	"context"
	"database/sql"

	"github.com/anvidev/project-time-tracker/internal/database"
)

func (s *Store) AllWeekdays(ctx context.Context, userId int64) ([]Weekday, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		select weekday, hours from users_hours where user_id = ?
	`

	rows, err := s.db.QueryContext(ctx, stmt, userId)
	if err != nil {
		return nil, err
	}

	days := make([]Weekday, 0, 0)
	for rows.Next() {
		var day Weekday
		if err = rows.Scan(
			&day.Weekday,
			&day.Hours,
		); err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}

func (s *Store) UpdateWeekdays(ctx context.Context, userId int64, data []Weekday) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	return database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		stmt := `
		update users_hours set hours = ? where user_id = ? and weekday = ?
		`

		for _, day := range data {
			if _, err := tx.ExecContext(ctx, stmt, day.Hours, userId, day.Weekday); err != nil {
				return err
			}
		}

		return nil
	})
}

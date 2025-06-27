package hours

import "context"

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

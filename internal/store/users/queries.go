package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/types"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrDublicateEmail     = errors.New("email is already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user is not active")
)

func (s *Store) Register(ctx context.Context, input RegisterUserInput) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	user, err := database.WithTxResult(ctx, s.db, func(tx *sql.Tx) (*User, error) {
		user, err := s.createUser(ctx, tx, input)
		if err != nil {
			return nil, err
		}

		defaultHours := Hours{
			UserId:    user.Id,
			Monday:    types.Duration{Duration: time.Duration(7*time.Hour + 30*time.Minute)},
			Tuesday:   types.Duration{Duration: time.Duration(7*time.Hour + 30*time.Minute)},
			Wednesday: types.Duration{Duration: time.Duration(7*time.Hour + 30*time.Minute)},
			Thursday:  types.Duration{Duration: time.Duration(7*time.Hour + 30*time.Minute)},
			Friday:    types.Duration{Duration: time.Duration(7 * time.Hour)},
			Saturday:  types.Duration{},
			Sunday:    types.Duration{},
		}

		if err := s.setHours(ctx, tx, &defaultHours); err != nil {
			return nil, err
		}

		return user, err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	var user User

	stmt := `
		select id, name, email, hash, is_active, role, created_at
		from users
		where email = ?
	`

	if err := s.db.
		QueryRowContext(ctx, stmt, email).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password.hash,
			&user.IsActive,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *Store) createUser(ctx context.Context, tx *sql.Tx, input RegisterUserInput) (*User, error) {
	user := User{
		Name:      input.Name,
		Email:     input.Email,
		Role:      RoleEmployee,
		IsActive:  true,
		CreatedAt: time.Now().Format(time.DateTime),
	}

	if err := user.Password.Set(input.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	stmt := `
		insert into users (name, email, hash, role, is_active, created_at)
		values (?, ?, ?, ?, ?, ?)
		returning rowid
	`

	if err := tx.QueryRowContext(
		ctx,
		stmt,
		user.Name,
		user.Email,
		user.Password.hash,
		user.Role,
		user.IsActive,
		user.CreatedAt,
	).Scan(&user.Id); err != nil {
		switch {
		case err.Error() == "failed to execute SQL:\nSQLite error: UNIQUE constraint failed: users.email":
			return nil, ErrDublicateEmail
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *Store) setHours(ctx context.Context, tx *sql.Tx, hours *Hours) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		insert or replace into users_hours (user_id, weekday, hours)
		values (?, ?, ?)
	`

	dayHours := map[int]types.Duration{
		1: hours.Monday,
		2: hours.Tuesday,
		3: hours.Wednesday,
		4: hours.Thursday,
		5: hours.Friday,
		6: hours.Saturday,
		0: hours.Sunday,
	}

	for dayOfWeek, duration := range dayHours {
		_, err := tx.ExecContext(ctx, stmt, hours.UserId, dayOfWeek, duration.Duration.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) List(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `select id, name, email, role, is_active, created_at from users`

	rows, err := s.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) GetById(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	var user User

	stmt := `
		select id, name, email, hash, is_active, role, created_at
		from users
		where id = ?
	`

	if err := s.db.
		QueryRowContext(ctx, stmt, id).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password.hash,
			&user.IsActive,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

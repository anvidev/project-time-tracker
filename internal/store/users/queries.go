package users

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

	if err := s.db.QueryRowContext(
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

func (s *Store) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	var user User

	stmt := `
		select rowid, name, email, hash, is_active, role, created_at
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

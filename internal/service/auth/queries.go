package auth

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrDublicateEmail  = errors.New("email is already in use")
)

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	user := User{
		Name:      req.Name,
		Email:     req.Email,
		Role:      RoleEmployee,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := user.Password.Set(req.Password); err != nil {
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
		case err.Error() == "unique something": // TODO: get exact error message from sqlite (libsql)
			return nil, ErrDublicateEmail
		default:
			return nil, err
		}
	}

	return &user, nil
}

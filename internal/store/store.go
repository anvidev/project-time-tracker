package store

import (
	"context"
	"database/sql"

	"github.com/anvidev/project-time-tracker/internal/store/sessions"
	"github.com/anvidev/project-time-tracker/internal/store/users"
)

type Store struct {
	Sessions SessionStorer
	Users    UserStorer
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Sessions: sessions.NewStore(db),
		Users:    users.NewStore(db),
	}
}

type SessionStorer interface {
	Create(ctx context.Context, userId int64) (*sessions.Session, error)
	Validate(ctx context.Context, token string) (*sessions.Session, error)
	Invalidate(ctx context.Context, token string) error
}

type UserStorer interface {
	Register(ctx context.Context, input users.RegisterUserInput) (*users.User, error)
	GetByEmail(ctx context.Context, email string) (*users.User, error)
}

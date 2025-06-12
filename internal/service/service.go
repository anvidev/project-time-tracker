package service

import (
	"context"
	"database/sql"

	"github.com/anvidev/project-time-tracker/internal/service/auth"
)

type Service struct {
	Auth AuthServicer
}

func NewService(db *sql.DB) *Service {
	return &Service{
		Auth: auth.NewService(db),
	}
}

type AuthServicer interface {
	Register(context.Context, auth.RegisterRequest) (*auth.User, error)
	// Login()
	// Validate()
	// CreateSession()
}

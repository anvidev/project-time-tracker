package service

import (
	"database/sql"

	"github.com/anvidev/project-time-tracker/internal/service/auth"
)

type Service struct {
	Auth auth.Service
}

func NewService(db *sql.DB) *Service {
	return &Service{
		Auth: auth.NewService(db),
	}
}

type AuthServicer interface {
	Register()
	Login()
	Validate()
	CreateSession()
}

package auth

import (
	"database/sql"
	"time"
)

type Service struct {
	db           *sql.DB
	queryTimeout time.Duration
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db:           db,
		queryTimeout: 5 * time.Second,
	}
}

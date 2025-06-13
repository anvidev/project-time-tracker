package sessions

import (
	"database/sql"
	"time"
)

type Store struct {
	db               *sql.DB
	queryTimeout     time.Duration
	sessionExpiresIn time.Duration
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:               db,
		queryTimeout:     5 * time.Second,
		sessionExpiresIn: 24 * time.Hour * 7,
	}
}

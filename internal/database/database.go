package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewContext(ctx context.Context, url, token string) (*sql.DB, error) {
	fullPath := fmt.Sprintf("%s?authToken=%s", url, token)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	db, err := sql.Open("libsql", fullPath)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

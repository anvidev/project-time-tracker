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

func WithTx(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func WithTxResult[T any](ctx context.Context, db *sql.DB, f func(tx *sql.Tx) (*T, error)) (*T, error) {
	var t T

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &t, err
	}

	pt, err := f(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &t, err
		}
		return &t, err
	}

	if err := tx.Commit(); err != nil {
		return &t, err
	}

	return pt, nil
}

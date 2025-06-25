package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/anvidev/project-time-tracker/internal/store/categories"
	"github.com/anvidev/project-time-tracker/internal/store/sessions"
	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
	"github.com/anvidev/project-time-tracker/internal/store/users"
)

type Store struct {
	TimeEntries TimeEntriesStorer
	Categories  CategoriesStorer
	Sessions    SessionStorer
	Users       UserStorer
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		TimeEntries: time_entries.NewStore(db),
		Categories:  categories.NewStore(db),
		Sessions:    sessions.NewStore(db),
		Users:       users.NewStore(db),
	}
}

type TimeEntriesStorer interface {
	Register(ctx context.Context, userId int64, input time_entries.RegisterTimeEntryInput) (*time_entries.TimeEntry, error)
	Update(ctx context.Context, userId, id int64, input time_entries.UpdateTimeEntryInput) (*time_entries.TimeEntry, error)
	Delete(ctx context.Context, id, userId int64) error
	SummaryDay(ctx context.Context, userId int64, date time.Time) (*time_entries.SummaryDay, error)
	SummaryMonth(ctx context.Context, userId int64, month time.Month, year int) (*time_entries.SummaryMonth, error)
}

type CategoriesStorer interface {
	Leafs(ctx context.Context, userId int64) ([]categories.Category, error)
	Follow(ctx context.Context, id, userId int64) error
	Unfollow(ctx context.Context, id, userId int64) error
	Tree(ctx context.Context, userId int64) ([]*categories.CategoryTree, error)
}

type SessionStorer interface {
	Create(ctx context.Context, userId int64) (*sessions.Session, error)
	Validate(ctx context.Context, token string) (*sessions.Session, error)
	InvalidateAll(ctx context.Context, userId int64) error
}

type UserStorer interface {
	Register(ctx context.Context, input users.RegisterUserInput) (*users.User, error)
	GetByEmail(ctx context.Context, email string) (*users.User, error)
	List(ctx context.Context) ([]users.User, error)
}

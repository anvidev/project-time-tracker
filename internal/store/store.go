package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/anvidev/project-time-tracker/internal/store/categories"
	"github.com/anvidev/project-time-tracker/internal/store/hours"
	"github.com/anvidev/project-time-tracker/internal/store/sessions"
	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
	"github.com/anvidev/project-time-tracker/internal/store/users"
)

type Store struct {
	TimeEntries TimeEntriesStorer
	Categories  CategoriesStorer
	Sessions    SessionStorer
	Users       UserStorer
	Hours       HourStorer
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		TimeEntries: time_entries.NewStore(db),
		Categories:  categories.NewStore(db),
		Sessions:    sessions.NewStore(db),
		Users:       users.NewStore(db),
		Hours:       hours.NewStore(db),
	}
}

type TimeEntriesStorer interface {
	Register(ctx context.Context, userId int64, input time_entries.RegisterTimeEntryInput) (*time_entries.TimeEntry, error)
	Update(ctx context.Context, userId, id int64, input time_entries.UpdateTimeEntryInput) (*time_entries.TimeEntry, error)
	Delete(ctx context.Context, id, userId int64) error
	SummaryDay(ctx context.Context, userId int64, date time.Time) (*time_entries.SummaryDay, error)
	SummaryMonth(ctx context.Context, userId int64, month time.Month, year int) (*time_entries.SummaryMonth, error)
	CategoryTotal(ctx context.Context, categoryId int64) (time.Duration, error)
	List(ctx context.Context, filters time_entries.Filters) ([]time_entries.TimeEntry, error)
}

type CategoriesStorer interface {
	Get(ctx context.Context, id int64) (*categories.Category, error)
	Create(ctx context.Context, input categories.CreateCategoryInput) (*categories.Category, error)
	Update(ctx context.Context, id int64, title string) (*categories.Category, error)
	ToggleRetire(ctx context.Context, id int64) error
	Leafs(ctx context.Context, userId int64) ([]categories.Category, error)
	Follow(ctx context.Context, id, userId int64) error
	Unfollow(ctx context.Context, id, userId int64) error
	Tree(ctx context.Context, userId int64) ([]*categories.CategoryTree, error)
	List(ctx context.Context) ([]categories.Category, error)
}

type SessionStorer interface {
	Create(ctx context.Context, userId int64) (*sessions.Session, error)
	Validate(ctx context.Context, token string) (*sessions.Session, error)
	InvalidateAll(ctx context.Context, userId int64) error
}

type UserStorer interface {
	Register(ctx context.Context, input users.RegisterUserInput) (*users.User, error)
	GetByEmail(ctx context.Context, email string) (*users.User, error)
	GetById(ctx context.Context, id int64) (*users.User, error)
	List(ctx context.Context) ([]users.User, error)
}

type HourStorer interface {
	AllWeekdays(ctx context.Context, userId int64) ([]hours.Weekday, error)
	UpdateWeekdays(ctx context.Context, userId int64, data []hours.Weekday) error
}

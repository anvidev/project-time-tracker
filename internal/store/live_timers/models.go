package live_timers

import (
	"github.com/anvidev/project-time-tracker/internal/types"
)

const (
	StatusRunning = "running"
	StatusPaused  = "paused"
)

// Timer is the live timer of a user. Elapsed is the time tracked at ServerTime;
// while Status is running, clients can keep counting from Elapsed + (now - ServerTime).
type Timer struct {
	CategoryId  int64          `json:"categoryId"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Date        string         `json:"date"` // yyyy-MM-dd (time.DateOnly)
	Elapsed     types.Duration `json:"elapsed"`
	ServerTime  string         `json:"serverTime"` // RFC3339
}

type StartTimerInput struct {
	CategoryId  int64  `json:"categoryId" validate:"required"`
	Description string `json:"description"`
	Date        string `json:"date"` // optional, defaults to today
}

type UpdateTimerInput struct {
	CategoryId  *int64  `json:"categoryId"`
	Description *string `json:"description"`
}

type SaveTimerInput struct {
	Duration    *types.Duration `json:"duration"`    // optional override of the tracked time
	Description *string         `json:"description"` // optional override
}

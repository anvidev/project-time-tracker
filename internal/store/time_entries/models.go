package time_entries

import (
	"github.com/anvidev/project-time-tracker/internal/types"
)

type TimeEntry struct {
	Id          int64          `json:"id"`
	CategoryId  int64          `json:"categoryId"`
	Category    string         `json:"category"`
	UserId      int64          `json:"userId"`
	Date        string         `json:"date"` // yyyy-MM-dd (time.DateOnly)
	Duration    types.Duration `json:"duration"`
	Description string         `json:"description"`
}

type SummaryDay struct {
	Date        string         `json:"date"`
	Weekday     string         `json:"weekday"`
	TotalHours  types.Duration `json:"totalHours"`
	MaxHours    types.Duration `json:"maxHours"`
	TimeEntries []TimeEntry    `json:"timeEntries"`
}

type SummaryMonth struct {
	Month      string         `json:"month"`
	TotalHours types.Duration `json:"totalHours"`
	MaxHours   types.Duration `json:"maxHours"`
	Days       []SummaryDay   `json:"days"`
}

type RegisterTimeEntryInput struct {
	CategoryId  int64          `json:"categoryId"`
	Date        string         `json:"date"`
	Duration    types.Duration `json:"duration"`
	Description string         `json:"description"`
}

package hours

import "github.com/anvidev/project-time-tracker/internal/types"

type Weekday struct {
	Weekday int64          `json:"weekday"`
	Hours   types.Duration `json:"hours"`
}

package hours

import (
	"github.com/anvidev/project-time-tracker/internal/types"
)

type Weekday struct {
	Weekday int64          `json:"weekday" apiduck:"desc=weekday is a number from 0-6 with Sunday at 0"`
	Hours   types.Duration `json:"hours"`
}

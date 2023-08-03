// Package structs contains structs for API responses
package structs

import (
	"fmt"
	"time"
)

// Health contains data about the status of the API and possible failures
type Health struct {
	HasFailures   bool `json:"failures"`
	CurrentStatus struct {
		UpdatedAt time.Time `json:"updatedInterval"`
		Fail      struct {
			System []string
			Beard  []string
		}
	}
}

func (h Health) String() string {
	var detail string
	if h.HasFailures {
		detail = fmt.Sprintf(
			"system=%v,beard=%v,ts=%s",
			h.CurrentStatus.Fail.System,
			h.CurrentStatus.Fail.Beard,
			h.CurrentStatus.UpdatedAt.Format(time.RFC3339),
		)
	} else {
		detail = "ok"
	}
	return fmt.Sprintf("Health{%s}", detail)
}

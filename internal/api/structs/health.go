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
	return fmt.Sprintf("Health{HasFailures=%t}", h.HasFailures)
}

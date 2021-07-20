package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

// ServiceLevel defines a Service Level
type ServiceLevel struct {
	Percentage float64
	Duration   time.Duration
}

// MarshalJSON marshals this into JSON
func (serviceLevel ServiceLevel) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Percentage float64 `json:"percentage"`
		Duration   int64   `json:"durationMs"`
	}{
		Percentage: serviceLevel.Percentage,
		Duration:   serviceLevel.Duration.Milliseconds(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (serviceLevel *ServiceLevel) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		Percentage float64 `json:"percentage"`
		Duration   int64   `json:"durationMs"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	serviceLevel.Percentage = inner.Percentage
	serviceLevel.Duration = time.Duration(inner.Duration) * time.Millisecond
	return
}

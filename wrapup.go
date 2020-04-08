package purecloud

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

// Wrapup describes a Wrapup
type Wrapup struct {
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	Notes       string        `json:"notes"`
	Tags        []string      `json:"tags"`
	Duration    time.Duration `json:"-"`
	EndTime     time.Time     `json:"endTime"`
	Provisional bool          `json:"provisional"`
}

// MarshalJSON marshals this into JSON
func (wrapup Wrapup) MarshalJSON() ([]byte, error) {
	type surrogate Wrapup
	data, err := json.Marshal(struct {
		surrogate
		DurationSeconds int64 `json:"durationSeconds"`
	}{
		surrogate:       surrogate(wrapup),
		DurationSeconds: int64(wrapup.Duration.Seconds()),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (wrapup *Wrapup) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Wrapup
	var inner struct {
		surrogate
		DurationSeconds int64 `json:"durationSeconds"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*wrapup = Wrapup(inner.surrogate)
	wrapup.Duration = time.Duration(inner.DurationSeconds) * time.Second
	return
}

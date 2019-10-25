package purecloud

import (
	"encoding/json"
  "time"

	"github.com/pkg/errors"
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
	return json.Marshal(struct {
    surrogate
		D          int64   `json:"durationSeconds"`
	}{
    surrogate: surrogate(wrapup),
		D:         int64(wrapup.Duration.Seconds()),
	})
}

// UnmarshalJSON unmarshals JSON into this
func (wrapup *Wrapup) UnmarshalJSON(payload []byte) (err error) {
  type surrogate Wrapup
	var inner struct {
    surrogate
		D          int64   `json:"durationSeconds"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
  }
  *wrapup = Wrapup(inner.surrogate)
	wrapup.Duration = time.Duration(inner.D) * time.Second
	return
}
package purecloud

import (
	"time"
)

// Wrapup describes a Wrapup
type Wrapup struct {
  Name        string        `json:"name"`
  Code        string        `json:"code"`
  Notes       string        `json:"notes"`
  Tags        []string      `json:"tags"`
  Duration    time.Duration `json:"durationSeconds"` // time.Duration
  EndTime     time.Time     `json:"endTime"`
  Provisional bool          `json:"provisional"`
}
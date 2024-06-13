package gcloudcx

import (
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

type OpenMessageFrom struct {
	ID        string `json:"id"`
	Type      string `json:"idType,omitempty"`
	Firstname string `json:"firstName,omitempty"`
	Lastname  string `json:"lastName,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (from OpenMessageFrom) Redact() interface{} {
	redacted := from
	if len(from.Firstname) > 0 {
		redacted.Firstname = logger.RedactWithHash(from.Firstname)
	}
	if len(from.Lastname) > 0 {
		redacted.Lastname = logger.RedactWithHash(from.Lastname)
	}
	if len(from.Nickname) > 0 {
		redacted.Nickname = logger.RedactWithHash(from.Nickname)
	}
	return redacted
}

// Validate checks if the object is valid
func (from *OpenMessageFrom) Validate() (err error) {
	if len(from.ID) == 0 {
		return errors.ArgumentMissing.With("ID")
	}
	return
}

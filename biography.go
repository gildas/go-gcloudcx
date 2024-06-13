package gcloudcx

import "github.com/gildas/go-logger"

// Biography describes a User's biography
type Biography struct {
	Biography string   `json:"biography"`
	Interests []string `json:"interests"`
	Hobbies   []string `json:"hobbies"`
	Spouse    string   `json:"spouse"`
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (biography Biography) Redact() interface{} {
	redacted := biography
	if len(biography.Biography) > 0 {
		redacted.Biography = logger.RedactWithHash(biography.Biography)
	}
	if len(biography.Spouse) > 0 {
		redacted.Spouse = logger.RedactWithHash(biography.Spouse)
	}
	return redacted
}

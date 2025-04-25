package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type Story struct {
	Type      string    `json:"type"` // Mention, Story
	ReplyToID uuid.UUID `json:"replyToId"`
	URL       *url.URL  `json:"url"`
}

// UnmarshalJSON unmarshals a Story from JSON
//
// Implements json.Unmarshaler
func (story *Story) UnmarshalJSON(data []byte) error {
	type surrogate Story
	var inner struct {
		surrogate
		ReplyToID core.UUID `json:"replyToId"`
		URL       *core.URL `json:"url"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*story = Story(inner.surrogate)
	story.ReplyToID = uuid.UUID(inner.ReplyToID)
	story.URL = (*url.URL)(inner.URL)

	return nil
}

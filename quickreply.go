package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type QuickReply struct {
	Text       string   `json:"text"`
	Payload    string   `json:"payload"`
	Action     string   `json:"action"`
	URL        *url.URL `json:"url,omitempty"`
	IsSelected bool     `json:"isSelected"`
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (quickReply *QuickReply) UnmarshalJSON(payload []byte) (err error) {
	type surrogate QuickReply
	var inner struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*quickReply = QuickReply(inner.surrogate)
	quickReply.URL = (*url.URL)(inner.URL)

	return nil
}

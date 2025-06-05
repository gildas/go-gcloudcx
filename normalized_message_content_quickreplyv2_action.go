package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageQuickReplyV2Action describes the action of a QuickReplyV2
type NormalizedMessageQuickReplyV2Action struct {
	Action   string   `json:"action,omitempty"` // "Message"
	Text     string   `json:"text"`
	Payload  string   `json:"payload"`
	ImageURL *url.URL `json:"image,omitempty"`
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (action NormalizedMessageQuickReplyV2Action) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageQuickReplyV2Action

	data, err := json.Marshal(struct {
		surrogate
		ImageURL *core.URL `json:"image,omitempty"`
	}{
		surrogate: surrogate(action),
		ImageURL:  (*core.URL)(action.ImageURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (action *NormalizedMessageQuickReplyV2Action) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageQuickReplyV2Action

	var inner struct {
		surrogate
		ImageURL *core.URL `json:"image"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = NormalizedMessageQuickReplyV2Action(inner.surrogate)
	action.ImageURL = (*url.URL)(inner.ImageURL)
	return
}

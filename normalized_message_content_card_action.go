package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageCardAction describes the action of a Card
type NormalizedMessageCardAction struct {
	ActionType string   `json:"type,omitempty"` // "Link", "Postback", "Unknown"
	Text       string   `json:"text"`
	Payload    string   `json:"payload"`
	URL        *url.URL `json:"url,omitempty"`
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (action NormalizedMessageCardAction) GetType() string {
	return action.ActionType
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (action NormalizedMessageCardAction) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCardAction

	data, err := json.Marshal(struct {
		surrogate
		URL *core.URL `json:"url"`
	}{
		surrogate: surrogate(action),
		URL:       (*core.URL)(action.URL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (action *NormalizedMessageCardAction) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCardAction

	var inner struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = NormalizedMessageCardAction(inner.surrogate)
	action.URL = (*url.URL)(inner.URL)
	return
}

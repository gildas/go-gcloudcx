package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// OpenMessageCardAction describes the action of a Card
type OpenMessageCardAction struct {
	ActionType string   `json:"type,omitempty"` // "Link", "Postback", "Unknown"
	Text       string   `json:"text"`
	Payload    string   `json:"payload"`
	URL        *url.URL `json:"url,omitempty"`
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (action OpenMessageCardAction) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageCardAction

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
func (action *OpenMessageCardAction) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageCardAction

	var inner struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = OpenMessageCardAction(inner.surrogate)
	action.URL = (*url.URL)(inner.URL)
	return
}

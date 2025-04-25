package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type Card struct {
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	URL           *url.URL     `json:"url"`
	DefaultAction CardAction   `json:"defaultAction"`
	Actions       []CardAction `json:"actions"`
}

type CardAction struct {
	Type       string   `json:"type"` // Link, Postback
	Text       string   `json:"text"`
	Payload    string   `json:"payload"`
	URL        *url.URL `json:"url,omitempty"`
	IsSelected bool     `json:"isSelected"`
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (card *Card) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Card
	var inner struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*card = Card(inner.surrogate)
	card.URL = (*url.URL)(inner.URL)

	return nil
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (action *CardAction) UnmarshalJSON(payload []byte) (err error) {
	type surrogate CardAction
	var inner struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = CardAction(inner.surrogate)
	action.URL = (*url.URL)(inner.URL)

	return nil
}

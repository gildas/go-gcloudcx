package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type Sticker struct {
	ID  string   `json:"id"`
	URL *url.URL `json:"url"`
}

// UnmarshalJSON unmarshals a Sticker from JSON
//
// Implements json.Unmarshaler
func (sticker *Sticker) UnmarshalJSON(data []byte) error {
	type surrogate Sticker
	var inner struct {
		surrogate
		ID  string    `json:"id"`
		URL *core.URL `json:"url"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*sticker = Sticker(inner.surrogate)
	sticker.ID = inner.ID
	sticker.URL = (*url.URL)(inner.URL)

	return nil
}

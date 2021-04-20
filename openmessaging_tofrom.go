package purecloud

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageTo struct {
	ID string `json:"id"`
}

type OpenMessageFrom struct {
	ID        string   `json:"id"`
	Type      string   `json:"idType"`
	Firstname string   `json:"firstName"`
	Lastname  string   `json:"lastName"`
	Nickname  string   `json:"nickname"`
	ImageURL  *url.URL `json:"-"`
}

// MarshalJSON marshals this into JSON
func (from OpenMessageFrom) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageFrom
	data, err := json.Marshal(struct{
		surrogate
		I *core.URL `json:"image"`
	}{
		surrogate: surrogate(from),
		I:         (*core.URL)(from.ImageURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (from *OpenMessageFrom) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageFrom
	var inner struct {
		surrogate
		I *core.URL `json:"image"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*from = OpenMessageFrom(inner.surrogate)
	from.ImageURL = (*url.URL)(inner.I)
	return
}
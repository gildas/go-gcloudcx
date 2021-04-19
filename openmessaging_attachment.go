package purecloud

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageAttachment struct {
	ID        string   `json:"id"`
	Type      string   `json:"mediaType"`
	URL       *url.URL `json:"-"`
	Mime      string   `json:"mime,omitempty"`
	Filename  string   `json:"filename,omitempty"`
	Text      string   `json:"text,omitempty"`
	Hash      string   `json:"sha256,omitempty"`
}

// MarshalJSON marshals this into JSON
func (attachment OpenMessageAttachment) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageAttachment
	data, err := json.Marshal(struct {
		surrogate
		U *core.URL `json:"url"`
	}{
		surrogate: surrogate(attachment),
		U:         (*core.URL)(attachment.URL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (attachment *OpenMessageAttachment) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageAttachment
	var inner struct {
		surrogate
		U *core.URL `json:"url"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*attachment = OpenMessageAttachment(inner.surrogate)
	attachment.URL = (*url.URL)(inner.U)
	return
}
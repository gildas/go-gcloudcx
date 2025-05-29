package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// HelpLink represents a link to help documentation or resources
type HelpLink struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	URI         *url.URL `json:"uri"`
}

// MarshalJSON marshals this into JSON
func (helpLink HelpLink) MarshalJSON() ([]byte, error) {
	type surrogate HelpLink
	data, err := json.Marshal(struct {
		surrogate
		URI *core.URL `json:"uri"`
	}{
		surrogate: surrogate(helpLink),
		URI:       (*core.URL)(helpLink.URI),
	})
	if err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}
	return data, nil
}

// UnmarshalJSON unmarshals JSON into this
func (helpLink *HelpLink) UnmarshalJSON(payload []byte) error {
	type surrogate HelpLink
	var inner struct {
		surrogate
		URI *core.URL `json:"uri"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*helpLink = HelpLink(inner.surrogate)
	helpLink.URI = (*url.URL)(inner.URI)
	return nil
}

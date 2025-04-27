package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type Media struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	MediaType     string    `json:"mediaType"`
	ContentLength uint64    `json:"contentLength"`
	URL           *url.URL  `json:"url"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (media Media) GetID() uuid.UUID {
	return media.ID
}

// Validate validates the media
func (media *Media) Validate() error {
	var merr errors.MultiError

	if media.ID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("id"))
	}
	if media.URL == nil {
		merr.Append(errors.ArgumentMissing.With("url"))
	}

	return merr.AsError()
}

// UnmarshalJSON unmarshals the media from JSON
//
// Implements json.Unmarshaler
func (media *Media) UnmarshalJSON(data []byte) error {
	type surrogate Media
	var inner struct {
		surrogate
		ID  core.UUID `json:"id"`
		URL *core.URL `json:"url"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*media = Media(inner.surrogate)
	media.ID = uuid.UUID(inner.ID)
	media.URL = (*url.URL)(inner.URL)

	return errors.JSONUnmarshalError.Wrap(media.Validate())
}

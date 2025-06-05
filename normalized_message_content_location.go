package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type NormalizedMessageLocationContent struct {
	Text      string   `json:"text"`
	URL       *url.URL `json:"url"`
	Address   string   `json:"address"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageLocationContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (location NormalizedMessageLocationContent) GetType() string {
	return "Location"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (location NormalizedMessageLocationContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageLocationContent
	type Location struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	data, err := json.Marshal(struct {
		ContentType string   `json:"contentType"`
		Location    Location `json:"location"`
	}{
		ContentType: location.GetType(),
		Location: Location{
			surrogate: surrogate(location),
			URL:       (*core.URL)(location.URL),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (location *NormalizedMessageLocationContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageLocationContent
	type Location struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	var inner struct {
		Location Location `json:"location"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*location = NormalizedMessageLocationContent(inner.Location.surrogate)
	location.URL = (*url.URL)(inner.Location.URL)
	return nil
}

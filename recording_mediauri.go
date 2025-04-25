package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type RecordingMediaURI struct {
	URI  *url.URL  `json:"mediaUri"`
	Data []float64 `json:"waveformData,omitempty"`
}

// UnmarshalJSON unmarshals the media URI from JSON
func (mediaURI *RecordingMediaURI) UnmarshalJSON(data []byte) error {
	type surrogate RecordingMediaURI
	var inner struct {
		surrogate
		URI *core.URL `json:"mediaUri"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*mediaURI = RecordingMediaURI(inner.surrogate)
	mediaURI.URI, _ = url.Parse(inner.URI.String())
	return nil
}

package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// Image represents an image used in various GCloudCX resources.
type Image struct {
	ImageURI   *url.URL `json:"imageUri"`
	Resolution string   `json:"resolution"`
}

// MarshalJSON marshals this into JSON
func (image Image) MarshalJSON() ([]byte, error) {
	type surrogate Image
	data, err := json.Marshal(struct {
		surrogate
		ImageURI *core.URL `json:"imageUri"`
	}{
		surrogate: surrogate(image),
		ImageURI:  (*core.URL)(image.ImageURI),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (image *Image) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Image
	var inner struct {
		surrogate
		ImageURI *core.URL `json:"imageUri"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*image = Image(inner.surrogate)
	image.ImageURI = (*url.URL)(inner.ImageURI)
	return
}

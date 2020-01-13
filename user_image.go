package purecloud

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// UserImage represents a User's Avatar image
type UserImage struct {
	ImageURL   *url.URL `json:"-"`
	Resolution string   `json:"resolution"`
}

// MarshalJSON marshals this into JSON
func (userImage UserImage) MarshalJSON() ([]byte, error) {
	type surrogate UserImage
	data, err := json.Marshal(struct {
		surrogate
		I *core.URL `json:"imageUrl"`
	}{
		surrogate: surrogate(userImage),
		I:         (*core.URL)(userImage.ImageURL),
	})
	if err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}
	return data, nil
}

// UnmarshalJSON unmarshals JSON into this
func (userImage *UserImage) UnmarshalJSON(payload []byte) (err error) {
	type surrogate UserImage
	var inner struct {
		surrogate
		I *core.URL `json:"imageUrl"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*userImage = UserImage(inner.surrogate)
	userImage.ImageURL = (*url.URL)(inner.I)
	return
}
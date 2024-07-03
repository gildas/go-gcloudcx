package gcloudcx

import "github.com/gildas/go-errors"

type OpenMessageTo struct {
	ID   string `json:"id"`
	Type string `json:"idType,omitempty"`
}

// Validate checks if the object is valid
func (to *OpenMessageTo) Validate() (err error) {
	if len(to.ID) == 0 {
		return errors.ArgumentMissing.With("channel.to.id")
	}
	return
}

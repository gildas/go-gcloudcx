package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type Flow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Division    *Division `json:"division,omitempty"`
	IsActive    bool      `json:"active"`
	IsSystem    bool      `json:"system"`
	IsDeleted   bool      `json:"deleted"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (flow Flow) GetID() uuid.UUID {
	return flow.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (flow Flow) GetURI() URI {
	return NewURI("/api/v2/flows/%s", flow.ID)
}

// String gets a string representation of this
//
// implements fmt.Stringer
func (flow Flow) String() string {
	return flow.Name
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (flow Flow) MarshalJSON() ([]byte, error) {
	type surrogate Flow
	data, err := json.Marshal(&struct {
		surrogate
		SelfURI URI `json:"selfUri"`
	}{
		surrogate: surrogate(flow),
		SelfURI:   flow.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

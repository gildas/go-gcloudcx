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

// Initialize initializes the object
//
// implements Initializable
func (flow *Flow) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			flow.ID = parameter
		}
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (flow Flow) GetID() uuid.UUID {
	return flow.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (flow Flow) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/flows/%s", ids[0])
	}
	if flow.ID != uuid.Nil {
		return NewURI("/api/v2/flows/%s", flow.ID)
	}
	return URI("/api/v2/flows/")
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

package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// ResponseManagementLibrary is the interface for the Response Management Library
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementLibrary struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	DateCreated time.Time        `json:"dateCreated,omitempty"`
	CreatedBy   *DomainEntityRef `json:"createdBy,omitempty"`
	Version     int              `json:"version"`
	SelfURI     URI              `json:"selfUri,omitempty"`
}

// Initialize initializes the object
//
// implements Initializable
func (library *ResponseManagementLibrary) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			library.ID = parameter
		}
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (library ResponseManagementLibrary) GetID() uuid.UUID {
	return library.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (library ResponseManagementLibrary) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/responsemanagement/libraries/%s", ids[0])
	}
	if library.ID != uuid.Nil {
		return NewURI("/api/v2/responsemanagement/libraries/%s", library.ID)
	}
	return URI("/api/v2/responsemanagement/libraries/")
}

// GetType gets the identifier of this
//
// implements Identifiable
func (library *ResponseManagementLibrary) GetType() string {
	return "responsemanagement.library"
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (library ResponseManagementLibrary) String() string {
	if len(library.Name) > 0 {
		return library.Name
	}
	return library.ID.String()
}

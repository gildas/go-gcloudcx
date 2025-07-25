package gcloudcx

import (
	"context"

	"github.com/gildas/go-core"
	"github.com/google/uuid"
)

// Identifiable describes that can get their Identifier as a UUID
type Identifiable interface {
	core.Identifiable
}

// Addressable describes things that carry a URI (typically /api/v2/things/{{uuid}})
type Addressable interface {
	// GetURI gets the URI
	//
	// if ids are provided, they are used to replace the {{uuid}} in the URI.
	//
	// if no ids are provided and the Addressable has a UUID, it is used to replace the {{uuid}} in the URI.
	//
	// else, the pattern for the URI is returned ("/api/v2/things/%s")
	GetURI(ids ...uuid.UUID) URI
}

// AddressableByStringID describes things that carry a URI (typically /api/v2/things/{{string}})
type AddressableByStringID interface {
	// GetURI gets the URI
	//
	// if ids are provided, they are used to replace the {{string}} in the URI.
	//
	// if no ids are provided and the Addressable has a UUID, it is used to replace the {{string}} in the URI.
	//
	// else, the pattern for the URI is returned ("/api/v2/things/%s")
	GetURI(ids ...string) URI
}

// Initializable describes things that can be initialized
type Initializable interface {
	Initialize(parameters ...interface{})
}

// Fetchable describes things that can be fetched from the Genesys Cloud API
type Fetchable interface {
	Identifiable
	Addressable
}

// FetchableByStringID describes things that can be fetched from the Genesys Cloud API
//
// These objects use a named ID instead of a UUID to fetch them.
//
// For example, IntegrationType uses a string ID like "genesyscloud-digital-bot-connector" or "amazon-lex-v2" instead of a UUID.
type FetchableByStringID interface {
	core.StringIdentifiable
	AddressableByStringID
}

// StateUpdater describes objects than can update the state of an Identifiable
type StateUpdater interface {
	UpdateState(context context.Context, identifiable Identifiable, state string) (correlationID string, err error)
}

// Disconnecter describes objects that can disconnect an Identifiable from themselves
type Disconnecter interface {
	Disconnect(context context.Context, identifiable Identifiable) (correlationID string, err error)
}

// Transferrer describes objects that can transfer an Identifiable somewhere else
type Transferrer interface {
	Transfer(context context.Context, identifiable Identifiable, target Identifiable) (correlationID string, err error)
}

// Address describes an Address (telno, etc)
type Address struct {
	Name               string `json:"name"`
	NameRaw            string `json:"nameRaw"`
	AddressDisplayable string `json:"addressDisplayable"`
	AddressRaw         string `json:"addressRaw"`
	AddressNormalized  string `json:"addressNormalized"`
}

// ErrorBody describes errors in PureCloud objects
type ErrorBody struct {
	Status            int               `json:"status"`
	Code              string            `json:"code"`
	EntityID          string            `json:"entityId"`
	EntityName        string            `json:"entityName"`
	Message           string            `json:"message"`
	MessageWithParams string            `json:"messageWithParams"`
	MessageParams     map[string]string `json:"messageParams"`
	ContextID         string            `json:"contextId"`
	Details           []ErrorDetail     `json:"details"`
	Errors            []*ErrorBody      `json:"errors"`
}

// ErrorDetail describes the details of an error
type ErrorDetail struct {
	ErrorCode  string `json:"errorCode"`
	Fieldname  string `json:"fieldName"`
	EntityID   string `json:"entityId"`
	EntityName string `json:"entityName"`
}

// Error returns a string representation of this error
func (err ErrorBody) Error() string {
	return err.MessageWithParams
}

// paginatedEntities describes a paginated list of entities
type paginatedEntities struct {
	PageSize    int `json:"pageSize"`
	PageNumber  int `json:"pageNumber"`
	PageCount   int `json:"pageCount"`
	EntityCount int `json:"total"`
	FirstURI    URI `json:"firstUri"`
	SelfURI     URI `json:"selfUri"`
	LastURI     URI `json:"lastUri"`
}

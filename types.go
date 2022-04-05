package gcloudcx

import (
	"context"

	"github.com/gildas/go-core"
)

// Identifiable describes that can get their Identifier as a UUID
type Identifiable interface {
	core.Identifiable
}

// Addressable describes things that carry a URI (typically /api/v2/things/{{uuid}})
type Addressable interface {
	GetURI() URI
}

// Initializable describes things that can be initialized
type Initializable interface {
	Initialize(parameters ...interface{}) error
}

// Fetchable describes things that can be fetched from the Genesys Cloud API
type Fetchable interface {
	Fetch(context context.Context, client *Client, parameters ...interface{}) error
}

// StateUpdater describes objects than can update the state of an Identifiable
type StateUpdater interface {
	UpdateState(context context.Context, identifiable Identifiable, state string) error
}

// Disconnecter describes objects that can disconnect an Identifiable from themselves
type Disconnecter interface {
	Disconnect(context context.Context, identifiable Identifiable) error
}

// Transferrer describes objects that can transfer an Identifiable somewhere else
type Transferrer interface {
	Transfer(context context.Context, identifiable Identifiable, target Identifiable) error
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

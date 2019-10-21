package purecloud

import (
	"net/url"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region              string             `json:"region"`
	DeploymentID        string             `json:"deploymentId"`
	Organization        *Organization      `json:"-"`
	API                 *url.URL           `json:"apiUrl,omitempty"`
	Proxy               *url.URL           `json:"proxyUrl,omitempty"`
	AuthorizationGrant  AuthorizationGrant `json:"auth"`
	Logger              *logger.Logger     `json:"-"`
}

// Identifiable describes things that carry an identifier
type Identifiable interface {
	// GetID gets the identifier of this
	GetID() string
}

// AddressableEntityRef describes an Entity that can be addressed
type AddressableEntityRef struct {
	ID      string `json:"id"`
	SelfURI string `json:"selfUri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref AddressableEntityRef) GetID() string {
	return ref.ID
}

// DomainEntityRef describes a DomainEntity Reference
type DomainEntityRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SelfURI  string `json:"self_uri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref DomainEntityRef) GetID() string {
	return ref.ID
}

// Address describes an Addres (telno, etc)
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
package gcloudcx

import "github.com/google/uuid"

// DomainRole describes a Role in a Domain
type DomainRole struct {
	// TODO: Find out what should be here!
	ID      uuid.UUID `json:"id"`
	SelfURI URI       `json:"selfUri,omitempty"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (role DomainRole) GetID() uuid.UUID {
	return role.ID
}

// GetURI gets the URI of this
//
//	implements Addressable
func (role DomainRole) GetURI() URI {
	return role.SelfURI
}

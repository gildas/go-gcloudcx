package gcloudcx

import "github.com/google/uuid"

// EntityRef describes an Entity that has an ID
type EntityRef struct {
	ID uuid.UUID `json:"id"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref EntityRef) GetID() uuid.UUID {
	return ref.ID
}

// AddressableEntityRef describes an Entity that can be addressed
type AddressableEntityRef struct {
	ID      uuid.UUID `json:"id"`
	SelfURI URI       `json:"selfUri,omitempty"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref AddressableEntityRef) GetID() uuid.UUID {
	return ref.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (ref AddressableEntityRef) GetURI() URI {
	return ref.SelfURI
}

// DomainEntityRef describes a DomainEntity Reference
type DomainEntityRef struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name,omitempty"`
	SelfURI URI       `json:"selfUri,omitempty"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref DomainEntityRef) GetID() uuid.UUID {
	return ref.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (ref DomainEntityRef) GetURI() URI {
	return ref.SelfURI
}

package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// UserStation describes a User Station
type UserStation struct {
	// TODO: Find out what should be here!
	ID             uuid.UUID         `json:"id"`
	SelfURI        URI               `json:"selfUri,omitempty"`
	Name           string            `json:"name"`
	Type           string            `json:"type"`
	AssociatedUser *User             `json:"associatedUser"`
	AssociatedDate time.Time         `json:"associatedDate"`
	DefaultUser    *User             `json:"defaultUser"`
	ProviderInfo   map[string]string `json:"providerInfo"`
}

// UserStations describes the stations of a user
type UserStations struct {
	AssociatedStation     *UserStation `json:"associatedStation"`
	LastAssociatedStation *UserStation `json:"lastAssociatedStation"`
	DefaultStation        *UserStation `json:"defaultStation"`
	EffectiveStation      *UserStation `json:"effectiveStation"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (station UserStation) GetID() uuid.UUID {
	return station.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (station UserStation) GetURI() URI {
	return station.SelfURI
}

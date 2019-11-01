package purecloud

import (
	"time"
)

// UserStations describes a User Station
type UserStation struct {
	// TODO: Find out what should be here!
	ID             string            `json:"id"`
	SelfURI        string            `json:"selfUri,omitempty"`
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
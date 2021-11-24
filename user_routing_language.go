package gcloudcx

import "github.com/google/uuid"

// UserRoutingLanguage describe a Routing Language for a User
type UserRoutingLanguage struct {
	ID          uuid.UUID `json:"id"`
	SelfURI     URI       `json:"selfUri"`
	Name        string    `json:"name"`
	LanguageURI string    `json:"languageUri"`
	State       string    `json:"state"`
	Proficiency float64   `json:"proficiency"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (userRoutingLanguage UserRoutingLanguage) GetID() uuid.UUID {
	return userRoutingLanguage.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (userRoutingLanguage UserRoutingLanguage) GetURI() URI {
	return userRoutingLanguage.SelfURI
}

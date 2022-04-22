package gcloudcx

import "github.com/google/uuid"

type AuthorizationGrantRole struct {
	ID          uuid.UUID                  `json:"id"`
	SelfUri     string                     `json:"selfUri"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	IsDefault   bool                       `json:"default"`
	Policies    []AuthorizationGrantPolicy `json:"policies"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (role AuthorizationGrantRole) GetID() uuid.UUID {
	return role.ID
}

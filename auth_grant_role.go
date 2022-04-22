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

// CheckScope checks if the grant role allows or denies the given scope
func (role AuthorizationGrantRole) CheckScope(scope AuthorizationScope) bool {
	for _, policy := range role.Policies {
		if policy.CheckScope(scope) {
			return true
		}
	}
	return false
}

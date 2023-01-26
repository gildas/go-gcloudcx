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
//
// If allowed, the policy that allows the scope is returned
func (role AuthorizationGrantRole) CheckScope(scope AuthorizationScope) (AuthorizationGrantPolicy, bool) {
	for _, policy := range role.Policies {
		if policy.CheckScope(scope) {
			return policy, true
		}
	}
	return AuthorizationGrantPolicy{}, false
}

// String returns a string representation of the AuthorizationDivision
//
// implements fmt.Stringer
func (role AuthorizationGrantRole) String() string {
	if len(role.Name) > 0 {
		return role.Name
	}
	return role.ID.String()
}

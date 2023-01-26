package gcloudcx

import "github.com/google/uuid"

type AuthorizationGrant struct {
	SubjectID uuid.UUID              `json:"subjectId"`
	Division  AuthorizationDivision  `json:"division"`
	Role      AuthorizationGrantRole `json:"role"`
	CreatedAt string                 `json:"grantMadeAt"` // TODO: this is an ISO8601 date
}

// CheckScope checks if the grant allows or denies the given scope
//
// If allowed, the policy that allows the scope is returned
func (grant AuthorizationGrant) CheckScope(scope AuthorizationScope) (AuthorizationGrantPolicy, bool) {
	return grant.Role.CheckScope(scope)
}

func (grant AuthorizationGrant) String() string {
	return grant.Role.String() + "@" + grant.Division.String()
}

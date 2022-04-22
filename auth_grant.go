package gcloudcx

import "github.com/google/uuid"

type AuthorizationGrant struct {
	SubjectID uuid.UUID              `json:"subjectId"`
	Division  AuthorizationDivision  `json:"division"`
	Role      AuthorizationGrantRole `json:"role"`
	CreatedAt string                 `json:"grantMadeAt"` // TODO: this is an ISO8601 date
}

// CheckScope checks if the grant allows or denies the given scope
func (grant AuthorizationGrant) CheckScope(scope AuthorizationScope) bool {
	return grant.Role.CheckScope(scope)
}

func (grant AuthorizationGrant) String() string {
	return grant.Role.String() + "@" + grant.Division.String()
}

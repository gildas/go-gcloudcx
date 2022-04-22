package gcloudcx

import "github.com/google/uuid"

type AuthorizationGrant struct {
	SubjectID uuid.UUID              `json:"subjectId"`
	Division  AuthorizationDivision  `json:"division"`
	Role      AuthorizationGrantRole `json:"role"`
	CreatedAt string                 `json:"grantMadeAt"` // TODO: this is an ISO8601 date
}

package gcloudcx

import (
	"github.com/google/uuid"
)

// AuthorizationSubject describes the roles and permissions of a Subject
type AuthorizationSubject struct {
	ID      uuid.UUID            `json:"id"`
	SelfUri string               `json:"selfUri"`
	Name    string               `json:"name"`
	Grants  []AuthorizationGrant `json:"grants"`
	Version int                  `json:"version"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (subject AuthorizationSubject) GetID() uuid.UUID {
	return subject.ID
}

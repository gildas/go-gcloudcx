package purecloud

import (
	"github.com/gildas/go-core"
	"github.com/google/uuid"
)

// Authorizer describes what a grants should do
type Authorizer interface {
	Authorize(client *Client) error // Authorize a client with PureCloud
	AccessToken() *AccessToken      // Get the Access Token obtained by the Authorizer
	core.Identifiable               // Implements core.Identifiable
}

// AuthorizationSubject describes the roles and permissions of a Subject
type AuthorizationSubject struct {
	ID      uuid.UUID            `json:"id"`
	SelfUri string               `json:"selfUri"`
	Name    string               `json:"name"`
	Grants  []AuthorizationGrant `json:"grants"`
	Version int                  `json:"version"`
}

type AuthorizationGrant struct {
	SubjectID      uuid.UUID              `json:"subjectId"`
	Division       AuthorizationDivision  `json:"division"`
	Role           AuthorizationGrantRole `json:"role"`
	CreatedAt      string                 `json:"grantMadeAt"` // TODO: this is an ISO8601 date
}

type AuthorizationDivision struct {
	ID           uuid.UUID      `json:"id"`
	SelfUri      string         `json:"selfUri"`
	Name         string         `json:"name"`
	Description  string         `json:"description"` // required
	IsHome       bool           `json:"homeDivision"`
	ObjectCounts map[string]int `json:"objectCounts"`
}

type AuthorizationGrantRole struct {
	ID          uuid.UUID                  `json:"id"`
	SelfUri     string                     `json:"selfUri"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	IsDefault   bool                       `json:"default"`
	Policies    []AuthorizationGrantPolicy `json:"policies"`
}

type AuthorizationGrantPolicy struct {
	EntityName string   `json:"entityName"`
	Domain     string   `json:"domain"`
	Condition  string   `json:"condition"`
	Actions    []string `json:"actions"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (subject AuthorizationSubject) GetID() uuid.UUID {
	return subject.ID
}

// GetID gets the identifier
//
// implements core.Identifiable
func (division AuthorizationDivision) GetID() uuid.UUID {
	return division.ID
}

// GetID gets the identifier
//
// implements core.Identifiable
func (role AuthorizationGrantRole) GetID() uuid.UUID {
	return role.ID
}

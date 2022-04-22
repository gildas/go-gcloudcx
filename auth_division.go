package gcloudcx

import "github.com/google/uuid"

type AuthorizationDivision struct {
	ID           uuid.UUID      `json:"id"`
	SelfUri      string         `json:"selfUri"`
	Name         string         `json:"name"`
	Description  string         `json:"description"` // required
	IsHome       bool           `json:"homeDivision"`
	ObjectCounts map[string]int `json:"objectCounts"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (division AuthorizationDivision) GetID() uuid.UUID {
	return division.ID
}

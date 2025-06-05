package gcloudcx

import (
	"encoding/json"
	"strings"

	"github.com/gildas/go-errors"
)

// StringCollectionSlotEntity represents a string collection slot entity in Genesys Cloud CX
type StringCollectionSlotEntity struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func init() {
	slotEntityRegistry.Add(StringCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity StringCollectionSlotEntity) GetType() string {
	return "StringCollection"
}

// GetName returns the name of the slot entity
func (entity StringCollectionSlotEntity) GetName() string {
	return entity.Name
}

// StringCollection returns the string representation of the slot entity's value
func (entity StringCollectionSlotEntity) String() string {
	return strings.Join(entity.Values, ",")
}

// MarshalJSON marshals the slot entity to JSON
func (entity StringCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate StringCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      entity.GetType(),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

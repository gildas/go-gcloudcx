package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// StringSlotEntity represents a string slot entity in Genesys Cloud CX
type StringSlotEntity struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func init() {
	slotEntityRegistry.Add(StringSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity StringSlotEntity) GetType() string {
	return "String"
}

// GetName returns the name of the slot entity
func (entity StringSlotEntity) GetName() string {
	return entity.Name
}

// String returns the string representation of the slot entity's value
func (entity StringSlotEntity) String() string {
	return entity.Value
}

// MarshalJSON marshals the slot entity to JSON
func (entity StringSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate StringSlotEntity

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      entity.GetType(),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

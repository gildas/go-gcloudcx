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

// ParseValue parses the value and returns a new SlotEntity instance
func (entity StringSlotEntity) ParseValue(value string) (SlotEntity, error) {
	return StringSlotEntity{
		Name:  entity.Name,
		Value: value,
	}, nil
}

// String returns the string representation of the slot entity's value
func (entity StringSlotEntity) String() string {
	return entity.Value
}

// Validate checks if the slot entity is valid
func (entity *StringSlotEntity) Validate() error {
	var merr errors.MultiError

	if len(entity.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("entity.name"))
	}
	if len(entity.Name) > 100 {
		merr.Append(errors.ArgumentInvalid.With("entity.name", "must be less than 100 characters"))
	}

	return merr.AsError()
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

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *StringSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate StringSlotEntity

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = StringSlotEntity(inner.surrogate)
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

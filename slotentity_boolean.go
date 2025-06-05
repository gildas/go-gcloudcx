package gcloudcx

import (
	"encoding/json"
	"strconv"

	"github.com/gildas/go-errors"
)

// BooleanSlotEntity represents a boolean slot entity
type BooleanSlotEntity struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

func init() {
	slotEntityRegistry.Add(BooleanSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity BooleanSlotEntity) GetType() string {
	return "Boolean"
}

// GetName returns the name of the slot entity
func (entity BooleanSlotEntity) GetName() string {
	return entity.Name
}

// Boolean returns the string representation of the slot entity's value
func (entity BooleanSlotEntity) String() string {
	return strconv.FormatBool(entity.Value)
}

// MarshalJSON marshals the slot entity to JSON
func (entity BooleanSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate BooleanSlotEntity

	data, err := json.Marshal(struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}{
		Type:      entity.GetType(),
		Value:     entity.String(),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *BooleanSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate BooleanSlotEntity

	var inner struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = BooleanSlotEntity(inner.surrogate)
	entity.Value, err = strconv.ParseBool(inner.Value)
	if err != nil {
		return errors.Join(errors.JSONUnmarshalError, errors.ArgumentInvalid.With("value", inner.Value, "boolean"), err)
	}

	return nil
}

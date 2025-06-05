package gcloudcx

import (
	"encoding/json"
	"strconv"

	"github.com/gildas/go-errors"
)

// IntegerSlotEntity represents an integer slot entity
type IntegerSlotEntity struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func init() {
	slotEntityRegistry.Add(IntegerSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity IntegerSlotEntity) GetType() string {
	return "Integer"
}

// GetName returns the name of the slot entity
func (entity IntegerSlotEntity) GetName() string {
	return entity.Name
}

// Integer returns the string representation of the slot entity's value
func (entity IntegerSlotEntity) String() string {
	return strconv.FormatInt(entity.Value, 10)
}

// MarshalJSON marshals the slot entity to JSON
func (entity IntegerSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate IntegerSlotEntity

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
func (entity *IntegerSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate IntegerSlotEntity

	var inner struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = IntegerSlotEntity(inner.surrogate)
	entity.Value, err = strconv.ParseInt(inner.Value, 10, 64)
	if err != nil {
		return errors.Join(errors.JSONUnmarshalError, errors.ArgumentInvalid.With("value", inner.Value, "integer"), err)
	}

	return nil
}

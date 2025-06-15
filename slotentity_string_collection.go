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

// ParseValue parses the value and returns a new SlotEntity instance
func (entity StringCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []string
	for _, raw := range strings.Split(value, ",") {
		values = append(values, strings.TrimSpace(raw))
	}
	return &StringCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// StringCollection returns the string representation of the slot entity's value
func (entity StringCollectionSlotEntity) String() string {
	return strings.Join(entity.Values, ",")
}

// Validate checks if the slot entity is valid
func (entity *StringCollectionSlotEntity) Validate() error {
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

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *StringCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate StringCollectionSlotEntity

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = StringCollectionSlotEntity(inner.surrogate)
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

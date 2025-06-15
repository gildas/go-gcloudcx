package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// IntegerCollectionSlotEntity represents an integer collection slot entity
type IntegerCollectionSlotEntity struct {
	Name   string  `json:"name"`
	Values []int64 `json:"values"`
}

func init() {
	slotEntityRegistry.Add(IntegerCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity IntegerCollectionSlotEntity) GetType() string {
	return "IntegerCollection"
}

// GetName returns the name of the slot entity
func (entity IntegerCollectionSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity IntegerCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []int64
	for index, raw := range strings.Split(value, ",") {
		intValue, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		values = append(values, intValue)
	}
	return &IntegerCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// IntegerCollection returns the string representation of the slot entity's value
func (entity IntegerCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value int64) string { return strconv.FormatInt(value, 10) }), ",")
}

// Validate checks if the slot entity is valid
func (entity *IntegerCollectionSlotEntity) Validate() error {
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
func (entity IntegerCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate IntegerCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(value int64) string { return strconv.FormatInt(value, 10) }),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *IntegerCollectionSlotEntity) UnmarshalJSON(payload []byte) error {
	type surrogate IntegerCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = IntegerCollectionSlotEntity(inner.surrogate)
	entity.Values = core.Map(inner.Values, func(value string) int64 {
		intValue, _ := strconv.ParseInt(value, 10, 64)
		return intValue
	})
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

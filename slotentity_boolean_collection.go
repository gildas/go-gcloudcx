package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// BooleanCollectionSlotEntity represents a boolean collection slot entity
type BooleanCollectionSlotEntity struct {
	Name   string `json:"name"`
	Values []bool `json:"values"`
}

func init() {
	slotEntityRegistry.Add(BooleanCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity BooleanCollectionSlotEntity) GetType() string {
	return "BooleanCollection"
}

// GetName returns the name of the slot entity
func (entity BooleanCollectionSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity BooleanCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var boolValues []bool
	for index, raw := range strings.Split(value, ",") {
		boolValue, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		boolValues = append(boolValues, boolValue)
	}
	return &BooleanCollectionSlotEntity{
		Name:   entity.Name,
		Values: boolValues,
	}, nil
}

// BooleanCollection returns the string representation of the slot entity's value
func (entity BooleanCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value bool) string { return strconv.FormatBool(value) }), ",")
}

// Validate checks if the slot entity is valid
func (entity *BooleanCollectionSlotEntity) Validate() error {
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
func (entity BooleanCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate BooleanCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(value bool) string { return strconv.FormatBool(value) }),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *BooleanCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate BooleanCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = BooleanCollectionSlotEntity(inner.surrogate)
	entity.Values = core.Map(inner.Values, func(value string) bool {
		b, _ := strconv.ParseBool(value)
		return b
	})
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// DecimalCollectionSlotEntity represents a decimal collection slot entity
type DecimalCollectionSlotEntity struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

func init() {
	slotEntityRegistry.Add(DecimalCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DecimalCollectionSlotEntity) GetType() string {
	return "DecimalCollection"
}

// GetName returns the name of the slot entity
func (entity DecimalCollectionSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity DecimalCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []float64
	for index, raw := range strings.Split(value, ",") {
		floatValue, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		values = append(values, floatValue)
	}
	return &DecimalCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// DecimalCollection returns the string representation of the slot entity's value
func (entity DecimalCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }), ",")
}

// Validate checks if the slot entity is valid
func (entity *DecimalCollectionSlotEntity) Validate() error {
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
func (entity DecimalCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DecimalCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *DecimalCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DecimalCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = DecimalCollectionSlotEntity(inner.surrogate)
	entity.Values = core.Map(inner.Values, func(value string) float64 {
		floatValue, _ := strconv.ParseFloat(value, 64)
		return floatValue
	})
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

package gcloudcx

import (
	"encoding/json"
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

// DecimalCollection returns the string representation of the slot entity's value
func (entity DecimalCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }), ",")
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

	return nil
}

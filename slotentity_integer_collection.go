package gcloudcx

import (
	"encoding/json"
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

// IntegerCollection returns the string representation of the slot entity's value
func (entity IntegerCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value int64) string { return strconv.FormatInt(value, 10) }), ",")
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

	return nil
}

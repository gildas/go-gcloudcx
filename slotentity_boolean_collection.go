package gcloudcx

import (
	"encoding/json"
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

// BooleanCollection returns the string representation of the slot entity's value
func (entity BooleanCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value bool) string { return strconv.FormatBool(value) }), ",")
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

	return nil
}

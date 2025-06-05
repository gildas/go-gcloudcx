package gcloudcx

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// DatetimeCollectionSlotEntity represents a datetime collection slot entity
type DatetimeCollectionSlotEntity struct {
	Name   string      `json:"name"`
	Values []time.Time `json:"values"`
}

func init() {
	slotEntityRegistry.Add(DatetimeCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DatetimeCollectionSlotEntity) GetType() string {
	return "DatetimeCollection"
}

// GetName returns the name of the slot entity
func (entity DatetimeCollectionSlotEntity) GetName() string {
	return entity.Name
}

// DatetimeCollection returns the string representation of the slot entity's value
func (entity DatetimeCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value time.Time) string { return value.Format(time.RFC3339) }), ",")
}

// MarshalJSON marshals the slot entity to JSON
func (entity DatetimeCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DatetimeCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(value time.Time) string { return value.Format(time.RFC3339) }),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *DatetimeCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DatetimeCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = DatetimeCollectionSlotEntity(inner.surrogate)
	entity.Values = core.Map(inner.Values, func(value string) time.Time {
		datetime, _ := time.Parse(time.RFC3339, value)
		return datetime
	})
	return nil
}

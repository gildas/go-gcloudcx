package gcloudcx

import (
	"encoding/json"
	"fmt"
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

// ParseValue parses the value and returns a new SlotEntity instance
func (entity DatetimeCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []time.Time
	for index, raw := range strings.Split(value, ",") {
		datetime, err := time.Parse(time.RFC3339, strings.TrimSpace(raw))
		if err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		values = append(values, datetime)
	}
	return &DatetimeCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// DatetimeCollection returns the string representation of the slot entity's value
func (entity DatetimeCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value time.Time) string { return value.Format(time.RFC3339) }), ",")
}

// Validate checks if the slot entity is valid
func (entity *DatetimeCollectionSlotEntity) Validate() error {
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
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// DurationCollectionSlotEntity represents a string collection slot entity
type DurationCollectionSlotEntity struct {
	Name   string          `json:"name"`
	Values []time.Duration `json:"values"`
}

func init() {
	slotEntityRegistry.Add(DurationCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DurationCollectionSlotEntity) GetType() string {
	return "DurationCollection"
}

// GetName returns the name of the slot entity
func (entity DurationCollectionSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity DurationCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []time.Duration
	for index, raw := range strings.Split(value, ",") {
		duration, err := core.ParseDuration(strings.TrimSpace(raw))
		if err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		values = append(values, duration)
	}
	return &DurationCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// DurationCollection returns the string representation of the slot entity's value
func (entity DurationCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(duration time.Duration) string { return core.Duration(duration).ToISO8601() }), ",")
}

// Validate checks if the slot entity is valid
func (entity *DurationCollectionSlotEntity) Validate() error {
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
func (entity DurationCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DurationCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(duration time.Duration) string { return core.Duration(duration).ToISO8601() }),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *DurationCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DurationCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = DurationCollectionSlotEntity(inner.surrogate)
	entity.Values = core.Map(inner.Values, func(value string) time.Duration {
		duration, err := core.ParseDuration(value)
		if err != nil {
			return 0 // or handle the error as needed
		}
		return duration
	})
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

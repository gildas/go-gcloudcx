package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// DurationSlotEntity represents a duration slot entity
type DurationSlotEntity struct {
	Name  string        `json:"name"`
	Value time.Duration `json:"value"`
}

func init() {
	slotEntityRegistry.Add(DurationSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DurationSlotEntity) GetType() string {
	return "Duration"
}

// GetName returns the name of the slot entity
func (entity DurationSlotEntity) GetName() string {
	return entity.Name
}

// Duration returns the string representation of the slot entity's value
func (entity DurationSlotEntity) String() string {
	return core.Duration(entity.Value).ToISO8601()
}

// MarshalJSON marshals the slot entity to JSON
func (entity DurationSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DurationSlotEntity

	data, err := json.Marshal(struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}{
		Type:      entity.GetType(),
		Value:     entity.String(),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *DurationSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DurationSlotEntity

	var inner struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = DurationSlotEntity(inner.surrogate)
	entity.Value, err = core.ParseDuration(inner.Value)
	if err != nil {
		return errors.Join(errors.JSONUnmarshalError, errors.ArgumentInvalid.With("value", inner.Value, "duration"), err)
	}

	return nil
}

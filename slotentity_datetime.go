package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// DatetimeSlotEntity represents a datetime slot entity
type DatetimeSlotEntity struct {
	Name  string    `json:"name"`
	Value time.Time `json:"value"`
}

func init() {
	slotEntityRegistry.Add(DatetimeSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DatetimeSlotEntity) GetType() string {
	return "Datetime"
}

// GetName returns the name of the slot entity
func (entity DatetimeSlotEntity) GetName() string {
	return entity.Name
}

// Datetime returns the string representation of the slot entity's value
func (entity DatetimeSlotEntity) String() string {
	return entity.Value.Format(time.RFC3339)
}

// MarshalJSON marshals the slot entity to JSON
func (entity DatetimeSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DatetimeSlotEntity

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
func (entity *DatetimeSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DatetimeSlotEntity

	var inner struct {
		Type  string    `json:"type"`
		Value core.Time `json:"value"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*entity = DatetimeSlotEntity(inner.surrogate)
	entity.Value = time.Time(inner.Value)
	return nil
}

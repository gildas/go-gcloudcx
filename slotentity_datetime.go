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

// ParseValue parses the value and returns a new SlotEntity instance
func (entity DatetimeSlotEntity) ParseValue(value string) (SlotEntity, error) {
	datetime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errors.ArgumentInvalid.With("value", value, "datetime")
	}
	return &DatetimeSlotEntity{
		Name:  entity.Name,
		Value: datetime,
	}, nil
}

// Datetime returns the string representation of the slot entity's value
func (entity DatetimeSlotEntity) String() string {
	return entity.Value.Format(time.RFC3339)
}

// Validate checks if the slot entity is valid
func (entity *DatetimeSlotEntity) Validate() error {
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
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

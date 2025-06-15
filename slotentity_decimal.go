package gcloudcx

import (
	"encoding/json"
	"strconv"

	"github.com/gildas/go-errors"
)

// DecimalSlotEntity represents a decimal slot entity
type DecimalSlotEntity struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func init() {
	slotEntityRegistry.Add(DecimalSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity DecimalSlotEntity) GetType() string {
	return "Decimal"
}

// GetName returns the name of the slot entity
func (entity DecimalSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity DecimalSlotEntity) ParseValue(value string) (SlotEntity, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, errors.ArgumentInvalid.With("value", value, "decimal")
	}
	return &DecimalSlotEntity{
		Name:  entity.Name,
		Value: floatValue,
	}, nil
}

// Decimal returns the string representation of the slot entity's value
func (entity DecimalSlotEntity) String() string {
	return strconv.FormatFloat(entity.Value, 'f', -1, 64)
}

// Validate checks if the slot entity is valid
func (entity *DecimalSlotEntity) Validate() error {
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
func (entity DecimalSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate DecimalSlotEntity

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
func (entity *DecimalSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate DecimalSlotEntity

	var inner struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*entity = DecimalSlotEntity(inner.surrogate)
	if len(inner.Value) > 0 {
		entity.Value, err = strconv.ParseFloat(inner.Value, 64)
		if err != nil {
			return errors.Join(errors.JSONUnmarshalError, errors.ArgumentInvalid.With("value", inner.Value, "decimal"), err)
		}
	}
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

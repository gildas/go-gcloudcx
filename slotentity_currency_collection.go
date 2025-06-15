package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// CurrencyCollectionSlotEntity represents a currency collection slot entity
type CurrencyCollectionSlotEntity struct {
	Name   string          `json:"name"`
	Values []CurrencyValue `json:"values"`
}

func init() {
	slotEntityRegistry.Add(CurrencyCollectionSlotEntity{})
}

// GetType returns the type of the slot entity
func (entity CurrencyCollectionSlotEntity) GetType() string {
	return "CurrencyCollection"
}

// GetName returns the name of the slot entity
func (entity CurrencyCollectionSlotEntity) GetName() string {
	return entity.Name
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity CurrencyCollectionSlotEntity) ParseValue(value string) (SlotEntity, error) {
	var values []CurrencyValue
	for index, raw := range strings.Split(value, ",") {
		var currency CurrencyValue
		if err := currency.Parse(strings.TrimSpace(raw)); err != nil {
			return nil, errors.ArgumentInvalid.With(fmt.Sprintf("value[%d]", index), value)
		}
		values = append(values, currency)
	}
	return &CurrencyCollectionSlotEntity{
		Name:   entity.Name,
		Values: values,
	}, nil
}

// CurrencyCollection returns the string representation of the slot entity's value
func (entity CurrencyCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value CurrencyValue) string { return value.String() }), ",")
}

// Validate checks if the slot entity is valid
func (entity *CurrencyCollectionSlotEntity) Validate() error {
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
func (entity CurrencyCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate CurrencyCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type: entity.GetType(),
		Values: core.Map(entity.Values, func(value CurrencyValue) string {
			return fmt.Sprintf(`{"amount": %.2f, "code": "%s"}`, value.Amount, value.Currency)
		}),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *CurrencyCollectionSlotEntity) UnmarshalJSON(data []byte) (err error) {
	type surrogate CurrencyCollectionSlotEntity

	var inner struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	entity.Name = inner.Name
	entity.Values = core.Map(inner.Values, func(value string) CurrencyValue {
		var currencyValue CurrencyValue
		if err := currencyValue.Parse(value); err != nil {
			return CurrencyValue{}
		}
		return currencyValue
	})

	if inner.Type != entity.GetType() {
		return errors.InvalidType.With(inner.Type, entity.GetType())
	}
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

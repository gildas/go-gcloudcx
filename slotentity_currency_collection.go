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

// CurrencyValue is a helper type to represent the currency value in JSON
type CurrencyValue struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
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

// CurrencyCollection returns the string representation of the slot entity's value
func (entity CurrencyCollectionSlotEntity) String() string {
	return strings.Join(core.Map(entity.Values, func(value CurrencyValue) string { return value.String() }), ",")
}

// String returns the string representation of the currency value
func (value CurrencyValue) String() string {
	return fmt.Sprintf(`{"amount": %.2f, "code": "%s"}`, value.Amount, value.Currency)
}

// MarshalJSON marshals the slot entity to JSON
func (entity CurrencyCollectionSlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate CurrencyCollectionSlotEntity

	data, err := json.Marshal(struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
		surrogate
	}{
		Type:      entity.GetType(),
		Values:    core.Map(entity.Values, func(value CurrencyValue) string { return value.String() }),
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
		var cv CurrencyValue
		_, _ = fmt.Sscanf(value, `{"amount": %f, "code": "%s"}`, &cv.Amount, &cv.Currency)
		return cv
	})

	if inner.Type != entity.GetType() {
		return errors.InvalidType.With(inner.Type, entity.GetType())
	}

	return nil
}

package gcloudcx

import (
	"encoding/json"
	"fmt"

	"github.com/gildas/go-errors"
)

// CurrencySlotEntity represents a currency slot entity
type CurrencySlotEntity struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"-"`
	Currency string  `json:"-"`
}

// CurrencyValue is a helper type to represent the currency value in JSON
type CurrencyValue struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func init() {
	slotEntityRegistry.Add(CurrencySlotEntity{})
}

// GetType returns the type of the slot entity
func (entity CurrencySlotEntity) GetType() string {
	return "Currency"
}

// GetName returns the name of the slot entity
func (entity CurrencySlotEntity) GetName() string {
	return entity.Name
}

// Parse parses a string representation of a currency value
func (value *CurrencyValue) Parse(raw string) error {
	var amount float64
	var currency string

	n, err := fmt.Sscanf(raw, "%f %s", &amount, &currency)
	if err != nil || n != 2 {
		return errors.ArgumentInvalid.With("value", raw)
	}
	value.Amount = amount
	value.Currency = currency
	return nil
}

// String returns the string representation of the currency value
func (value CurrencyValue) String() string {
	return fmt.Sprintf(`{"amount": %.2f, "code": "%s"}`, value.Amount, value.Currency)
}

// ParseValue parses the value and returns a new SlotEntity instance
func (entity CurrencySlotEntity) ParseValue(value string) (SlotEntity, error) {
	var currencyValue CurrencyValue

	if err := currencyValue.Parse(value); err != nil {
		return nil, errors.ArgumentInvalid.With("value", value)
	}
	return &CurrencySlotEntity{
		Name:     entity.Name,
		Amount:   currencyValue.Amount,
		Currency: currencyValue.Currency,
	}, nil
}

// Currency returns the string representation of the slot entity's value
func (entity CurrencySlotEntity) String() string {
	return fmt.Sprintf("%.2f %s", entity.Amount, entity.Currency)
}

// Validate checks if the slot entity is valid
func (entity *CurrencySlotEntity) Validate() error {
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
func (entity CurrencySlotEntity) MarshalJSON() ([]byte, error) {
	type surrogate CurrencySlotEntity

	data, err := json.Marshal(struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}{
		Type:      entity.GetType(),
		Value:     fmt.Sprintf(`{"amount": %.2f, "code": "%s"}`, entity.Amount, entity.Currency),
		surrogate: surrogate(entity),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the slot entity from JSON
func (entity *CurrencySlotEntity) UnmarshalJSON(payload []byte) (err error) {
	type surrogate CurrencySlotEntity

	var inner struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		surrogate
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*entity = CurrencySlotEntity(inner.surrogate)

	var data struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"code"`
	}
	if err = json.Unmarshal([]byte(inner.Value), &data); err != nil {
		return errors.Join(errors.JSONUnmarshalError, errors.ArgumentInvalid.With("value", inner.Value, "currency"), err)
	}
	entity.Amount = data.Amount
	entity.Currency = data.Currency
	return errors.JSONUnmarshalError.Wrap(entity.Validate())
}

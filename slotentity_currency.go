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

// Currency returns the string representation of the slot entity's value
func (entity CurrencySlotEntity) String() string {
	return fmt.Sprintf(`{"amount": %.2f, "code": "%s"}`, entity.Amount, entity.Currency)
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
		Value:     entity.String(),
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
	return nil
}

package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type AuthorizationDivision struct {
	ID           uuid.UUID      `json:"-"`
	SelfUri      string         `json:"selfUri"`
	Name         string         `json:"name"`
	Description  string         `json:"description"` // required
	IsHome       bool           `json:"homeDivision"`
	ObjectCounts map[string]int `json:"objectCounts"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (division AuthorizationDivision) GetID() uuid.UUID {
	return division.ID
}

// String returns a string representation of the AuthorizationDivision
//
// implements fmt.Stringer
func (division AuthorizationDivision) String() string {
	if len(division.Name) > 0 {
		return division.Name
	}
	return division.ID.String()
}

func (division AuthorizationDivision) MarshalJSON() ([]byte, error) {
	type surrogate AuthorizationDivision
	inner := struct {
		surrogate
		ID string `json:"id"`
	}{
		surrogate: surrogate(division),
		ID:        "*",
	}
	if division.ID != uuid.Nil {
		inner.ID = division.ID.String()
	}
	data, err := json.Marshal(inner)
	return data, errors.JSONMarshalError.Wrap(err)
}

func (division *AuthorizationDivision) UnmarshalJSON(payload []byte) (err error) {
	type surrogate AuthorizationDivision
	var inner struct {
		surrogate
		ID string `json:"id"`
	}
	if err = json.Unmarshal(payload, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	(*division) = AuthorizationDivision(inner.surrogate)
	if inner.ID != "*" {
		if division.ID, err = uuid.Parse(inner.ID); err != nil {
			return errors.WrapErrors(errors.JSONMarshalError, errors.ArgumentInvalid.With("id", inner.ID), err)
		}
	}
	if division.ObjectCounts == nil {
		division.ObjectCounts = make(map[string]int)
	}
	return nil
}

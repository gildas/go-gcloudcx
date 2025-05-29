package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// CrendentialType represents the type of credential
type CredentialType struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder []string  `json:"displayOrder"`
	Required     []string  `json:"required"`
	Properties   any       `json:"properties"`
}

// MarshalJSON customizes the JSON encoding of CredentialType
func (credentialType CredentialType) MarshalJSON() ([]byte, error) {
	type surrogate CredentialType
	data, err := json.Marshal(struct {
		surrogate
		ID core.UUID `json:"id"`
	}{
		surrogate: surrogate(credentialType),
		ID:        core.UUID(credentialType.ID),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON customizes the JSON decoding of CredentialType
func (credentialType *CredentialType) UnmarshalJSON(payload []byte) (err error) {
	type surrogate CredentialType
	var inner struct {
		surrogate
		ID core.UUID `json:"id"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*credentialType = CredentialType(inner.surrogate)
	credentialType.ID = uuid.UUID(inner.ID)
	return nil
}

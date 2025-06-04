package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// CrendentialType represents the type of credential
type IntegrationCredentialType struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder []string  `json:"displayOrder,omitempty"`
	Required     []string  `json:"required,omitempty"`
	Properties   any       `json:"properties,omitempty"`
}

// MarshalJSON customizes the JSON encoding of CredentialType
func (credentialType IntegrationCredentialType) MarshalJSON() ([]byte, error) {
	type surrogate IntegrationCredentialType
	data, err := json.Marshal(struct {
		surrogate
		ID core.UUID `json:"id,omitempty"`
	}{
		surrogate: surrogate(credentialType),
		ID:        core.UUID(credentialType.ID),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON customizes the JSON decoding of CredentialType
func (credentialType *IntegrationCredentialType) UnmarshalJSON(payload []byte) (err error) {
	type surrogate IntegrationCredentialType
	var inner struct {
		surrogate
		ID core.UUID `json:"id"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*credentialType = IntegrationCredentialType(inner.surrogate)
	credentialType.ID = uuid.UUID(inner.ID)
	return nil
}

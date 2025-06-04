package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// IntegrationCredential represents a credential for an integration.
type IntegrationCredential struct {
	ID               uuid.UUID                 `json:"id"`
	Name             string                    `json:"name,omitempty"`
	Type             IntegrationCredentialType `json:"type"`
	CredentialFields map[string]string         `json:"credentialFields"`
}

// MarshalJSON marshals into JSON
//
// implements json.Marshaler
func (credential IntegrationCredential) MarshalJSON() ([]byte, error) {
	type surrogate IntegrationCredential
	data, err := json.Marshal(struct {
		ID core.UUID `json:"id,omitempty"`
		surrogate
	}{
		ID:        core.UUID(credential.ID),
		surrogate: surrogate(credential),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UnmarshalJSON unmarshals from JSON
//
// implements json.Unmarshaler
func (credential *IntegrationCredential) UnmarshalJSON(payload []byte) error {
	type surrogate IntegrationCredential
	var inner struct {
		surrogate
		ID core.UUID `json:"id"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*credential = IntegrationCredential(inner.surrogate)
	credential.ID = uuid.UUID(inner.ID)
	return nil
}

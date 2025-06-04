package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// IntegrationCredentialInfo represents the information about a credential
type IntegrationCredentialInfo struct {
	ID         uuid.UUID                 `json:"id"`
	Name       string                    `json:"name"`
	Type       IntegrationCredentialType `json:"type"`
	CreatedAt  time.Time                 `json:"createdDate"`
	ModifiedAt time.Time                 `json:"modifiedDate"`
}

// MarshalJSON customizes the JSON encoding of CredentialInfo
func (credentialInfo IntegrationCredentialInfo) MarshalJSON() ([]byte, error) {
	type surrogate IntegrationCredentialInfo

	var credentialType *IntegrationCredentialType

	if credentialInfo.Type.ID != uuid.Nil {
		credentialType = &credentialInfo.Type
	}

	data, err := json.Marshal(struct {
		surrogate
		ID         string                     `json:"id"`
		Type       *IntegrationCredentialType `json:"type,omitempty"`
		CreatedAt  core.Time                  `json:"createdDate,omitempty"`
		ModifiedAt core.Time                  `json:"modifiedDate,omitempty"`
	}{
		surrogate:  surrogate(credentialInfo),
		ID:         credentialInfo.ID.String(),
		Type:       credentialType,
		CreatedAt:  core.Time(credentialInfo.CreatedAt),
		ModifiedAt: core.Time(credentialInfo.ModifiedAt),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UnmarshalJSON customizes the JSON decoding of CredentialInfo
func (credentialInfo *IntegrationCredentialInfo) UnmarshalJSON(payload []byte) error {
	type surrogate IntegrationCredentialInfo
	var inner struct {
		surrogate
		ID         core.UUID `json:"id"`
		CreatedAt  core.Time `json:"createdDate"`
		ModifiedAt core.Time `json:"modifiedDate"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*credentialInfo = IntegrationCredentialInfo(inner.surrogate)
	credentialInfo.ID = uuid.UUID(inner.ID)
	credentialInfo.CreatedAt = inner.CreatedAt.AsTime()
	credentialInfo.ModifiedAt = inner.ModifiedAt.AsTime()
	return nil
}

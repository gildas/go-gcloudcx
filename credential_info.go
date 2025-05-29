package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// CredentialInfo represents the information about a credential
type CredentialInfo struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Type       CredentialType `json:"type"`
	CreatedAt  time.Time      `json:"createdDate"`
	ModifiedAt time.Time      `json:"modifiedDate"`
}

// MarshalJSON customizes the JSON encoding of CredentialInfo
func (credentialInfo CredentialInfo) MarshalJSON() ([]byte, error) {
	type surrogate CredentialInfo
	data, err := json.Marshal(struct {
		surrogate
		ID         string    `json:"id"`
		CreatedAt  core.Time `json:"createdDate"`
		ModifiedAt core.Time `json:"modifiedDate"`
	}{
		surrogate:  surrogate(credentialInfo),
		ID:         credentialInfo.ID.String(),
		CreatedAt:  core.Time(credentialInfo.CreatedAt),
		ModifiedAt: core.Time(credentialInfo.ModifiedAt),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UnmarshalJSON customizes the JSON decoding of CredentialInfo
func (credentialInfo *CredentialInfo) UnmarshalJSON(payload []byte) error {
	type surrogate CredentialInfo
	var inner struct {
		surrogate
		ID         core.UUID `json:"id"`
		CreatedAt  core.Time `json:"createdDate"`
		ModifiedAt core.Time `json:"modifiedDate"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*credentialInfo = CredentialInfo(inner.surrogate)
	credentialInfo.ID = uuid.UUID(inner.ID)
	credentialInfo.CreatedAt = inner.CreatedAt.AsTime()
	credentialInfo.ModifiedAt = inner.ModifiedAt.AsTime()
	return nil
}

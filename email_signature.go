package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type EmailSignature struct {
	Enabled          bool      `json:"enabled"`
	CannedResponseID uuid.UUID `json:"cannedResponseId"`
	AlwaysIncluded   bool      `json:"alwaysIncluded"`
	InclusionType    string    `json:"inclusionType"` // Draft, Send, SendOnce
}

// UnmarshalJSON unmarshals the email signature from JSON
//
// Implements json.UnMarshaler
func (signature *EmailSignature) UnmarshalJSON(data []byte) error {
	type surrogate EmailSignature
	var inner struct {
		surrogate
		CannedResponseID core.UUID `json:"cannedResponseId"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*signature = EmailSignature(inner.surrogate)
	signature.CannedResponseID = uuid.UUID(inner.CannedResponseID)

	return nil
}

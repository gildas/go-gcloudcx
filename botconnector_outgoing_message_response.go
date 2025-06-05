package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type BotConnectorOutgoingMessageResponse struct {
	MessageID uuid.UUID `json:"messageId"` // The ID of the message being responded to.
}

// Validate validates the outgoing message response
func (response *BotConnectorOutgoingMessageResponse) Validate() error {
	var merr errors.MultiError

	if response.MessageID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("messageId"))
	}
	return merr.AsError()
}

// MarshalJSON marshals the BotConnectorOutgoingMessageResponse to JSON
//
// implements core.Marshaler
func (response BotConnectorOutgoingMessageResponse) MarshalJSON() ([]byte, error) {
	type surrogate BotConnectorOutgoingMessageResponse
	data, err := json.Marshal(struct {
		surrogate
		MessageID core.UUID `json:"messageId"`
	}{
		surrogate: surrogate(response),
		MessageID: core.UUID(response.MessageID),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON data into a BotConnectorOutgoingMessageResponse
//
// implements core.Unmarshaler
func (response *BotConnectorOutgoingMessageResponse) UnmarshalJSON(data []byte) error {
	type surrogate BotConnectorOutgoingMessageResponse
	var inner struct {
		surrogate
		MessageID core.UUID `json:"messageId"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*response = BotConnectorOutgoingMessageResponse(inner.surrogate)
	response.MessageID = uuid.UUID(inner.MessageID)
	return nil
}

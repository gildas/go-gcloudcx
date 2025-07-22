package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// TODO: This will need to go to go-gcloudcx

// BotConnectorIncomingMessageRequest represents a message request received from a Genesys Cloud Digital Bot
type BotConnectorIncomingMessageRequest struct {
	BotID             string            `json:"botId"`
	BotVersion        string            `json:"botVersion"`
	BotSessionID      uuid.UUID         `json:"botSessionId"`
	ConversationID    uuid.UUID         `json:"genesysConversationId"`
	MessageID         uuid.UUID         `json:"messageId"`
	Language          string            `json:"languageCode"`
	BotSessionTimeout int64             `json:"botSessionTimeout"`
	Message           NormalizedMessage `json:"inputMessage"`
	Parameters        map[string]string `json:"parameters"`
	CorrelationID     string            `json:"correlationId,omitempty"` // Optional correlation ID for tracking purposes
}

// GetType returns the type of the message
//
// implements core.TypeCarrier
func (request BotConnectorIncomingMessageRequest) GetType() string {
	return request.Message.GetType()
}

// Validate validates the incoming message request
func (request *BotConnectorIncomingMessageRequest) Validate() error {
	var merr errors.MultiError
	if len(request.BotID) == 0 {
		merr.Append(errors.ArgumentMissing.With("botId"))
	}
	if request.BotSessionID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("botSessionId"))
	}
	if request.ConversationID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("genesysConversationId"))
	}
	if request.MessageID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("messageId"))
	}
	if request.Language == "" {
		merr.Append(errors.ArgumentMissing.With("languageCode"))
	}
	if request.BotSessionTimeout <= 0 {
		merr.Append(errors.ArgumentInvalid.With("botSessionTimeout", request.BotSessionTimeout, "must be greater than 0"))
	}
	if request.Message.GetType() == "" {
		merr.Append(errors.ArgumentMissing.With("inputMessage.type"))
	}

	// Only ButtonResponse is supported for Structured message in Genesys Cloud today
	if request.Message.GetType() == NormalizedMessageTypeStructured {
		for _, content := range request.Message.Content {
			if _, ok := content.(*NormalizedMessageButtonResponseContent); !ok {
				return errors.ArgumentInvalid.With("inputMessage.content", content.GetType(), "ButtonResponse")
			}
		}
	}

	return merr.AsError()
}

// MarshalJSON marshals to JSON
//
// implements core.Marshaller
func (request BotConnectorIncomingMessageRequest) MarshalJSON() ([]byte, error) {
	type surrogate BotConnectorIncomingMessageRequest

	data, err := json.Marshal(struct {
		surrogate
		BotSessionID   core.UUID `json:"botSessionId"`
		ConversationID core.UUID `json:"genesysConversationId"`
		MessageID      core.UUID `json:"messageId"`
	}{
		surrogate:      surrogate(request),
		BotSessionID:   core.UUID(request.BotSessionID),
		ConversationID: core.UUID(request.ConversationID),
		MessageID:      core.UUID(request.MessageID),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals from JSON
//
// implements core.Unmarshaller
func (request *BotConnectorIncomingMessageRequest) UnmarshalJSON(data []byte) error {
	type surrogate BotConnectorIncomingMessageRequest
	var inner struct {
		surrogate
		BotSessionID   core.UUID `json:"botSessionId"`
		ConversationID core.UUID `json:"genesysConversationId"`
		MessageID      core.UUID `json:"messageId"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*request = BotConnectorIncomingMessageRequest(inner.surrogate)
	request.BotSessionID = uuid.UUID(inner.BotSessionID)
	request.ConversationID = uuid.UUID(inner.ConversationID)
	request.MessageID = uuid.UUID(inner.MessageID)

	return nil
}

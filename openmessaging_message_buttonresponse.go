package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// OpenMessageButtonResponse is a message sent to the Open Messaging API
//
// See: https://developer.genesys.cloud/api/rest/v2/conversations/#post-api-v2-conversations-messages-inbound-open
type OpenMessageButtonResponse struct {
	ID                   string             `json:"id,omitempty"` // Can be anything
	Channel              OpenMessageChannel `json:"channel"`
	Direction            string             `json:"direction"` // inbound or outbound
	ButtonResponse       ButtonResponse     `json:"buttonResponse"`
	OriginatingMessageID string             `json:"originatingMessageId"`
	OriginatingEntity    string             `json:"originatingEntity,omitempty"` // Bot or Human
	Metadata             map[string]string  `json:"metadata,omitempty"`
	ConversationID       uuid.UUID          `json:"conversationId,omitempty"`
	KeysToRedact         []string           `json:"-"`
}

// init initializes this type
func init() {
	openMessageRegistry.Add(OpenMessageButtonResponse{})
}

// GetType tells the type of this OpenMessage
//
// implements core.TypeCarrier
func (message OpenMessageButtonResponse) GetType() string {
	return "ButtonResponse"
}

// GetID gets the identifier of this
//
//	implements OpenMessage
func (message OpenMessageButtonResponse) GetID() string {
	return message.ID
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageButtonResponse) Redact() interface{} {
	redacted := message
	redacted.Channel = message.Channel.Redact().(OpenMessageChannel)
	if core.GetEnvAsBool("REDACT_MESSAGE_TEXT", true) && len(message.ButtonResponse.Text) > 0 {
		redacted.ButtonResponse.Text = logger.RedactWithHash(message.ButtonResponse.Text)
	}
	if core.GetEnvAsBool("REDACT_MESSAGE_TEXT", true) && len(message.ButtonResponse.Payload) > 0 {
		redacted.ButtonResponse.Payload = logger.RedactWithHash(message.ButtonResponse.Payload)
	}
	for _, key := range message.KeysToRedact {
		if value, found := message.Metadata[key]; found {
			redacted.Metadata[key] = logger.RedactWithHash(value)
		}
	}
	return redacted
}

// MarshalJSON marshals this into JSON
func (message OpenMessageButtonResponse) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageButtonResponse
	data, err := json.Marshal(struct {
		surrogate
		Type string `json:"type"`
	}{
		surrogate: surrogate(message),
		Type:      message.GetType(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (message *OpenMessageButtonResponse) UnmarshalJSON(data []byte) (err error) {
	type surrogate OpenMessageButtonResponse
	var inner struct {
		surrogate
		Type           string    `json:"type"`
		KeysToRedact   []string  `json:"keysToRedact"`
		ConversationID core.UUID `json:"conversationId,omitempty"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageButtonResponse{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageButtonResponse(inner.surrogate)
	message.ConversationID = uuid.UUID(inner.ConversationID)
	message.KeysToRedact = append(message.KeysToRedact, inner.KeysToRedact...)
	return
}

package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// OpenMessageText is a text message sent or received by the Open Messaging API
//
// See: https://developer.genesys.cloud/api/rest/v2/conversations/#post-api-v2-conversations-messages-inbound-open
type OpenMessageText struct {
	ID                string               `json:"id,omitempty"` // Can be anything
	Channel           OpenMessageChannel   `json:"channel"`
	Direction         string               `json:"direction"` // inbound or outbound
	Text              string               `json:"text"`
	Content           []OpenMessageContent `json:"content,omitempty"`
	OriginatingEntity string               `json:"originatingEntity,omitempty"` // Bot or Human
	Metadata          map[string]string    `json:"metadata,omitempty"`
	ConversationID    uuid.UUID            `json:"conversationId,omitempty"`
	KeysToRedact      []string             `json:"-"`
}

// init initializes this type
func init() {
	openMessageRegistry.Add(OpenMessageText{})
}

// GetType tells the type of this OpenMessage
//
// implements core.TypeCarrier
func (message OpenMessageText) GetType() string {
	return "Text"
}

// GetID gets the identifier of this
//
//	implements OpenMessage
func (message OpenMessageText) GetID() string {
	return message.ID
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageText) Redact() interface{} {
	redacted := message
	redacted.Channel = message.Channel.Redact().(OpenMessageChannel)
	if core.GetEnvAsBool("REDACT_MESSAGE_TEXT", true) && len(message.Text) > 0 {
		redacted.Text = logger.RedactWithHash(message.Text)
	}
	redacted.Content = make([]OpenMessageContent, 0, len(message.Content))
	for _, content := range message.Content {
		if redactable, ok := content.(logger.Redactable); ok {
			redacted.Content = append(redacted.Content, redactable.Redact().(OpenMessageContent))
		}
	}
	for _, key := range message.KeysToRedact {
		if value, found := message.Metadata[key]; found {
			redacted.Metadata[key] = logger.RedactWithHash(value)
		}
	}
	return redacted
}

// MarshalJSON marshals this into JSON
func (message OpenMessageText) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageText
	data, err := json.Marshal(struct {
		surrogate
		Type           string    `json:"type"`
		ConversationID core.UUID `json:"conversationId,omitempty"`
	}{
		surrogate:      surrogate(message),
		Type:           message.GetType(),
		ConversationID: core.UUID(message.ConversationID),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (message *OpenMessageText) UnmarshalJSON(data []byte) (err error) {
	type surrogate OpenMessageText
	var inner struct {
		surrogate
		Type           string            `json:"type"`
		KeysToRedact   []string          `json:"keysToRedact"`
		Content        []json.RawMessage `json:"content,omitempty"`
		ConversationID core.UUID         `json:"conversationId,omitempty"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageText{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageText(inner.surrogate)
	message.ConversationID = uuid.UUID(inner.ConversationID)
	message.KeysToRedact = append(message.KeysToRedact, inner.KeysToRedact...)
	unmarshalMode := core.GetEnvAsString("JSON_UNMARSHAL_MODE", "strict") // "strict" or "ignore_unknown_keys"
	isUnmarshalIgnoreUnknownKeys := unmarshalMode == "ignore_unknown_keys"
	if len(inner.Content) > 0 {
		message.Content = make([]OpenMessageContent, 0, len(inner.Content))
		for _, content := range inner.Content {
			content, err := UnmarshalOpenMessageContent(content)
			if errors.Is(err, errors.InvalidType) && isUnmarshalIgnoreUnknownKeys {
				continue
			} else if errors.Is(err, errors.ArgumentMissing) {
				return err
			} else if errors.Is(err, errors.JSONUnmarshalError) {
				return err
			} else if err != nil {
				return errors.JSONUnmarshalError.Wrap(err)
			}
			message.Content = append(message.Content, content)
		}
	}
	return
}

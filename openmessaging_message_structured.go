package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

// OpenMessageStructed describes a structured message
//
// See: https://developer.genesys.cloud/api/rest/v2/conversations/#post-api-v2-conversations-messages-inbound-open
type OpenMessageStructured struct {
	ID                string               `json:"id,omitempty"` // Can be anything
	Channel           OpenMessageChannel   `json:"channel"`
	Direction         string               `json:"direction"` // inbound or outbound
	Text              string               `json:"text"`
	Content           []OpenMessageContent `json:"content,omitempty"`
	OriginatingEntity string               `json:"originatingEntity,omitempty"` // Bot or Human
	Metadata          map[string]string    `json:"metadata,omitempty"`
	KeysToRedact      []string             `json:"-"`
}

// init initializes this type
func init() {
	openMessageRegistry.Add(OpenMessageStructured{})
}

// GetType tells the type of this OpenMessage
//
// implements core.TypeCarrier
func (message OpenMessageStructured) GetType() string {
	return "Structured"
}

// GetID gets the identifier of this
//
//	implements OpenMessage
func (message OpenMessageStructured) GetID() string {
	return message.ID
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageStructured) Redact() interface{} {
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
		if value, found := redacted.Metadata[key]; found {
			redacted.Metadata[key] = logger.RedactWithHash(value)
		}
	}
	return redacted
}

// MarshalJSON marshals this into JSON
func (message OpenMessageStructured) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageStructured
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
func (message *OpenMessageStructured) UnmarshalJSON(data []byte) (err error) {
	type surrogate OpenMessageStructured
	var inner struct {
		surrogate
		Type         string            `json:"type"`
		KeysToRedact []string          `json:"keysToRedact"`
		Content      []json.RawMessage `json:"content,omitempty"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageStructured{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageStructured(inner.surrogate)
	message.KeysToRedact = append(message.KeysToRedact, inner.KeysToRedact...)
	if len(inner.Content) > 0 {
		message.Content = make([]OpenMessageContent, len(inner.Content))
		for i, content := range inner.Content {
			if message.Content[i], err = UnmarshalOpenMessageContent(content); errors.Is(err, errors.JSONUnmarshalError) {
				return err
			} else if err != nil {
				return errors.JSONUnmarshalError.Wrap(err)
			}
		}
	}
	return
}

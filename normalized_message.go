package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// TODO: This will need to go to go-gcloudcx

// NormalizedMessage represents a Genesys digital message for a Genesys Cloud bot
//
// It can be either a text message or a structured message with content
// The Text is mandatory for text messages, while Content is mandatory for structured messages
// The Content can only contain ButtonResponse content type
type NormalizedMessage struct {
	Type    string                     `json:"type"` // Text or Structured
	Text    string                     `json:"text,omitempty"`
	Content []NormalizedMessageContent `json:"content,omitempty"`
}

const (
	NormalizedMessageTypeText       = "Text"       // NormalizedMessageTypeText represents a text input message
	NormalizedMessageTypeStructured = "Structured" // NormalizedMessageTypeStructured represents a structured input message
)

// GetType returns the type of the input message
//
// implements core.TypeCarrier
func (message NormalizedMessage) GetType() string {
	return message.Type
}

// Validate validates the input message
func (message *NormalizedMessage) Validate() error {
	var merr errors.MultiError

	if message.Type == "" {
		merr.Append(errors.ArgumentMissing.With("inputMessage.type"))
	} else if message.Type != NormalizedMessageTypeText && message.Type != NormalizedMessageTypeStructured {
		merr.Append(errors.ArgumentInvalid.With("inputMessage.type", message.Type, "Text, Structured"))
	}

	if message.Type == NormalizedMessageTypeText && len(message.Text) == 0 {
		merr.Append(errors.ArgumentMissing.With("inputMessage.text"))
	}

	if message.Type == NormalizedMessageTypeStructured && len(message.Content) == 0 {
		merr.Append(errors.ArgumentMissing.With("inputMessage.content"))
	}

	return merr.AsError()
}

// MarshalJSON marshals to JSON
//
// implements core.Marshaller
func (message NormalizedMessage) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessage

	data, err := json.Marshal(struct {
		surrogate
	}{
		surrogate: surrogate(message),
	})
	if err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}

	return data, nil
}

// UnmarshalJSON unmarshals from JSON
//
// implements core.Unmarshaller
func (message *NormalizedMessage) UnmarshalJSON(data []byte) error {
	type surrogate NormalizedMessage
	var inner struct {
		surrogate
		Content []json.RawMessage `json:"content"`
	}

	if err := errors.JSONUnmarshalError.Wrap(json.Unmarshal(data, &inner)); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}

	*message = NormalizedMessage(inner.surrogate)

	if len(inner.Content) > 0 {
		message.Content = make([]NormalizedMessageContent, 0, len(inner.Content))

		for _, raw := range inner.Content {
			content, err := UnmarshalOpenMessageContent(raw)
			if err != nil {
				return errors.JSONUnmarshalError.Wrap(err)
			}
			message.Content = append(message.Content, content)
		}
	}

	return message.Validate()
}

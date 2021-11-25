package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

type OpenMessageText struct {
	ID              string                `json:"id,omitempty"` // Can be anything
	Channel         *OpenMessageChannel   `json:"channel"`
	Direction       string                `json:"direction"`
	Text            string                `json:"text"`
	Content         []*OpenMessageContent `json:"content,omitempty"`
	RelatedMessages []*OpenMessageText    `json:"relatedMessages,omitempty"`
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

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageText) Redact() interface{} {
	redacted := message
	if message.Channel != nil {
		redacted.Channel = message.Channel.Redact().(*OpenMessageChannel)
	}
	return &redacted
}

// MarshalJSON marshals this into JSON
func (message OpenMessageText) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageText
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
func (message *OpenMessageText) UnmarshalJSON(data []byte) (err error) {
	type surrogate OpenMessageText
	var inner struct {
		surrogate
		Type string `json:"type"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageText{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageText(inner.surrogate)
	return
}

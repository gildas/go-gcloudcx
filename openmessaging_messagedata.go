package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

type OpenMessageData struct {
	ID                string       `json:"id,omitempty"` // Can be anything
	Name              string       `json:"name,omitempty"`
	ProviderMessageID string       `json:"providerMessageId,omitempty"`
	Timestamp         time.Time    `json:"-"`
	From              string       `json:"fromAddress,omitempty"`
	To                string       `json:"toAddress,omitempty"`
	Direction         string       `json:"direction"`     // inbound or outbound
	MessengerType     string       `json:"messengerType"` // sms, facebook, twitter, etc
	Text              string       `json:"textBody"`
	NormalizedMessage OpenMessage  `json:"-"`
	Status            string       `json:"status"`              // sent, received, delivered, undelivered, etc
	CreatedBy         *User        `json:"createdBy,omitempty"` // nil unless NormalizedMessage.OriginatingEntity is "Human"
	Conversation      Conversation `json:"-"`
	SelfURI           URI          `json:"selfUri"`
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (messageData OpenMessageData) Redact() interface{} {
	redacted := messageData
	if messageData.CreatedBy != nil {
		redactedUser := messageData.CreatedBy.Redact().(User)
		redacted.CreatedBy = &redactedUser
	}
	if messageData.NormalizedMessage != nil {
		redacted.NormalizedMessage = messageData.NormalizedMessage.Redact().(OpenMessage)
	}
	if core.GetEnvAsBool("REDACT_MESSAGE_TEXT", true) && len(messageData.Text) > 0 {
		redacted.Text = logger.RedactWithHash(messageData.Text)
	}
	return redacted
}

// MarshalJSON marshals this into JSON
func (messageData OpenMessageData) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageData
	data, err := json.Marshal(struct {
		surrogate
		Timestamp         core.Time   `json:"timestamp"`
		NormalizedMessage OpenMessage `json:"normalizedMessage"`
		ConversationID    uuid.UUID   `json:"conversationId"`
	}{
		surrogate:         surrogate(messageData),
		Timestamp:         core.Time(messageData.Timestamp),
		NormalizedMessage: messageData.NormalizedMessage,
		ConversationID:    messageData.Conversation.ID,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (messageData *OpenMessageData) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageData
	var inner struct {
		surrogate
		Timestamp         core.Time       `json:"timestamp"`
		NormalizedMessage json.RawMessage `json:"normalizedMessage"`
		ConversationID    uuid.UUID       `json:"conversationId"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*messageData = OpenMessageData(inner.surrogate)
	messageData.Timestamp = inner.Timestamp.AsTime()
	messageData.Conversation = Conversation{ID: inner.ConversationID, SelfURI: NewURI("/conversations/%s", inner.ConversationID)}
	messageData.NormalizedMessage, err = UnmarshalOpenMessage(inner.NormalizedMessage)
	return
}

package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type OpenMessageChannel struct {
	ID        uuid.UUID                   `json:"id,omitempty"`
	Platform  string                      `json:"platform"` // Open
	Type      string                      `json:"type"` // Private, Public
	MessageID string                      `json:"messageId"`
	Time      time.Time                   `json:"-"`
	To        *OpenMessageTo              `json:"to"`
	From      *OpenMessageFrom            `json:"from"`
	Metadata  *OpenMessageChannelMetadata `json:"metadata,omitempty"`
}

type OpenMessageChannelMetadata struct {
	Attributes map[string]string `json:"customAttributes,omitempty"`
}

func NewOpenMessageChannel(messageID string, to *OpenMessageTo, from *OpenMessageFrom) *OpenMessageChannel {
	return &OpenMessageChannel{
		Platform:  "Open",
		Type:      "Private",
		MessageID: messageID,
		Time:      time.Now().UTC(),
		To:        to,
		From:      from,
	}
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (channel OpenMessageChannel) Redact() interface{} {
	redacted := channel
	if channel.From != nil {
		redacted.From = channel.From.Redact().(*OpenMessageFrom)
	}
	return &redacted
}

// MarshalJSON marshals this into JSON
func (channel OpenMessageChannel) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageChannel
	var id string
	if channel.ID != uuid.Nil {
		id = channel.ID.String()
	}
	data, err := json.Marshal(struct{
		ID        string `json:"id,omitempty"`
		surrogate
		Time core.Time `json:"time"`
	}{
		ID:        id,
		surrogate: surrogate(channel),
		Time:      (core.Time)(channel.Time),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (channel *OpenMessageChannel) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageChannel
	var inner struct {
		ID        string `json:"id"`
		surrogate
		Time core.Time `json:"time"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*channel = OpenMessageChannel(inner.surrogate)
	channel.Time = inner.Time.AsTime()
	if len(inner.ID) > 0 {
		channel.ID, err = uuid.Parse(inner.ID)
	}
	return
}

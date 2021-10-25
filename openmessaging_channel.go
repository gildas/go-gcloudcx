package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageChannel struct {
	Platform  string           `json:"platform"` // Open
	Type      string           `json:"type"` // Private, Public
	MessageID string           `json:"messageId"`
	Time      time.Time        `json:"-"`
	To        *OpenMessageTo   `json:"to"`
	From      *OpenMessageFrom `json:"from"`
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
	data, err := json.Marshal(struct{
		surrogate
		T core.Time `json:"time"`
	}{
		surrogate: surrogate(channel),
		T:         (core.Time)(channel.Time),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (channel *OpenMessageChannel) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageChannel
	var inner struct {
		surrogate
		T core.Time `json:"time"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*channel = OpenMessageChannel(inner.surrogate)
	channel.Time = inner.T.AsTime()
	return
}
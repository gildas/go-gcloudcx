package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

type OpenMessageChannel struct {
	ID               uuid.UUID         `json:"id,omitempty"`
	Platform         string            `json:"platform"` // Open
	Type             string            `json:"type"`     // Private, Public
	MessageID        string            `json:"messageId,omitempty"`
	Time             time.Time         `json:"-"`
	To               *OpenMessageTo    `json:"to,omitempty"`
	From             *OpenMessageFrom  `json:"from"`
	CustomAttributes map[string]string `json:"-"`
	KeysToRedact     []string          `json:"-"`
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (channel OpenMessageChannel) Redact() interface{} {
	redacted := channel
	if channel.From != nil {
		redactedFrom := channel.From.Redact().(OpenMessageFrom)
		redacted.From = &redactedFrom
	}
	if len(channel.KeysToRedact) > 0 {
		redacted.CustomAttributes = make(map[string]string, len(channel.CustomAttributes))
		for key, value := range channel.CustomAttributes {
			if core.Contains(channel.KeysToRedact, key) {
				redacted.CustomAttributes[key] = logger.RedactWithHash(value)
			} else {
				redacted.CustomAttributes[key] = value
			}
		}
	}
	return redacted
}

// Validate checks if the object is valid
func (channel *OpenMessageChannel) Validate() (err error) {
	if channel.Platform != "Open" {
		return errors.ArgumentInvalid.With("channel.platform", channel.Platform)
	}
	if channel.Type != "Private" && channel.Type != "Public" {
		return errors.ArgumentInvalid.With("channel.type", channel.Type)
	}
	if channel.From == nil {
		return errors.ArgumentMissing.With("channel.from")
	}
	if err = channel.From.Validate(); err != nil {
		return
	}
	if channel.To == nil {
		return errors.ArgumentMissing.With("To")
	}
	if err = channel.To.Validate(); err != nil {
		return
	}
	// TODO: Not that simple.... mandatory fields depend on the type of OpenMessage the channel belongs to...
	return
}

// MarshalJSON marshals this into JSON
func (channel OpenMessageChannel) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageChannel
	type OpenMessageChannelMetadata struct {
		Attributes map[string]string `json:"customAttributes,omitempty"`
	}
	var id string
	var metadata *OpenMessageChannelMetadata

	if channel.ID != uuid.Nil {
		id = channel.ID.String()
	}
	if len(channel.CustomAttributes) > 0 {
		metadata = &OpenMessageChannelMetadata{
			Attributes: channel.CustomAttributes,
		}
	}
	data, err := json.Marshal(struct {
		ID string `json:"id,omitempty"`
		surrogate
		Time     core.Time                   `json:"time"`
		Metadata *OpenMessageChannelMetadata `json:"metadata,omitempty"`
	}{
		ID:        id,
		surrogate: surrogate(channel),
		Time:      (core.Time)(channel.Time),
		Metadata:  metadata,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (channel *OpenMessageChannel) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageChannel
	type OpenMessageChannelMetadata struct {
		Attributes map[string]string `json:"customAttributes,omitempty"`
	}
	var inner struct {
		ID string `json:"id"`
		surrogate
		Time         core.Time                   `json:"time"`
		Metadata     *OpenMessageChannelMetadata `json:"metadata,omitempty"`
		KeysToRedact []string                    `json:"keysToRedact"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*channel = OpenMessageChannel(inner.surrogate)
	channel.Time = inner.Time.AsTime()
	if len(inner.ID) > 0 {
		channel.ID, err = uuid.Parse(inner.ID)
	}
	if inner.Metadata != nil {
		channel.CustomAttributes = inner.Metadata.Attributes
	}
	channel.KeysToRedact = append(channel.KeysToRedact, inner.KeysToRedact...)
	return
}

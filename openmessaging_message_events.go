package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

// OpenMessageText is a text message sent or received by the Open Messaging API
//
// See https://developer.genesys.cloud/commdigital/digital/openmessaging/inboundEventMessages
type OpenMessageEvents struct {
	ID           string             `json:"id,omitempty"` // Can be anything
	Channel      OpenMessageChannel `json:"channel"`
	Direction    string             `json:"direction,omitempty"` // Can be "Inbound" or "Outbound"
	Events       []OpenMessageEvent `json:"events"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
	KeysToRedact []string           `json:"-"`
}

// init initializes this type
func init() {
	openMessageRegistry.Add(OpenMessageEvents{})
}

// GetType returns the type of this event
//
// implements core.TypeCarrier
func (message OpenMessageEvents) GetType() string {
	return "Event"
}

// GetID gets the identifier of this
//
//	implements OpenMessage
func (message OpenMessageEvents) GetID() string {
	return message.ID
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageEvents) Redact() interface{} {
	redacted := message
	redacted.Channel = message.Channel.Redact().(OpenMessageChannel)
	for i, event := range message.Events {
		if redactable, ok := event.(logger.Redactable); ok {
			message.Events[i] = redactable.Redact().(OpenMessageEvent)
		}
	}
	for _, key := range message.KeysToRedact {
		if value, found := redacted.Metadata[key]; found {
			redacted.Metadata[key] = logger.RedactWithHash(value)
		}
	}
	return &redacted
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (message OpenMessageEvents) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageEvents

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
//
// implements json.Unmarshaler
func (message *OpenMessageEvents) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageEvents
	var inner struct {
		surrogate
		Type         string            `json:"type"`
		Events       []json.RawMessage `json:"events"`
		KeysToRedact []string          `json:"keysToRedact"`
	}
	if err = json.Unmarshal(payload, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*message = OpenMessageEvents(inner.surrogate)

	message.Events = make([]OpenMessageEvent, 0, len(inner.Events))
	for _, raw := range inner.Events {
		event, err := UnmarshalOpenMessageEvent(raw)
		if err != nil {
			return err
		}
		message.Events = append(message.Events, event)
	}
	message.KeysToRedact = append(message.KeysToRedact, inner.KeysToRedact...)
	return
}

package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageText is a text message sent or received by the Open Messaging API
//
// See https://developer.genesys.cloud/commdigital/digital/openmessaging/inboundEventMessages
type OpenMessageEvents struct {
	ID        string              `json:"id,omitempty"` // Can be anything
	Channel   *OpenMessageChannel `json:"channel"`
	Direction string              `json:"direction,omitempty"` // Can be "Inbound" or "Outbound"
	Events    []OpenMessageEvent  `json:"events"`
	Metadata  map[string]string   `json:"metadata,omitempty"`
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
		Type   string            `json:"type"`
		Events []json.RawMessage `json:"events"`
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
	return
}

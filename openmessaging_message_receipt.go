package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageReceipt is the message receipt returned by the Open Message Integration API.
//
// See: https://developer.genesys.cloud/api/digital/openmessaging/receipts
type OpenMessageReceipt struct {
	ID           string              `json:"id,omitempty"`      // Can be anything, message ID this receipt relates to
	Channel      *OpenMessageChannel `json:"channel"`
	Direction    string              `json:"direction"`         // Can be "Inbound" or "Outbound"
	Status       string              `json:"status"`            // Can be "Published" (Inbound), "Delivered" (Outbound), "Failed"
	Reasons      []*StatusReason     `json:"reasons,omitempty"` // Contains the reason for the failure
	FinalReceipt bool                `json:"isFinalReceipt"`    // True if this is the last receipt about this message ID
	Metadata     map[string]string   `json:"metadata,omitempty"`
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// init initializes this type
func init() {
	openMessageRegistry.Add(OpenMessageReceipt{})
}

// GetType tells the type of this OpenMessage
//
// implements core.TypeCarrier
func (message OpenMessageReceipt) GetType() string {
	return "Receipt"
}

// GetID gets the identifier of this
//   implements OpenMessage
func (message OpenMessageReceipt) GetID() string {
	return message.ID
}

// IsFailed tells if the receipt is failed
func (message OpenMessageReceipt) IsFailed() bool {
	return message.Status == "Failed"
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageReceipt) Redact() interface{} {
	redacted := message
	if message.Channel != nil {
		redacted.Channel = message.Channel.Redact().(*OpenMessageChannel)
	}
	return &redacted
}

// MarshalJSON marshals this into JSON
func (message OpenMessageReceipt) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageReceipt
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
func (message *OpenMessageReceipt) UnmarshalJSON(data []byte) (err error) {
	type surrogate OpenMessageReceipt
	var inner struct {
		surrogate
		Type string `json:"type"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageReceipt{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageReceipt(inner.surrogate)
	return
}

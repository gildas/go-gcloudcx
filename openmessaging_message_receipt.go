package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageReceipt is the message receipt returned by the Open Message Integration API.
//
// See: https://developer.genesys.cloud/api/digital/openmessaging/receipts
type OpenMessageReceipt struct {
	ID          string              `json:"id,omitempty"` // Can be anything
	Channel     *OpenMessageChannel `json:"channel"`
	Direction   string              `json:"direction"`
	Status      string              `json:"status"`
	Reasons     []*StatusReason     `json:"reasons,omitempty"`
	FinalReceipt bool               `json:"isFinalReceipt"`
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

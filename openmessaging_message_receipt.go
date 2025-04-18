package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// OpenMessageReceipt is the message receipt returned by the Open Message Integration API.
//
// See: https://developer.genesys.cloud/api/digital/openmessaging/receipts
type OpenMessageReceipt struct {
	ID             string             `json:"id,omitempty"` // Can be anything, message ID this receipt relates to
	Channel        OpenMessageChannel `json:"channel"`
	Direction      string             `json:"direction"`         // Can be "Inbound" or "Outbound"
	Status         string             `json:"status"`            // Can be "Published" (Inbound), "Delivered" (Outbound), "Sent", "Read", "Failed", "Removed"
	Reasons        []StatusReason     `json:"reasons,omitempty"` // Contains the reason for the failure
	FinalReceipt   bool               `json:"isFinalReceipt"`    // True if this is the last receipt about this message ID
	Metadata       map[string]string  `json:"metadata,omitempty"`
	ConversationID uuid.UUID          `json:"conversationId,omitempty"`
	KeysToRedact   []string           `json:"-"`
}

type StatusReason struct {
	Code    string `json:"code,omitempty"` // MessageExpired, RateLimited, MessageNotAllowed, GeneralError, UnsupportedMessage, UnknownMessage, InvalidMessageStructure, InvalidDestination, ServerError, MediaTypeNotAllowed, InvalidMediaContentLength, RecipientOptedOut
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
//
//	implements OpenMessage
func (message OpenMessageReceipt) GetID() string {
	return message.ID
}

// IsFailed tells if the receipt is failed
func (message OpenMessageReceipt) IsFailed() bool {
	return message.Status == "Failed"
}

// AsError converts this to an error
func (reason StatusReason) AsError() error {
	switch reason.Code {
	case "MessageExpired":
		return MessageExpired.With(reason.Message)
	case "RateLimited":
		return RateLimited.With(reason.Message)
	case "MessageNotAllowed":
		return MessageNotAllowed.With(reason.Message)
	case "GeneralError":
		return GeneralError.With(reason.Message)
	case "UnsupportedMessage":
		return UnsupportedMessage.With(reason.Message)
	case "UnknownMessage":
		return UnknownMessage.With(reason.Message)
	case "InvalidMessageStructure":
		return InvalidMessageStructure.With(reason.Message)
	case "InvalidDestination":
		return InvalidDestination.With(reason.Message)
	case "ServerError":
		return ServerError.With(reason.Message)
	case "MediaTypeNotAllowed":
		return MediaTypeNotAllowed.With(reason.Message)
	case "InvalidMediaContentLength":
		return InvalidMediaContentLength.With(reason.Message)
	case "RecipientOptedOut":
		return RecipientOptedOut.With(reason.Message)
	}
	return GeneralError.With(reason.Message)
}

// AsError converts this to an error
func (message OpenMessageReceipt) AsError() error {
	if !message.IsFailed() {
		return nil
	}
	return errors.Join(core.Map(message.Reasons, func(reason StatusReason) error { return reason.AsError() })...)
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessageReceipt) Redact() interface{} {
	redacted := message
	redacted.Channel = message.Channel.Redact().(OpenMessageChannel)
	for _, key := range message.KeysToRedact {
		if value, found := redacted.Metadata[key]; found {
			redacted.Metadata[key] = logger.RedactWithHash(value)
		}
	}
	return redacted
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
		Type           string    `json:"type"`
		KeysToRedact   []string  `json:"keysToRedact"`
		ConversationID core.UUID `json:"conversationId,omitempty"`
	}

	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != (OpenMessageReceipt{}).GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type))
	}
	*message = OpenMessageReceipt(inner.surrogate)
	message.ConversationID = uuid.UUID(inner.ConversationID)
	message.KeysToRedact = append(message.KeysToRedact, inner.KeysToRedact...)
	return
}

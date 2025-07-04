package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type BotConnectorOutgoingMessageRequest struct {
	BotID         string                 `json:"botId"`
	BotVersion    string                 `json:"botVersion"`
	BotState      string                 `json:"botState"` // Complete, Failed, MoreData
	BotSessionID  uuid.UUID              `json:"botSessionId"`
	MessageID     uuid.UUID              `json:"messageId"`
	Language      string                 `json:"languageCode"`
	Intent        string                 `json:"intent,omitempty"`        // The discovered intent, mandatory if BotState is "Complete"
	Confidence    float64                `json:"confidence,omitempty"`    // Confidence score for the intent, optional
	Entities      []SlotEntity           `json:"entities,omitempty"`      // List of slot entities, optional
	ReplyMessages []NormalizedMessage    `json:"replyMessages,omitempty"` // List of reply messages to be sent back, optional
	Parameters    map[string]string      `json:"parameters,omitempty"`    // Message Parameters, optional
	ErrorInfo     *BotConnectorErrorInfo `json:"errorInfo,omitempty"`     // Error information if BotState is "Failed", optional
	CorrelationID string                 `json:"correlationId,omitempty"` // Optional correlation ID for tracking purposes
}

// Validate validates the outgoing message request
func (request *BotConnectorOutgoingMessageRequest) Validate() error {
	var merr errors.MultiError

	if len(request.BotID) == 0 {
		merr.Append(errors.ArgumentMissing.With("botId"))
	}
	if request.BotSessionID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("botSessionId"))
	}
	if request.MessageID == uuid.Nil {
		merr.Append(errors.ArgumentMissing.With("messageId"))
	}
	if request.Language == "" {
		merr.Append(errors.ArgumentMissing.With("languageCode"))
	}
	if request.BotState == "" {
		merr.Append(errors.ArgumentMissing.With("botState"))
	} else if request.BotState != BotStateComplete && request.BotState != BotStateFailed && request.BotState != BotStateMoreData {
		merr.Append(errors.ArgumentInvalid.With("botState", request.BotState, "Complete, Failed, MoreData"))
	}

	if request.BotState == "" {
		merr.Append(errors.ArgumentMissing.With("botState"))
	} else if request.BotState != BotStateComplete && request.BotState != BotStateFailed && request.BotState != BotStateMoreData {
		merr.Append(errors.ArgumentInvalid.With("botState", request.BotState, "Complete, Failed, MoreData"))
	}

	if request.BotState == "Complete" && request.Intent == "" {
		merr.Append(errors.ArgumentMissing.With("intent"))
	}

	return merr.AsError()
}

// MarshalJSON marshals to JSON
//
// implements core.Marshaller
func (request BotConnectorOutgoingMessageRequest) MarshalJSON() ([]byte, error) {
	type surrogate BotConnectorOutgoingMessageRequest

	var confidence *float64
	if request.Confidence != 0 {
		confidence = &request.Confidence
	}

	data, err := json.Marshal(struct {
		surrogate
		BotSessionID core.UUID `json:"botSessionId"`
		MessageID    core.UUID `json:"messageId"`
		Confidence   *float64  `json:"confidence,omitempty"` // Confidence is optional, so we use a pointer to allow nil
	}{
		surrogate:    surrogate(request),
		BotSessionID: core.UUID(request.BotSessionID),
		MessageID:    core.UUID(request.MessageID),
		Confidence:   confidence,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals from JSON
//
// implements core.Unmarshaller
func (request *BotConnectorOutgoingMessageRequest) UnmarshalJSON(data []byte) error {
	type surrogate BotConnectorOutgoingMessageRequest
	var inner struct {
		surrogate
		BotSessionID core.UUID         `json:"botSessionId"`
		MessageID    core.UUID         `json:"messageId"`
		Entities     []json.RawMessage `json:"entities"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*request = BotConnectorOutgoingMessageRequest(inner.surrogate)
	request.BotSessionID = uuid.UUID(inner.BotSessionID)
	request.MessageID = uuid.UUID(inner.MessageID)
	if len(inner.Entities) > 0 {
		request.Entities = make([]SlotEntity, 0, len(inner.Entities))
		for _, raw := range inner.Entities {
			entity, err := UnmarshalSlotEntity(raw)
			if err != nil {
				return errors.JSONUnmarshalError.Wrap(err)
			}
			request.Entities = append(request.Entities, entity)
		}
	}
	if request.Parameters == nil {
		request.Parameters = make(map[string]string)
	}
	return nil
}

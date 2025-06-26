package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// TODO: This will need to go to go-gcloudcx

type BotConnectorIncomingMessageResponse struct {
	BotState      string                 `json:"botState"`                // Complete, Failed, MoreData
	Intent        string                 `json:"intent,omitempty"`        // The discovered intent, mandatory if BotState is "Complete"
	Confidence    float64                `json:"confidence,omitempty"`    // Confidence score for the intent, optional
	Entities      []SlotEntity           `json:"entities,omitempty"`      // List of slot entities, optional
	ReplyMessages []NormalizedMessage    `json:"replyMessages,omitempty"` // List of reply messages to be sent back, optional
	Parameters    map[string]string      `json:"parameters,omitempty"`    // Message Parameters, optional
	ErrorInfo     *BotConnectorErrorInfo `json:"errorInfo,omitempty"`     // Error information if BotState is "Failed", optional
}

const (
	BotStateComplete = "Complete" // BotState when the message processing is complete
	BotStateFailed   = "Failed"   // BotState when the message processing has failed
	BotStateMoreData = "MoreData" // BotState when the message processing requires more data
)

// Validate validates the response
func (response *BotConnectorIncomingMessageResponse) Validate() error {
	var merr errors.MultiError

	if response.BotState == "" {
		merr.Append(errors.ArgumentMissing.With("botState"))
	} else if response.BotState != BotStateComplete && response.BotState != BotStateFailed && response.BotState != BotStateMoreData {
		merr.Append(errors.ArgumentInvalid.With("botState", response.BotState, "Complete, Failed, MoreData"))
	}

	if response.BotState == "Complete" && response.Intent == "" {
		merr.Append(errors.ArgumentMissing.With("intent"))
	}

	return merr.AsError()
}

// MarshalJSON marshals to JSON
//
// implements core.Marshaller
func (response BotConnectorIncomingMessageResponse) MarshalJSON() ([]byte, error) {
	type surrogate BotConnectorIncomingMessageResponse

	if err := response.Validate(); err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}

	var confidence *float64
	if response.Confidence != 0 {
		confidence = &response.Confidence
	}

	data, err := json.Marshal(struct {
		surrogate
		Confidence *float64 `json:"confidence,omitempty"` // Confidence is optional, so we use a pointer to allow nil
	}{
		surrogate:  surrogate(response),
		Confidence: confidence,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals from JSON
//
// implements core.Unmarshaller
func (response *BotConnectorIncomingMessageResponse) UnmarshalJSON(data []byte) error {
	type surrogate BotConnectorIncomingMessageResponse
	var inner struct {
		surrogate
		Entities []json.RawMessage `json:"entities"`
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*response = BotConnectorIncomingMessageResponse(inner.surrogate)

	if len(inner.Entities) > 0 {
		response.Entities = make([]SlotEntity, 0, len(inner.Entities))
		for _, raw := range inner.Entities {
			entity, err := UnmarshalSlotEntity(raw)
			if err != nil {
				return errors.JSONUnmarshalError.Wrap(err)
			}
			response.Entities = append(response.Entities, entity)
		}
	}
	if response.Parameters == nil {
		response.Parameters = make(map[string]string)
	}

	return response.Validate()
}

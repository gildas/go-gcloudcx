package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type IntegrationState struct {
	Code          string      `json:"code"`             // ACTIVE, ACTIVATING, INACTIVE, DEACTIVATING, ERROR
	Effective     string      `json:"effective"`        // Human readable state
	Detail        MessageInfo `json:"detail,omitempty"` // Details about the state
	LastUpdatedAt time.Time   `json:"lastUpdated"`      // Last time the state was updated

}

// MessageInfo represents a message with additional information
type MessageInfo struct {
	LocalizableMessageCode string            `json:"localizableMessageCode,omitempty"` // The code for the localizable message
	Message                string            `json:"message,omitempty"`                // The message text
	MessageWithParams      string            `json:"messageWithParams,omitempty"`      // The message with parameters
	MessageParams          map[string]string `json:"messageParams,omitempty"`          // Parameters for the message
}

// MarshalJSON marshals to JSON
func (state IntegrationState) MarshalJSON() ([]byte, error) {
	type surrogate IntegrationState
	data, err := json.Marshal(struct {
		surrogate
		LastUpdatedAt core.Time `json:"lastUpdated"` // Format time as string
	}{
		surrogate:     surrogate(state),
		LastUpdatedAt: core.Time(state.LastUpdatedAt),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UnmarshalJSON unmarshals from JSON
func (state *IntegrationState) UnmarshalJSON(payload []byte) error {
	type surrogate IntegrationState
	var inner struct {
		surrogate
		LastUpdatedAt core.Time `json:"lastUpdated"` // Format time as string
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}

	*state = IntegrationState(inner.surrogate)
	state.LastUpdatedAt = time.Time(inner.LastUpdatedAt)

	return nil
}

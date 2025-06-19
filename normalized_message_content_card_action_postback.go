package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// NormalizedMessageCardPostbackAction describes a Postback action of a Card
type NormalizedMessageCardPostbackAction struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

func init() {
	cardActionTypeRegistry.Add(NormalizedMessageCardPostbackAction{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (action NormalizedMessageCardPostbackAction) GetType() string {
	return "Postback"
}

// String returns a string representation of this action
//
// implements fmt.Stringer
func (action NormalizedMessageCardPostbackAction) String() string {
	return action.Text
}

// Validate checks that this action is valid
func (action *NormalizedMessageCardPostbackAction) Validate() error {
	var merr errors.MultiError

	if action.Text == "" {
		merr.Append(errors.ArgumentMissing.With("text"))
	}

	return merr.AsError()
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (action NormalizedMessageCardPostbackAction) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCardPostbackAction

	data, err := json.Marshal(struct {
		ActionType string `json:"type"`
		surrogate
	}{
		ActionType: action.GetType(),
		surrogate:  surrogate(action),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (action *NormalizedMessageCardPostbackAction) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCardPostbackAction

	var inner struct {
		ActionType string `json:"actionType"`
		surrogate
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = NormalizedMessageCardPostbackAction(inner.surrogate)
	return errors.JSONUnmarshalError.Wrap(action.Validate())
}

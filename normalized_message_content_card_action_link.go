package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageCardLinkAction describes a Link action of a Card
type NormalizedMessageCardLinkAction struct {
	Text string   `json:"text"`
	URL  *url.URL `json:"url"`
}

func init() {
	cardActionTypeRegistry.Add(NormalizedMessageCardLinkAction{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (action NormalizedMessageCardLinkAction) GetType() string {
	return "Link"
}

// String returns a string representation of this action
//
// implements fmt.Stringer
func (action NormalizedMessageCardLinkAction) String() string {
	return action.Text
}

// Validate checks that this action is valid
func (action *NormalizedMessageCardLinkAction) Validate() error {
	var merr errors.MultiError

	if action.Text == "" {
		merr.Append(errors.ArgumentMissing.With("text"))
	}

	if action.URL == nil {
		merr.Append(errors.ArgumentMissing.With("url"))
	}

	return merr.AsError()
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (action NormalizedMessageCardLinkAction) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCardLinkAction

	if action.URL == nil {
		return nil, errors.JSONMarshalError.Wrap(errors.ArgumentMissing.With("url"))
	}

	data, err := json.Marshal(struct {
		ActionType string `json:"type"`
		surrogate
		URL *core.URL `json:"url"`
	}{
		ActionType: action.GetType(),
		surrogate:  surrogate(action),
		URL:        (*core.URL)(action.URL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (action *NormalizedMessageCardLinkAction) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCardLinkAction

	var inner struct {
		ActionType string `json:"actionType"`
		surrogate
		URL *core.URL `json:"url"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*action = NormalizedMessageCardLinkAction(inner.surrogate)
	action.URL = (*url.URL)(inner.URL)

	return errors.JSONUnmarshalError.Wrap(action.Validate())
}

package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

type NormalizedMessageButtonResponseContent struct {
	ButtonType string `json:"type,omitempty"` // "Button", "QuickReply"
	Text       string `json:"text"`
	Payload    string `json:"payload"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageButtonResponseContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (buttonResponse NormalizedMessageButtonResponseContent) GetType() string {
	return "ButtonResponse"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (buttonResponse NormalizedMessageButtonResponseContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageButtonResponseContent
	type ButtonResponse struct {
		surrogate
	}
	data, err := json.Marshal(struct {
		ContentType    string         `json:"contentType"`
		ButtonResponse ButtonResponse `json:"buttonResponse"`
	}{
		ContentType: buttonResponse.GetType(),
		ButtonResponse: ButtonResponse{
			surrogate: surrogate(buttonResponse),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (buttonResponse *NormalizedMessageButtonResponseContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageButtonResponseContent
	type ButtonResponse struct {
		surrogate
	}
	var inner struct {
		ButtonResponse ButtonResponse `json:"buttonResponse"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*buttonResponse = NormalizedMessageButtonResponseContent(inner.ButtonResponse.surrogate)
	return
}

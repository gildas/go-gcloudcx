package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

type OpenMessageButtonResponseContent struct {
	ButtonType string `json:"type,omitempty"` // "Button", "QuickReply"
	Text       string `json:"text"`
	Payload    string `json:"payload"`
}

func init() {
	openMessageContentRegistry.Add(OpenMessageButtonResponseContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (buttonResponse OpenMessageButtonResponseContent) GetType() string {
	return "ButtonResponse"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (buttonResponse OpenMessageButtonResponseContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageButtonResponseContent
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
func (buttonResponse *OpenMessageButtonResponseContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageButtonResponseContent
	type ButtonResponse struct {
		surrogate
	}
	var inner struct {
		ButtonResponse ButtonResponse `json:"buttonResponse"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*buttonResponse = OpenMessageButtonResponseContent(inner.ButtonResponse.surrogate)
	return
}

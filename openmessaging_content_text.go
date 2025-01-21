package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageTextContent describes the content of a Text Message
type OpenMessageTextContent struct {
	Text string `json:"body"`
}

func init() {
	openMessageContentRegistry.Add(OpenMessageTextContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (text OpenMessageTextContent) GetType() string {
	return "Text"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (text OpenMessageTextContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageTextContent
	type Text struct {
		surrogate
	}
	data, err := json.Marshal(struct {
		ContentType string `json:"contentType"`
		Text        Text   `json:"text"`
	}{
		ContentType: text.GetType(),
		Text: Text{
			surrogate: surrogate(text),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (text *OpenMessageTextContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageTextContent
	type Text struct {
		surrogate
	}
	var inner struct {
		Text Text `json:"text"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*text = OpenMessageTextContent(inner.Text.surrogate)
	return
}

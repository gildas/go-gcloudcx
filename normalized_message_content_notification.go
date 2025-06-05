package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// NormalizedMessageNotificationContent describes a Notification Content for an OpenMessage
type NormalizedMessageNotificationContent struct {
	ID       string                         `json:"id,omitempty"`
	Language string                         `json:"language"`
	Header   *OpenMessageNotificationHeader `json:"header,omitempty"`
	Body     OpenMessageNotificationBody    `json:"body"`
	Footer   *OpenMessageNotificationFooter `json:"footer,omitempty"`
	Text     string                         `json:"text"`
}

// OpenMessageNotificationHeader describes the Header of a Notification Content
type OpenMessageNotificationHeader struct {
	HeaderType string                              `json:"type"` // "Text", "Media"
	Text       string                              `json:"text,omitempty"`
	Media      *NormalizedMessageAttachmentContent `json:"media,omitempty"`
	Parameters OpenMessageNotificationParameters   `json:"parameters,omitempty"`
}

// OpenMessageNotificationBody describes the Body of a Notification Content
type OpenMessageNotificationBody struct {
	Text       string                            `json:"text"`
	Parameters OpenMessageNotificationParameters `json:"parameters,omitempty"`
}

// OpenMessageNotificationFooter describes the Footer of a Notification Content
type OpenMessageNotificationFooter struct {
	Text string `json:"text"`
}

type OpenMessageNotificationParameters map[string]string

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageNotificationContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (template NormalizedMessageNotificationContent) GetType() string {
	return "Notification"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (notification NormalizedMessageNotificationContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageNotificationContent

	data, err := json.Marshal(struct {
		ContentType string    `json:"contentType"`
		Template    surrogate `json:"template"`
	}{
		ContentType: notification.GetType(),
		Template:    surrogate(notification),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (notification *NormalizedMessageNotificationContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageNotificationContent
	var inner struct {
		Template surrogate `json:"template"`
	}
	if err = json.Unmarshal(payload, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*notification = NormalizedMessageNotificationContent(inner.Template)
	return
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (parameters OpenMessageNotificationParameters) MarshalJSON() ([]byte, error) {
	type Parameter struct {
		Name string `json:"name,omitempty"`
		Text string `json:"text"`
	}
	values := make([]Parameter, 0, len(parameters))
	for name, text := range parameters {
		values = append(values, Parameter{Name: name, Text: text})
	}
	return json.Marshal(values)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (parameters *OpenMessageNotificationParameters) UnmarshalJSON(payload []byte) (err error) {
	var values []struct {
		Name string `json:"name"`
		Text string `json:"text"`
	}
	if err = json.Unmarshal(payload, &values); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*parameters = make(OpenMessageNotificationParameters)
	for _, value := range values {
		(*parameters)[value.Name] = value.Text
	}
	return
}

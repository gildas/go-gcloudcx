package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type NotificationTemplate struct {
	ID          uuid.UUID                    `json:"id"`
	Language    string                       `json:"language"`
	ContentType string                       `json:"contentType"`
	Header      NotificationTemplateHeader   `json:"header"`
	Body        NotificationTemplateBody     `json:"body"`
	Buttons     []NotificationTemplateButton `json:"buttons"`
	Footer      NotificationTemplateFooter   `json:"footer"`
}

type NotificationTemplateHeader struct {
	Type  string `json:"type"` // Text, Media
	Text  string `json:"text"`
	Media Media  `json:"media"` // TODO: Use the full definition
}

type NotificationTemplateBody struct {
	Text string `json:"text"`
}

type NotificationTemplateButton struct {
	Type        string `json:"type"` // QuickReply, PhoneNumber, Url
	Text        string `json:"text"`
	Index       int    `json:"index"`
	PhoneNumber string `json:"phoneNumber"`
	URL         string `json:"url"`
	IsSelected  bool   `json:"isSelected"`
}

type NotificationTemplateFooter struct {
	Text string `json:"text"`
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (template *NotificationTemplate) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NotificationTemplate
	var inner struct {
		surrogate
		ID uuid.UUID `json:"id"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*template = NotificationTemplate(inner.surrogate)
	template.ID = uuid.UUID(inner.ID)

	return nil
}

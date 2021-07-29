package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

type OpenMessageContent struct {
	Type       string                 `json:"contentType"` // Attachment, Location, QuickReply, ButtonResponse, Notification, GenericTemplate, ListTemplate, Postback, Reactions, Mention
	Attachment *OpenMessageAttachment `json:"attachment"`
}

// UnmarshalJSON unmarshals JSON into this
func (content *OpenMessageContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageContent
	var inner surrogate

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*content = OpenMessageContent(inner)
	if content.Attachment == nil {
		return errors.ArgumentMissing.With("attachment")
	}
	// if !Contains(content.Type, []string{"Attachment", "Location", "QuickReply", "ButtonResponse", "Notification", "GenericTemplate", "ListTemplate", "Postback", "Reactions", "Mention"}) {
	if !Contains(content.Type, []string{"Attachment", "Notification"}) {
		return errors.ArgumentInvalid.With("contentType", content.Type)
	}
	return
}

func Contains(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
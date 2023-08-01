package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageContent struct {
	Type       string                 `json:"contentType"` // Attachment, Location, QuickReply, ButtonResponse, Notification, GenericTemplate, ListTemplate, Postback, Reactions, Mention
	Template   *OpenMessageTemplate   `json:"template,omitempty"`
	Attachment *OpenMessageAttachment `json:"attachment,omitempty"`
}

// UnmarshalJSON unmarshals JSON into this
func (content *OpenMessageContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageContent
	var inner surrogate

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*content = OpenMessageContent(inner)
	// if !Contains([]string{"Attachment", "Location", "QuickReply", "ButtonResponse", "Notification", "GenericTemplate", "ListTemplate", "Postback", "Reactions", "Mention"}, content.Type) {
	if !core.Contains([]string{"Attachment", "Notification"}, content.Type) {
		return errors.ArgumentInvalid.With("contentType", content.Type)
	}
	if content.Template == nil && content.Attachment == nil {
		return errors.ArgumentMissing.With("template/attachment")
	}
	return
}

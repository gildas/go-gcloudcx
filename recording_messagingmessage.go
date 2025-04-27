package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type RecordingMessagingMessage struct {
	ID                   string               `json:"id"`
	To                   string               `json:"to"`
	From                 string               `json:"from"`
	FromUser             User                 `json:"fromUser"`
	FromExternalContact  DomainEntityRef      `json:"fromExternalContact"` // TODO: Use the full definition
	ParticipantID        uuid.UUID            `json:"participantId"`
	Purpose              string               `json:"purpose"`
	Queue                DomainEntityRef      `json:"queue"`
	Workflow             DomainEntityRef      `json:"workflow"`
	ContentType          string               `json:"contentType"`
	MessageText          string               `json:"messageText"`
	MediaAttachments     []Media              `json:"messageMediaAttachments"`
	StickerAttachments   []Sticker            `json:"messageStickerAttachments"`
	QuickReplies         []QuickReply         `json:"quickReplies"`
	ButtonResponse       ButtonResponse       `json:"buttonResponse"`
	Story                Story                `json:"story"`
	Cards                []Card               `json:"cards"`
	NotificationTemplate NotificationTemplate `json:"notificationTemplate"`
	Events               []MessageEvent       `json:"events"`
	Timestamp            time.Time            `json:"timestamp"`
}

// UnmarshalJSON unmarshals the recording messaging message from JSON
//
// Implements json.Unmarshaler
func (message *RecordingMessagingMessage) UnmarshalJSON(data []byte) error {
	type surrogate RecordingMessagingMessage
	var inner struct {
		surrogate
		ParticipantID core.UUID `json:"participantId"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*message = RecordingMessagingMessage(inner.surrogate)
	message.ParticipantID = uuid.UUID(inner.ParticipantID)
	return nil
}

package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type GCMessageDetails struct {
	MessageID           uuid.UUID        `json:"messageId"`
	MessageStatus       string           `json:"messageStatus,omitempty"`
	MessageURI          string           `json:"messageUri,omitempty"`
	MessageSegmentCount int64            `json:"messageSegmentCount,omitempty"`
	MessageTime         time.Time        `json:"messageTime"`
	Media               []Media          `json:"media,omitempty"`
	Stickers            []Sticker        `json:"stickers,omitempty"`
	MessageMetadata     *MessageMetadata `json:"messageMetadata,omitempty"`
	SocialVisibility    string           `json:"socialVisibility,omitempty"` // private, public
	ErrorInfo           *APIError        `json:"errorInfo,omitempty"`
}

type MessageMetadata struct {
	Type    string           `json:"type"` // Text, Structured, Receipt, Event, Message, Unknown
	Events  []MessageEvent   `json:"events"`
	Content []MessageContent `json:"content"`
}

type MessageEvent struct {
	EventType string         `json:"eventType"` // CoBrowse, Typing, Presence, Video, Unknown
	SubType   string         `json:"subType"`   // On, Join, Offering, OfferingExpired, OfferingAccepted, OfferingRejected, Disconnect, Clear, SignIn, SessionExpired, Unknown
	CoBrowse  *CobrowseEvent `json:"coBrowse,omitempty"`
	Typing    *TypingEvent   `json:"typing,omitempty"`
	Presence  *PresenceEvent `json:"presence,omitempty"`
	Video     *VideoEvent    `json:"video,omitempty"`
}

type CobrowseEvent struct {
	Type             string `json:"type"` // Offering, OfferingExpired, OfferingAccepted, OfferingRejected
	SessionID        string `json:"sessionId"`
	SessionJoinToken string `json:"sessionJoinToken"`
}

type TypingEvent struct {
	Type     string `json:"type"`     // On
	Duration int64  `json:"duration"` // Duration in milliseconds
}

type PresenceEvent struct {
	Type       string `json:"type"` // Offering, OfferingExpired, OfferingAccepted, OfferingRejected
	OfferingID string `json:"offeringId"`
	JWT        string `json:"jwt"`
}

type VideoEvent struct {
	Type       string `json:"type"` // On, Off
	OfferingID string `json:"offeringId"`
	JWT        string `json:"jwt"`
}

type MessageContent struct {
	ContentType string `json:"contentType"` // Reactions, Attachment, Location, QuickReply, Notification, ButtonResponse, Story, Mention, Card, Carousel, Text, QuickReplyV2, DatePicker, Unknown
	SubType     string `json:"subType"`     // Image, Video, Audio, File, Link, Mention, Reply, Button, QuickReply, Postback. Unknown
}

// UnmarshalJSON unmarshals the message details from JSON
//
// Implements json.UnMarshaler
func (message *GCMessageDetails) UnmarshalJSON(data []byte) error {
	type surrogate GCMessageDetails
	var inner struct {
		surrogate
		MessageID   core.UUID `json:"messageId"`
		MessageTime core.Time `json:"messageTime"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*message = GCMessageDetails(inner.surrogate)
	message.MessageID = uuid.UUID(inner.MessageID)
	message.MessageTime = time.Time(inner.MessageTime)

	return nil
}

package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// ConversationMessage describes a Message (like belonging to Participant)
type ConversationMessage struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Direction string    `json:"direction"` // inbound,outbound
	State     string    `json:"state"`     // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Held      bool      `json:"held"`

	RecordingID string `json:"recordingId,omitempty"`

	Segments         []Segment `json:"segments"`
	Provider         string    `json:"provider"`
	ScriptID         string    `json:"scriptId,omitempty"`
	PeerID           uuid.UUID `json:"peerId"`
	RecipientCountry string    `json:"recipientCountry,omitempty"`
	ToAddress        Address   `json:"toAddress"`
	FromAddress      Address   `json:"fromAddress"`

	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`
	StartHoldTime     time.Time `json:"startHoldTime"`

	Messages []MessageDetails `json:"messages"`

	DisconnectType string    `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	ErrorInfo      ErrorBody `json:"errorInfo"`
	// Screenshares []ScreenShare `json:"screenshares"`
	// SocialExpressions []SocialExpression `json:"socialExpressions"`
	// Videos []Video `json:"videos"`
}

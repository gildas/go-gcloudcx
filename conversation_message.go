package purecloud

import (
	"time"
)

// ConversationMessage describes a Message (like belonging to Participant)
type ConversationMessage struct {
	ID        string `json:"id"`
	Direction string `json:"direction"` // inbound,outbound
	State     string `json:"state"`     // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Held      bool   `json:"held"`

	RecordingID string `json:"recordingId"`

	Segments         []Segment `json:"segments"`
	Provider         string    `json:"provider"`
	ScriptID         string    `json:"scriptId"`
	PeerID           string    `json:"peerId"`
	Type             string    `json:"type"`
	RecipientCountry string    `json:"recipientCountry"`
	RecipientType    string    `json:"recipientType"`
	ToAddress        Address   `json:"toAddress"`
	FromAddress      Address   `json:"fromAddress"`

	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`
	StartHoldTime     time.Time `json:"startHoldTime"`

	Messages []MessageDetail `json:"messages"`

	DisconnectType string    `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	ErrorInfo      ErrorBody `json:"errorInfo"`
}

// MessageDetail  describes details about a Message
type MessageDetail struct {
	ID           string           `json:"messageId"`
	MessageURI   string           `json:"messageURI"`
	Status       string           `json:"messageStatus"`
	SegmentCount string           `json:"messageSegmentCount"`
	Time         time.Time        `json:"messageTime"`
	Media        MessageMedia     `json:"media"`
	Stickers     []MessageSticker `json:"stickers"`
}

// MessageMedia  describes the Media of a Message
type MessageMedia struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	MediaType     string `json:"mediaType"`
	ContentLength string `json:"contentLengthBytes"`
}

// MessageSticker  describes a Message Sticker
type MessageSticker struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

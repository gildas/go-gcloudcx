package purecloud

import (
	"net/url"
	"time"
)

// ConversationChat describes a Chat (like belonging to Participant)
type ConversationChat struct {
	ID        string `json:"id"`
	State     string `json:"state"`     // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Direction string `json:"direction"` // inbound,outbound
	Held      bool   `json:"held"`

	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`
	StartHoldTime     time.Time `json:"startHoldTime"`

	DisconnectType string `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable

	Segments       []Segment       `json:"segments"`
	Provider       string          `json:"provider"`
	PeerID         string          `json:"peerId"`
	RoomID         string          `json:"roomId"`
	ScriptID       string          `json:"scriptId"`
	RecordingID    string          `json:"recordingId"`
	AvatarImageURL *url.URL        `json:"avatarImageUrl"`
	JourneyContext *JourneyContext `json:"journeyContext"`
}

// JourneyContext  describes a Journey Context
type JourneyContext struct {
	Customer         struct {
		ID     string `json:"id"`
		IDType string `json:"idType"`
	} `json:"customer"`
	CustomerSession  struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"customerSession"`
	TriggeringAction struct {
		ID        string    `json:"id"`
		ActionMap struct {
			ID      string `json:"id"`
			Version int    `json:"version"`
		}                   `json:"actionMap"`
	} `json:"triggeringAction"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (conversation ConversationChat) GetID() string {
	return conversation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (conversation ConversationChat) String() string {
	return conversation.ID
}
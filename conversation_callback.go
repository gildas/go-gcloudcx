package purecloud

import (
	"time"

	"github.com/google/uuid"
)

// ConversationCallback describes a Callback (like belonging to Participant)
type ConversationCallback struct {
	ID        uuid.UUID `json:"id"`
	State     string    `json:"state"`     // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Direction string    `json:"direction"` // inbound,outbound
	Held      bool      `json:"held"`

	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`
	StartHoldTime     time.Time `json:"startHoldTime"`
	ScheduledTime     time.Time `json:"callbackScheduledTime"`

	DisconnectType string `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable

	Segments                  []Segment      `json:"segments"`
	Provider                  string         `json:"provider"`
	PeerID                    string         `json:"peerId"`
	DialerPreview             *DialerPreview `json:"dialerPreview"`
	Voicemail                 *Voicemail     `json:"voicemail"`
	CallbackNumbers           []string       `json:"callbackNumbers"`
	CallbackUserName          string         `json:"callbackUserName"`
	ScriptID                  string         `json:"scriptId"`
	AutomatedCallbackConfigID string         `json:"automatedCallbackConfigId"`
}

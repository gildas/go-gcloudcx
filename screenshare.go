package purecloud

import (
	"time"
)

// ScreenShare describes a Screen Share (like belonging to Participant)
type ScreenShare struct {
	ID                string     `json:"id"`
	State             string     `json:"state"`          // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Sharing           bool       `json:"sharing"`

	Segments          []Segment  `json:"segments"`
	PeerCount         int        `json:"peerCount"`
	Provider          string     `json:"provider"`
	PeerID            string     `json:"peerId"`

	ConnectedTime     time.Time  `json:"connectedTime"`
	DisconnectedTime  time.Time  `json:"disconnectedTime"`
	StartAlertingTime time.Time  `json:"startAlertingTime"`

	DisconnectType    string             `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
}
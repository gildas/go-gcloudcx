package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// ConversationVideo describes a Video (like belonging to Participant)
type ConversationVideo struct {
	ID    uuid.UUID  `json:"id"`
	Self  Address    `json:"self"`
	State string     `json:"state"` // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none

	Segments      []Segment `json:"segments"`
	Provider      string    `json:"provider"`
	PeerID        string    `json:"peerId"`
	PeerCount     int       `json:"peerCount"`
	Context       string    `json:"context"`
	AudioMuted    bool      `json:"audioMuted"`
	VideoMuted    bool      `json:"videoMuted"`
	SharingScreen bool      `json:"sharingScreen"`
	MSIDs         []string  `json:"msids"`

	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`

	DisconnectType string `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
}

package purecloud

import (
	"time"
)

// SocialExpression describes a SocialExpression (like belonging to Participant)
type SocialExpression struct {
	ID                string     `json:"id"`
	Direction         string     `json:"direction"`      // inbound,outbound
	State             string     `json:"state"`          // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Held              bool       `json:"held"`

	RecordingID       string     `json:"recordingId"`

	Segments          []Segment  `json:"segments"`
	Provider          string     `json:"provider"`
	ScriptID          string     `json:"scriptId"`
	PeerID            string     `json:"peerId"`
	SocialMediaID     string     `json:"socialMediaId"`
	SocialMediaHub    string     `json:"socialMediaHub"`
	SocialMediaName   string     `json:"socialMediaName"`
	PreviewText       string     `json:"previewText"`


	ConnectedTime     time.Time  `json:"connectedTime"`
	DisconnectedTime  time.Time  `json:"disconnectedTime"`
	StartAlertingTime time.Time  `json:"startAlertingTime"`
	StartHoldTime     time.Time  `json:"startHoldTime"`

	DisconnectType    string             `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
}

// GetID gets the identifier of this
//   implements Identifiable
func (socialExpression SocialExpression) GetID() string {
	return socialExpression.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (socialExpression SocialExpression) String() string {
	return socialExpression.ID
}
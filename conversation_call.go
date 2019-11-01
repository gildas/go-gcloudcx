package purecloud

import (
	"time"
)

// FaxStatus describes a FAX status
type FaxStatus struct {
	Direction        string `json:"direction"`      // inbound,outbound
	ActivePage       int    `json:"activePage"`
	ExpectedPages    int    `json:"expectedPages"`
	LinesTransmitted int    `json:"linesTransmitted"`
	BytesTransmitted int    `json:"bytesTransmitted"`
	BaudRate         int    `json:"baudRate"`
	PageErrors       int    `json:"pageErrors"`
	LineErrors       int    `json:"lineErrors"`
}

// ConversationCall describes a Call (like belonging to Participant)
type ConversationCall struct {
	ID                string     `json:"id"`
	Self              *Address   `json:"self"`
	Direction         string     `json:"direction"`      // inbound,outbound
	State             string     `json:"state"`          // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Muted             bool       `json:"muted"`
	Held              bool       `json:"held"`
	Confined          bool       `json:"confined"`

	Recording         bool       `json:"recording"`
	RecordingState    string     `json:"recodingState"`  // none,active,paused
	RecordingID       string     `json:"recordingId"`

	Segments          []Segment  `json:"segments"`
	DocumentID        string     `json:"documentId"`
	Provider          string     `json:"provider"`
	ScriptID          string     `json:"scriptId"`
	PeerID            string     `json:"peerId"`
	UUIData           string     `json:"uuiData"`
	Other             *Address   `json:"other"`

	ConnectedTime     time.Time  `json:"connectedTime"`
	DisconnectedTime  time.Time  `json:"disconnectedTime"`
	StartAlertingTime time.Time  `json:"startAlertingTime"`
	StartHoldTime     time.Time  `json:"startHoldTime"`

	DisconnectType    string              `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	DisconnectReasons []*DisconnectReason `json:"disconnectReasons"`

	FaxStatus         FaxStatus  `json:"faxStatus"`
	ErrorInfo         ErrorBody  `json:"errorInfo"`
}
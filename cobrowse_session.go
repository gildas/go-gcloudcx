package gcloudcx

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

// CobrowseSession describes a Cobrowse Session (like belonging to Participant)
type CobrowseSession struct {
	ID    uuid.UUID `json:"id"`
	State string    `json:"state"` // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Self  Address   `json:"self"`
	Held  bool      `json:"held"`

	ProviderEventTime time.Time `json:"providerEventTime"`
	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`

	DisconnectType string `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable

	Segments          []Segment `json:"segments"`
	Provider          string    `json:"provider"`
	PeerID            string    `json:"peerId"`
	CobrowseSessionID string    `json:"cobrowseSessionId"`
	CobrowseRole      string    `json:"cobrowseRole"`
	Controlling       []string  `json:"controlling"`
	ViewerURL         *url.URL  `json:"viewerUrl"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (session CobrowseSession) GetID() uuid.UUID {
	return session.ID
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (session CobrowseSession) String() string {
	return session.ID.String()
}

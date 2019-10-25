package purecloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/pkg/errors"
)

// ConversationChat describes a Chat (like belonging to Participant)
type ConversationChat struct {
	ID             string `json:"id"`
	SelfURI        string `json:"selfUri,omitempty"`
	State          string `json:"state"`     // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Direction      string `json:"direction"` // inbound,outbound
	DisconnectType string `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	Held           bool   `json:"held"`


	ConnectedTime     time.Time `json:"connectedTime"`
	DisconnectedTime  time.Time `json:"disconnectedTime"`
	StartAlertingTime time.Time `json:"startAlertingTime"`
	StartHoldTime     time.Time `json:"startHoldTime"`


	Participants   []*Participant  `json:"participants"`
	Segments       []Segment       `json:"segments"`
	Provider       string          `json:"provider"`
	PeerID         string          `json:"peerId"`
	RoomID         string          `json:"roomId"`
	ScriptID       string          `json:"scriptId"`
	RecordingID    string          `json:"recordingId"`
	AvatarImageURL *url.URL        `json:"-"`
	JourneyContext *JourneyContext `json:"journeyContext"`

	Client          *Client         `json:"-"`
	Logger          *logger.Logger  `json:"-"`
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

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *ConversationChat) Initialize(parameters ...interface{}) error {
	client := conversation.Client
	var log *logger.Logger
	for _, parameter := range parameters {
		if paramClient, ok := parameter.(*Client); ok {
			client = paramClient
		}
		if paramLogger, ok := parameter.(*logger.Logger); ok {
			log = paramLogger.Topic("conversation").Scope("conversation").Record("media", "chat")
		}
	}
	if client == nil {
		return errors.Errorf("Missing Client in initialization of %s %s", reflect.TypeOf(conversation).String(), conversation.GetID())
	}
	if log == nil {
		log = client.Logger.Topic("conversation").Scope("conversation").Record("media", "chat")
	}
	conversation.Client = client
	conversation.Logger = log.Topic("conversation").Scope("conversation").Record("media", "chat")
	return conversation.Client.Get("/conversations/chat/" + conversation.GetID(), &conversation)
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

// UpdateState update the state of an identifiable in this
//   implements StateUpdater
func (conversation ConversationChat) UpdateState(identifiable Identifiable, state string) error {
	return conversation.Client.Patch(
		fmt.Sprintf("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: state},
		nil,
	)
}

// UnmarshalJSON unmarshals JSON into this
func (conversation *ConversationChat) UnmarshalJSON(payload []byte) (err error) {
	type surrogate ConversationChat
	var inner struct {
		surrogate
		A *core.URL `json:"avatarImageUrl"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	*conversation = ConversationChat(inner.surrogate)
	conversation.AvatarImageURL = (*url.URL)(inner.A)
	return
}
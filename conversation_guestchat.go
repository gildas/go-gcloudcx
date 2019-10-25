package purecloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ConversationGuestChat describes a Guest Chat
type ConversationGuestChat struct {
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


	Guest          *ChatMember            `json:"member,omitempty"`
	Members        map[string]*ChatMember `json:"-"`

	Segments       []Segment       `json:"segments"`
	Provider       string          `json:"provider"`
	PeerID         string          `json:"peerId"`
	RoomID         string          `json:"roomId"`
	ScriptID       string          `json:"scriptId"`
	RecordingID    string          `json:"recordingId"`
	AvatarImageURL *url.URL        `json:"-"`
	JourneyContext *JourneyContext `json:"journeyContext"`

	JWT             string          `json:"jwt,omitempty"`
	EventStream     string          `json:"eventStreamUri,omitempty"`
	Socket          *websocket.Conn `json:"-"`
	Client          *Client         `json:"-"`
	Logger          *logger.Logger  `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *ConversationGuestChat) Initialize(parameters ...interface{}) error {
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
func (conversation ConversationGuestChat) GetID() string {
	return conversation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (conversation ConversationGuestChat) String() string {
	return conversation.ID
}

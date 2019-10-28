package purecloud

import (
	"fmt"
	"reflect"

	"github.com/gildas/go-logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ConversationGuestChat describes a Guest Chat
type ConversationGuestChat struct {
	ID             string `json:"id"`
	SelfURI        string `json:"selfUri,omitempty"`

	Target         *RoutingTarget         `json:"-"`
	Guest          *ChatMember            `json:"member,omitempty"`
	Members        map[string]*ChatMember `json:"-"`

	JWT             string          `json:"jwt,omitempty"`
	EventStream     string          `json:"eventStreamUri,omitempty"`
	Socket          *websocket.Conn `json:"-"`
	Client          *Client         `json:"-"`
	Logger          *logger.Logger  `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *ConversationGuestChat) Initialize(parameters ...interface{}) (err error) {
	for _, parameter := range parameters {
		if paramClient, ok := parameter.(*Client); ok {
			conversation.Client = paramClient
		}
		if paramLogger, ok := parameter.(*logger.Logger); ok {
			conversation.Logger = paramLogger.Topic("conversation").Scope("conversation").Record("media", "chat")
		}
		if paramGuest, ok := parameter.(*ChatMember); ok {
			conversation.Guest = paramGuest
		}
		if paramTarget, ok := parameter.(*RoutingTarget); ok {
			conversation.Target = paramTarget
		}
	}
	if conversation.Client == nil {
		return errors.Errorf("Missing Client in initialization of %s %s", reflect.TypeOf(conversation).String(), conversation.GetID())
	}
	if conversation.Logger == nil {
		conversation.Logger = conversation.Client.Logger.Topic("conversation").Scope("conversation").Record("media", "chat")
	}
	if conversation.Guest == nil {
		return errors.New("Missing ChatMember guest")
	}
	if conversation.Target == nil {
		return errors.New("Missing ChatMember Target")
	}
	if conversation.Client.Organization == nil {
		return errors.New("Missing Organization in Client")
	}
	if len(conversation.Client.DeploymentID) == 0 {
		return errors.New("Missing Deployment ID in Client")
	}
	conversation.Members = map[string]*ChatMember{}

	if err = conversation.Client.Post("/webchat/guest/conversations",
		struct {
			OrganizationID string         `json:"organizationId"`
			DeploymentID   string         `json:"deploymentId"`
			RoutingTarget  *RoutingTarget `json:"routingTarget"`
			Guest          *ChatMember    `json:"memberInfo"`
		}{
			OrganizationID: conversation.Client.Organization.ID,
			DeploymentID:   conversation.Client.DeploymentID,
			RoutingTarget:  conversation.Target,
			Guest:          conversation.Guest,
		},
		&conversation,
	); err != nil {
		return errors.WithStack(err)
	}
	conversation.Logger = conversation.Logger.Record("conversation", conversation.ID)
	conversation.Members[conversation.Guest.GetID()] = conversation.Guest

	conversation.Socket, _, err = websocket.DefaultDialer.Dial(conversation.EventStream, nil)
	if err != nil {
		conversation.Close()
	}
	return
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

func (conversation *ConversationGuestChat) Close() (err error) {
	log := conversation.Logger.Scope("close")

	if conversation.Socket != nil  {
		log.Debugf("Disconnecting websocket")
		if err = conversation.Socket.Close(); err != nil {
			log.Errorf("Failed while close websocket", err)
			return errors.WithStack(err)
		}
		log.Infof("Disconnected websocket")
	}
	if conversation.Guest != nil {
		log.Debugf("Disconnecting Guest Member")
		if err = conversation.Client.Delete(fmt.Sprintf("/webchat/guest/conversations/%s/members/%s", conversation.ID, conversation.Guest.ID), nil); err != nil {
			log.Errorf("Failed while disconnecting Guest Member", err)
			return errors.WithStack(err)
		}
		log.Infof("Disconnected Guest Member")
	}
	return
}
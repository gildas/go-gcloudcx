package purecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ConversationGuestChat describes a Guest Chat
type ConversationGuestChat struct {
	ID             string                  `json:"id"`
	SelfURI        string                  `json:"selfUri,omitempty"`
	Target         *RoutingTarget          `json:"-"`
	Guest          *ChatMember             `json:"member,omitempty"`
	Members        map[string]*ChatMember  `json:"-"`
	JWT             string                 `json:"jwt,omitempty"`
	EventStream     string                 `json:"eventStreamUri,omitempty"`
	Socket          *websocket.Conn        `json:"-"`
	TopicReceived   chan NotificationTopic `json:"-"`
	LogHeartbeat    bool                   `json:"logHeartbeat"`
	Client          *Client                `json:"-"`
	Logger          *logger.Logger         `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *ConversationGuestChat) Initialize(parameters ...interface{}) (err error) {
	client, logger, err := ExtractClientAndLogger(parameters...)
	if err != nil {
		return err
	}
	guest  := conversation.Guest
	target := conversation.Target
	for _, parameter := range parameters {
		if paramGuest, ok := parameter.(*ChatMember); ok {
			guest = paramGuest
		}
		if paramTarget, ok := parameter.(*RoutingTarget); ok {
			target = paramTarget
		}
	}
	if guest == nil {
		return errors.New("Missing ChatMember guest")
	}
	if target == nil {
		return errors.New("Missing ChatMember Target")
	}
	if client.Organization == nil {
		return errors.New("Missing Organization in Client")
	}
	if len(client.DeploymentID) == 0 {
		return errors.New("Missing Deployment ID in Client")
	}

	if err = client.Post("/webchat/guest/conversations",
		struct {
			OrganizationID string         `json:"organizationId"`
			DeploymentID   string         `json:"deploymentId"`
			RoutingTarget  *RoutingTarget `json:"routingTarget"`
			Guest          *ChatMember    `json:"memberInfo"`
		}{
			OrganizationID: client.Organization.ID,
			DeploymentID:   client.DeploymentID,
			RoutingTarget:  target,
			Guest:          guest,
		},
		&conversation,
	); err != nil {
		return err
	}
	conversation.Client            = client
	conversation.Logger            = logger.Child("conversation", "conversation", "media", "chat", "conversation", conversation.ID)
	// We get the guest's ID from PureCloud, the other fields should be from Initialize
	conversation.Guest.DisplayName = guest.DisplayName
	conversation.Guest.AvatarURL   = guest.AvatarURL
	conversation.Guest.Role        = guest.Role
	conversation.Guest.State       = guest.State
	conversation.Guest.Custom      = guest.Custom
	conversation.Members           = map[string]*ChatMember{}
	conversation.Members[conversation.Guest.ID] = conversation.Guest
	conversation.TopicReceived     = make(chan NotificationTopic)
	conversation.LogHeartbeat      = core.GetEnvAsBool("PURECLOUD_LOG_HEARTBEAT", false)
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

// Connect connects a Guest Chat to its websocket and starts its message loop
//   If the websocket was already connected, nothing happens
//   If the environment variable PURECLOUD_LOG_HEARTBEAT is set to true, the Heartbeat topic will be logged
func (conversation *ConversationGuestChat) Connect() (err error) {
	if conversation.Socket != nil {
		return
	}
	conversation.Socket, _, err = websocket.DefaultDialer.Dial(conversation.EventStream, nil)
	if err != nil {
		conversation.Close()
	}
	go conversation.messageLoop()
	return
}

// Close disconnects the websocket and the guest
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

func (conversation *ConversationGuestChat) messageLoop() (err error) {
	log := conversation.Logger.Scope("receive")

	if conversation.Socket == nil {
		return errors.New("Conversation Not Connected")
	}

	for {
		// get a message body and decode it. (ReadJSON is nice, but in case of unknown message, I cannot get the original string)
		var body []byte

		if _, body, err = conversation.Socket.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("Websocket was closed, stopping receive handler")
				return nil
			}
			log.Errorf("Failed to read incoming message", err)
			continue
		}

		// we have to use a custom version since topics are the same between agent chats and guest chats
		// TODO: Maybe, we should create a Conversation interface...
		topic, err := conversation.notificationTopicFromJSON(body)
		if err != nil {
			log.Warnf("%s, Body size: %d, Content: %s", err.Error(), len(body), string(body))
			continue
		}
		switch topic.(type) {
		case *MetadataTopic:
			if conversation.LogHeartbeat {
				log.Tracef("Request %d bytes: %s", len(body), string(body))
			}
		default:
			log.Tracef("Request %d bytes: %s", len(body), string(body))
		}
		// Make a fake channel object so Notification Topics can be sent through
		topic.Send(&NotificationChannel{
			ID:            conversation.ID,
			LogHeartbeat:  conversation.LogHeartbeat,
			Logger:        conversation.Logger,
			Client:        conversation.Client,
			Socket:        conversation.Socket,
			TopicReceived: conversation.TopicReceived,
		})
	}
}

func (conversation *ConversationGuestChat) notificationTopicFromJSON(payload []byte) (NotificationTopic, error) {
	var header struct {
		TopicName string `json:"topicName"`
		Data      json.RawMessage
	}
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, errors.WithStack(err)
	}
	switch {
	case ConversationGuestChatMessageTopic{}.Match(header.TopicName):
		var topic ConversationGuestChatMessageTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.WithStack(err)
		}
		return &topic, nil
	case ConversationGuestChatMemberTopic{}.Match(header.TopicName):
		var topic ConversationGuestChatMemberTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.WithStack(err)
		}
		return &topic, nil
	case MetadataTopic{}.Match(header.TopicName):
		var topic MetadataTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.WithStack(err)
		}
		return &topic, nil
	default:
		return nil, errors.Errorf("Unsupported Topic: %s", header.TopicName)
	}
}

// GetMember fetches the given member of this Conversation (caches the member)
func (conversation *ConversationGuestChat) GetMember(identifiable Identifiable) (*ChatMember, error) {
	if member, ok := conversation.Members[identifiable.GetID()]; ok {
		return member, nil
	}
	member := &ChatMember{}
	err := conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s", conversation.ID, identifiable.GetID()),
		&request.Options{
			Authorization: "bearer " + conversation.JWT,
		},
		&member,
	)
	if err != nil {
		return nil, err
	}
	conversation.Logger.Scope("getmember").Debugf("Response: %+v", member)
	conversation.Members[member.ID] = member
	return member, nil
}

// SendTyping sends a typing indicator to PureCloud as the chat guest
func (conversation *ConversationGuestChat) SendTyping() (err error) {
	response := &struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name,omitempty"`
		Conversation Conversation `json:"conversation,omitempty"`
		Sender       *ChatMember  `json:"sender,omitempty"`
		Timestamp    time.Time    `json:"timestamp,omitempty"`
	}{}
	if err = conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s/typing", conversation.ID, conversation.Guest.ID),
		&request.Options{
			Method:        http.MethodPost, // since payload is empty
			Authorization: "bearer " + conversation.JWT,
		},
		&response,
	); err == nil {
		conversation.Client.Logger.Record("scope", "sendtyping").Infof("Sent successfuly. Response: %+v", response)
	}
	return
}

// SendMessage sends a message as the chat guest
func (conversation *ConversationGuestChat) SendMessage(text string) (err error) {
	return conversation.sendBody("standard", text)
}

// SendNotice sends a notice as the chat guest
func (conversation *ConversationGuestChat) SendNotice(text string) (err error) {
	return conversation.sendBody("notice", text)
}

// sendBody sends a body message as the chat guest
func (conversation *ConversationGuestChat) sendBody(bodyType, body string) (err error) {
	response := &struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name,omitempty"`
		Conversation Conversation `json:"conversation,omitempty"`
		Sender       *ChatMember  `json:"sender,omitempty"`
		Body         string       `json:"body,omitempty"`
		BodyType     string       `json:"bodyType,omitempty"`
		Timestamp    time.Time    `json:"timestamp,omitempty"`
		SelfURI      string       `json:"selfUri,omitempty"`
	}{}
	if err = conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s/messages", conversation.ID, conversation.Guest.ID),
		&request.Options{
			Authorization: "bearer " + conversation.JWT,
			Payload:       struct {
				BodyType string `json:"bodyType"`
				Body     string `json:"body"`
			}{
				BodyType: bodyType,
				Body:     body,

			},
		},
		&response,
	); err == nil {
		conversation.Client.Logger.Record("scope", "sendbody").Infof("Sent successfuly. Response: %+v", response)
	}
	return
}
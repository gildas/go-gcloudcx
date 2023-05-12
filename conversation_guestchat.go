package gcloudcx

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ConversationGuestChat describes a Guest Chat
type ConversationGuestChat struct {
	ID            uuid.UUID                 `json:"id"`
	SelfURI       URI                       `json:"selfUri,omitempty"`
	Target        *RoutingTarget            `json:"-"`
	Guest         *ChatMember               `json:"member,omitempty"`
	Members       map[uuid.UUID]*ChatMember `json:"-"`
	JWT           string                    `json:"jwt,omitempty"`
	EventStream   string                    `json:"eventStreamUri,omitempty"`
	Socket        *websocket.Conn           `json:"-"`
	TopicReceived chan NotificationTopic    `json:"-"`
	LogHeartbeat  bool                      `json:"logHeartbeat"`
	client        *Client                   `json:"-"`
	logger        *logger.Logger            `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (conversation *ConversationGuestChat) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			conversation.ID = parameter
		case *Client:
			conversation.client = parameter
		case *logger.Logger:
			conversation.logger = parameter.Child("conversation", "conversation", "id", conversation.ID, "media", "guestchat")
		}
	}
	if conversation.logger == nil {
		conversation.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (conversation ConversationGuestChat) GetID() uuid.UUID {
	return conversation.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (conversation ConversationGuestChat) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/conversations/%s", ids[0])
	}
	if conversation.ID != uuid.Nil {
		return NewURI("/api/v2/conversations/%s", conversation.ID)
	}
	return URI("/api/v2/conversations/")
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (conversation ConversationGuestChat) String() string {
	return conversation.ID.String()
}

// Connect connects a Guest Chat to its websocket and starts its message loop
//
//	If the websocket was already connected, nothing happens
//	If the environment variable PURECLOUD_LOG_HEARTBEAT is set to true, the Heartbeat topic will be logged
func (conversation *ConversationGuestChat) Connect(context context.Context) (err error) {
	if conversation.Socket != nil {
		return
	}
	conversation.Socket, _, err = websocket.DefaultDialer.Dial(conversation.EventStream, nil)
	if err != nil {
		_ = conversation.Close(context)
		// return errors.NotConnectedError.With("Conversation")
		return errors.NotConnected.Wrap(err)
	}
	go conversation.messageLoop()
	return
}

// Start starts a Conversation Guest Chat
func (conversation *ConversationGuestChat) Start(ctx context.Context, guest *ChatMember, target *RoutingTarget) error {
	if conversation == nil || conversation.client == nil || conversation.logger == nil {
		return errors.NotInitialized.With("Conversation")
	}
	log := conversation.logger
	client := conversation.client

	if guest == nil {
		return errors.ArgumentMissing.With("Guest")
	}
	if target == nil {
		return errors.ArgumentMissing.With("Target")
	}
	if client.Organization == nil {
		return errors.ArgumentMissing.With("Organization")
	}
	if len(client.DeploymentID) == 0 {
		return errors.ArgumentMissing.With("DeploymentID")
	}

	if err := client.Post(ctx, "/webchat/guest/conversations",
		struct {
			OrganizationID string         `json:"organizationId"`
			DeploymentID   string         `json:"deploymentId"`
			RoutingTarget  *RoutingTarget `json:"routingTarget"`
			Guest          *ChatMember    `json:"memberInfo"`
		}{
			OrganizationID: conversation.client.Organization.ID.String(),
			DeploymentID:   conversation.client.DeploymentID.String(),
			RoutingTarget:  target,
			Guest:          guest,
		},
		&conversation,
	); err != nil {
		return err
	}

	conversation.logger = log
	conversation.client = client
	conversation.Guest.DisplayName = guest.DisplayName
	conversation.Guest.AvatarURL = guest.AvatarURL
	conversation.Guest.Role = guest.Role
	conversation.Guest.State = guest.State
	conversation.Guest.Custom = guest.Custom
	conversation.Members = map[uuid.UUID]*ChatMember{}
	conversation.Members[conversation.Guest.ID] = conversation.Guest
	conversation.TopicReceived = make(chan NotificationTopic)
	conversation.LogHeartbeat = core.GetEnvAsBool("PURECLOUD_LOG_HEARTBEAT", false)
	return nil
}

// Close disconnects the websocket and the guest
func (conversation *ConversationGuestChat) Close(context context.Context) (err error) {
	log := conversation.logger.Scope("close")

	if conversation.Socket != nil {
		log.Debugf("Disconnecting websocket")
		if err = conversation.Socket.Close(); err != nil {
			log.Errorf("Failed while close websocket", err)
			return errors.WithMessage(err, "Failed while closing websocket")
		}
		log.Infof("Disconnected websocket")
	}
	if conversation.Guest != nil {
		log.Debugf("Disconnecting Guest Member")
		if err = conversation.client.Delete(context, NewURI("/webchat/guest/conversations/%s/members/%s", conversation.ID, conversation.Guest.ID), nil); err != nil {
			log.Errorf("Failed while disconnecting Guest Member", err)
			return err
		}
		log.Infof("Disconnected Guest Member")
	}
	return
}

func (conversation *ConversationGuestChat) messageLoop() {
	log := conversation.logger.Scope("receive")

	for {
		// get a message body and decode it. (ReadJSON is nice, but in case of unknown message, I cannot get the original string)
		var err error
		var body []byte

		if _, body, err = conversation.Socket.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("Websocket was closed, stopping receive handler")
				return
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
			Logger:        conversation.logger,
			Client:        conversation.client,
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
		return nil, errors.JSONUnmarshalError.Wrap(err)
	}
	switch {
	case ConversationGuestChatMessageTopic{}.Match(header.TopicName):
		var topic ConversationGuestChatMessageTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.JSONUnmarshalError.Wrap(err)
		}
		return &topic, nil
	case ConversationGuestChatMemberTopic{}.Match(header.TopicName):
		var topic ConversationGuestChatMemberTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.JSONUnmarshalError.Wrap(err)
		}
		return &topic, nil
	case MetadataTopic{}.Match(header.TopicName):
		var topic MetadataTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.JSONUnmarshalError.Wrap(err)
		}
		return &topic, nil
	default:
		return nil, errors.Unsupported.With("Topic", header.TopicName)
	}
}

// GetMember fetches the given member of this Conversation (caches the member)
func (conversation *ConversationGuestChat) GetMember(context context.Context, identifiable Identifiable) (*ChatMember, error) {
	if member, ok := conversation.Members[identifiable.GetID()]; ok {
		return member, nil
	}
	member := &ChatMember{}
	err := conversation.client.SendRequest(
		context,
		NewURI("/webchat/guest/conversations/%s/members/%s", conversation.ID, identifiable.GetID()),
		&request.Options{
			Authorization: "bearer " + conversation.JWT,
		},
		&member,
	)
	if err != nil {
		return nil, err
	}
	conversation.logger.Scope("getmember").Debugf("Response: %+v", member)
	conversation.Members[member.ID] = member
	return member, nil
}

// SendTyping sends a typing indicator to Gcloud as the chat guest
func (conversation *ConversationGuestChat) SendTyping(context context.Context) (err error) {
	response := &struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name,omitempty"`
		Conversation Conversation `json:"conversation,omitempty"`
		Sender       *ChatMember  `json:"sender,omitempty"`
		Timestamp    time.Time    `json:"timestamp,omitempty"`
	}{}
	if err = conversation.client.SendRequest(
		context,
		NewURI("/webchat/guest/conversations/%s/members/%s/typing", conversation.ID, conversation.Guest.ID),
		&request.Options{
			Method:        http.MethodPost,
			Authorization: "bearer " + conversation.JWT,
		},
		&response,
	); err == nil {
		conversation.client.Logger.Record("scope", "sendtyping").Infof("Sent successfully. Response: %+v", response)
	}
	return
}

// SendMessage sends a message as the chat guest
func (conversation *ConversationGuestChat) SendMessage(context context.Context, text string) (err error) {
	return conversation.sendBody(context, "standard", text)
}

// SendNotice sends a notice as the chat guest
func (conversation *ConversationGuestChat) SendNotice(context context.Context, text string) (err error) {
	return conversation.sendBody(context, "notice", text)
}

// sendBody sends a body message as the chat guest
func (conversation *ConversationGuestChat) sendBody(context context.Context, bodyType, body string) (err error) {
	response := &struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name,omitempty"`
		Conversation Conversation `json:"conversation,omitempty"`
		Sender       *ChatMember  `json:"sender,omitempty"`
		Body         string       `json:"body,omitempty"`
		BodyType     string       `json:"bodyType,omitempty"`
		Timestamp    time.Time    `json:"timestamp,omitempty"`
		SelfURI      URI          `json:"selfUri,omitempty"`
	}{}
	if err = conversation.client.SendRequest(
		context,
		NewURI("/webchat/guest/conversations/%s/members/%s/messages", conversation.ID, conversation.Guest.ID),
		&request.Options{
			Authorization: "bearer " + conversation.JWT,
			Payload: struct {
				BodyType string `json:"bodyType"`
				Body     string `json:"body"`
			}{
				BodyType: bodyType,
				Body:     body,
			},
		},
		&response,
	); err == nil {
		conversation.client.Logger.Record("scope", "sendbody").Infof("Sent successfully. Response: %+v", response)
	}
	return
}

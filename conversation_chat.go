package gcloudcx

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// ConversationChat describes a Agent-side Chat
type ConversationChat struct {
	ID                uuid.UUID       `json:"id"`
	SelfURI           URI             `json:"selfUri,omitempty"`
	State             string          `json:"state"`          // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Direction         string          `json:"direction"`      // inbound,outbound
	DisconnectType    string          `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	Held              bool            `json:"held"`
	ConnectedTime     time.Time       `json:"connectedTime"`
	DisconnectedTime  time.Time       `json:"disconnectedTime"`
	StartAlertingTime time.Time       `json:"startAlertingTime"`
	StartHoldTime     time.Time       `json:"startHoldTime"`
	Participants      []*Participant  `json:"participants"`
	Segments          []Segment       `json:"segments"`
	Provider          string          `json:"provider"`
	PeerID            string          `json:"peerId"`
	RoomID            string          `json:"roomId"`
	ScriptID          string          `json:"scriptId"`
	RecordingID       string          `json:"recordingId"`
	AvatarImageURL    *url.URL        `json:"-"`
	JourneyContext    *JourneyContext `json:"journeyContext"`
	client            *Client         `json:"-"`
	logger            *logger.Logger  `json:"-"`
}

// JourneyContext  describes a Journey Context
type JourneyContext struct {
	Customer struct {
		ID     uuid.UUID `json:"id"`
		IDType string    `json:"idType"`
	} `json:"customer"`
	CustomerSession struct {
		ID   uuid.UUID `json:"id"`
		Type string    `json:"type"`
	} `json:"customerSession"`
	TriggeringAction struct {
		ID        uuid.UUID `json:"id"`
		ActionMap struct {
			ID      uuid.UUID `json:"id"`
			Version int       `json:"version"`
		} `json:"actionMap"`
	} `json:"triggeringAction"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloufcx.Client, *logger.Logger
//
// implements Initializable
func (conversation *ConversationChat) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case *Client:
			conversation.client = parameter
		case *logger.Logger:
			conversation.logger = parameter.Child("conversation", "conversation", "id", conversation.ID, "media", "chat")
		}
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (conversation ConversationChat) GetID() uuid.UUID {
	return conversation.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (conversation ConversationChat) GetURI(ids ...uuid.UUID) URI {
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
// implements the fmt.Stringer interface
func (conversation ConversationChat) String() string {
	return conversation.ID.String()
}

// Disconnect disconnect an Identifiable from this
//
// implements Disconnecter
func (conversation ConversationChat) Disconnect(context context.Context, identifiable Identifiable) error {
	return conversation.client.Patch(
		conversation.logger.ToContext(context),
		NewURI("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: "disconnected"},
		nil,
	)
}

// UpdateState update the state of an identifiable in this
//
// implements StateUpdater
func (conversation ConversationChat) UpdateState(context context.Context, identifiable Identifiable, state string) error {
	return conversation.client.Patch(
		conversation.logger.ToContext(context),
		NewURI("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: state},
		nil,
	)
}

// Transfer transfers a participant of this Conversation to the given Queue
//
// implement Transferrer
func (conversation ConversationChat) Transfer(context context.Context, identifiable Identifiable, queue Identifiable) error {
	return conversation.client.Post(
		conversation.logger.ToContext(context),
		NewURI("/conversations/chats/%s/participants/%s/replace", conversation.ID, identifiable.GetID()),
		struct {
			ID string `json:"queueId"`
		}{ID: queue.GetID().String()},
		nil,
	)
}

// Post sends a text message to a chat member
func (conversation ConversationChat) Post(context context.Context, member Identifiable, text string) error {
	return conversation.client.Post(
		conversation.logger.ToContext(context),
		NewURI("/conversations/chats/%s/communications/%s/messages", conversation.ID, member.GetID()),
		struct {
			BodyType string `json:"bodyType"`
			Body     string `json:"body"`
		}{
			BodyType: "standard",
			Body:     text,
		},
		nil,
	)
}

// SetTyping send a typing indicator to the chat member
func (conversation ConversationChat) SetTyping(context context.Context, member Identifiable) error {
	return conversation.client.Post(
		conversation.logger.ToContext(context),
		NewURI("/conversations/chats/%s/communications/%s/typing", conversation.ID, member.GetID()),
		nil,
		nil,
	)
}

// Wrapup wraps up a Participant of this Conversation
func (conversation ConversationChat) Wrapup(context context.Context, identifiable Identifiable, wrapup *Wrapup) error {
	return conversation.client.Patch(
		context,
		NewURI("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{Wrapup: wrapup},
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
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*conversation = ConversationChat(inner.surrogate)
	conversation.AvatarImageURL = (*url.URL)(inner.A)
	return
}

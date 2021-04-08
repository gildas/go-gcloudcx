package purecloud

import (
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
	SelfURI           string          `json:"selfUri,omitempty"`
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
	Client            *Client         `json:"-"`
	Logger            *logger.Logger  `json:"-"`
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

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *ConversationChat) Initialize(parameters ...interface{}) error {
	client, logger, id, err := parseParameters(parameters...)
	if err != nil {
		return err
	}
	// TODO: get /conversations/chats/$id when that REST call works better
	//  At the moment, chat participants do not have any chats even if they are connected. /conversations/$id looks fine
	if id != uuid.Nil {
		if err := conversation.Client.Get(NewURI("/conversations/%s", id), &conversation); err != nil {
			return err
		}
	}
	conversation.Client = client
	conversation.Logger = logger.Topic("conversation").Scope("conversation").Record("media", "chat")
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (conversation ConversationChat) GetID() uuid.UUID {
	return conversation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (conversation ConversationChat) String() string {
	return conversation.ID.String()
}

// Disconnect disconnect an Identifiable from this
//   implements Disconnecter
func (conversation ConversationChat) Disconnect(identifiable Identifiable) error {
	return conversation.Client.Patch(
		NewURI("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: "disconnected"},
		nil,
	)
}

// UpdateState update the state of an identifiable in this
//   implements StateUpdater
func (conversation ConversationChat) UpdateState(identifiable Identifiable, state string) error {
	return conversation.Client.Patch(
		NewURI("/conversations/chats/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: state},
		nil,
	)
}

// Transfer transfers a participant of this Conversation to the given Queue
//   implement Transferrer
func (conversation ConversationChat) Transfer(identifiable Identifiable, queue Identifiable) error {
	return conversation.Client.Post(
		NewURI("/conversations/chats/%s/participants/%s/replace", conversation.ID, identifiable.GetID()),
		struct {
			ID string `json:"queueId"`
		}{ID: queue.GetID().String()},
		nil,
	)
}

// Post sends a text message to a chat member
func (conversation ConversationChat) Post(member Identifiable, text string) error {
	return conversation.Client.Post(
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
func (conversation ConversationChat) SetTyping(member Identifiable) error {
	return conversation.Client.Post(
		NewURI("/conversations/chats/%s/communications/%s/typing", conversation.ID, member.GetID()),
		nil,
		nil,
	)
}

// Wrapup wraps up a Participant of this Conversation
func (conversation ConversationChat) Wrapup(identifiable Identifiable, wrapup *Wrapup) error {
	return conversation.Client.Patch(
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

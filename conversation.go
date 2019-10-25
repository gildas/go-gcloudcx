package purecloud

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gildas/go-logger"
	"github.com/pkg/errors"
)

// Conversation contains the details of a live conversation
//   See: https://developer.mypurecloud.com/api/rest/v2/conversations
type Conversation struct {
	ID              string          `json:"id"`
	SelfURI         string          `json:"selfUri,omitempty"`
	Name            string          `json:"name"`
	StartTime       time.Time       `json:"startTime"`
	EndTime         time.Time       `json:"endTime"`
	Address         string          `json:"address"`
	Participants    []*Participant  `json:"participants"`
	ConversationIDs []string        `json:"conversationIds"`
	MaxParticipants int             `json:"maxParticipants"`
	RecordingState  string          `json:"recordingState"`
	State           string          `json:"state"`
	Divisions       []struct {
		Division DomainEntityRef     `json:"division"`
		Entities []DomainEntityRef   `json:"entities"`
	}                               `json:"divisions"`

	Client          *Client         `json:"-"`
	Logger          *logger.Logger  `json:"-"`
}

// ConversationRoutingData  defines routing details of a Conversation
type ConversationRoutingData struct {
	Queue    AddressableEntityRef   `json:"queue"`
	Language AddressableEntityRef   `json:"language"`
	Priority int                    `json:"priority"`
	Skills   []AddressableEntityRef `json:"skills"`
	ScoredAgents []struct{
		Agent AddressableEntityRef    `json:"agent"`
		Score int                     `json:"score"`
	}                               `json:"scoredAgents"`
}

// Segment describes a fragment of a Conversation
type Segment struct {
	Type            string    `json:"type"`
	DisconnectType  string    `json:"disconnectType"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	HowEnded        string    `json:"howEnded"`
}

// DisconnectReason describes the reason of a disconnect
type DisconnectReason struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Phrase string `json:"phrase"`
}

// DialerPreview describes a Diapler Preview
type DialerPreview struct {
	ID                 string     `json:"id"`
	ContactID          string     `json:"contactId"`
	ContactListID      string     `json:"contactListId"`
	CaompaignID        string     `json:"campaignId"`
	PhoneNumberColumns []struct {
		Type       string `json:"type"`
		ColumnName string `json:"columnName"`
	}                             `json:"phoneNumberColumns"`
}

// Voicemail describes a voicemail
type Voicemail struct {
	ID           string `json:"id"`
	UploadStatus string `json:"uploadStatus"`

}

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *Conversation) Initialize(parameters ...interface{}) error {
	client := conversation.Client
	log    := conversation.Logger
	for _, parameter := range parameters {
		if paramClient, ok := parameter.(*Client); ok {
			client = paramClient
		}
		if paramLogger, ok := parameter.(*logger.Logger); ok {
			log = paramLogger.Topic("conversation").Scope("conversation")
		}
	}
	if client == nil {
		return errors.Errorf("Missing Client in initialization of %s %s", reflect.TypeOf(conversation).String(), conversation.GetID())
	}
	if log == nil {
		log = client.Logger.Topic("conversation").Scope("conversation")
	}
	conversation.Client = client
	conversation.Logger = log.Topic("conversation").Scope("conversation").Record("media", "chat")
	return conversation.Client.Get("/conversations/" + conversation.GetID(), &conversation)
}

// GetID gets the identifier of this
//   implements Identifiable
func (conversation Conversation) GetID() string {
	return conversation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (conversation Conversation) String() string {
	if len(conversation.Name) != 0 {
		return conversation.Name
	}
	return conversation.ID
}

// Post sends a text message to a chat member
func (conversation Conversation) Post(member Identifiable, text string) error {
	return conversation.Client.Post(
		fmt.Sprintf("/conversations/chats/%s/communications/%s/messages", conversation.ID, member.GetID()),
		struct{
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
func (conversation Conversation) SetTyping(member Identifiable) error {
	return conversation.Client.Post(fmt.Sprintf("/conversations/chats/%s/communications/%s/typing", conversation.ID, member.GetID()), nil, nil,)
}

// DisconnectParticipant set the Conversation State of  a participant
func (conversation Conversation) DisconnectParticipant(participant *Participant) error {
	return conversation.Client.Patch(
		fmt.Sprintf("/conversations/chats/%s/participants/%s", conversation.ID, participant.ID),
		MediaParticipantRequest{State: "disconnected"},
		nil,
	)
}

// SetStateParticipant set the Conversation State of  a participant
func (conversation Conversation) SetStateParticipant(participant *Participant, state string) error {
	return conversation.Client.Patch(
		fmt.Sprintf("/conversations/chats/%s/participants/%s", conversation.ID, participant.ID),
		MediaParticipantRequest{State: state},
		nil,
	)
}

// WrapupParticipant wraps up a Participant of this Conversation
func (conversation Conversation) WrapupParticipant(participant *Participant, wrapup *Wrapup) error {
	return conversation.Client.Patch(
		fmt.Sprintf("/conversations/chats/%s/participants/%s", conversation.ID, participant.ID),
		MediaParticipantRequest{Wrapup: wrapup},
		nil,
	)
}

// TransferParticipant transfers a participant of this Conversation to the given Queue
func (conversation Conversation) TransferParticipant(participant *Participant, queue Identifiable) error {
	return conversation.Client.Post(
		fmt.Sprintf("/conversations/chats/%s/participants/%s/replace", conversation.ID, participant.ID),
		struct{ID string `json:"queueId"`}{ID: queue.GetID()},
		nil,
	)
}
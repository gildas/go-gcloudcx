package purecloud

import (
	"fmt"
	"time"

	"github.com/gildas/go-logger"
)

// Conversation contains the details of a live conversation
type Conversation struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	StartTime       time.Time      `json:"startTime"`
	EndTime         time.Time      `json:"endTime"`
	Address         string         `json:"address"`
	Participants    []Participant  `json:"participants"`
	ConversationIDs []string       `json:"conversationIds"`
	MaxParticipants int            `json:"maxParticipants"`
	RecordingState  string         `json:"recordingState"`
	State           string         `json:"state"`
	Divisions       []struct {
		Division DomainEntityRef     `json:"division"`
		Entities []DomainEntityRef   `json:"entities"`
	}                              `json:"divisions"`
	SelfURI         string         `json:"selfUri,omitempty"`

	Client          *Client        `json:"-"`
	Logger          *logger.Logger `json:"-"`
}

// GetConversation get a ConversationChat from its ID
func (client *Client) GetConversation(conversationID string) (*Conversation, error) {
	conversation := &Conversation{}

	if err := client.Get("/conversations/" + conversationID, &conversation); err != nil {
		return nil, err
	}
	conversation.Client = client
	conversation.Logger = client.Logger.Topic("conversation").Scope("conversation")
	return conversation, nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (conversation Conversation) GetID() string {
	return conversation.ID
}

// WrapupParticipant wraps up a Participant of this Conversation
func (conversation Conversation) WrapupParticipant(participant *Participant, wrapup *Wrapup) error {
	return conversation.Client.Patch(
		fmt.Sprintf("/conversations/chat/%s/participants/%s", conversation.ID, participant.ID),
		struct{Wrapup *Wrapup `json:"wrapup"`}{Wrapup: wrapup},
		nil,
	)
}

// TransferParticipant transfers a participant of this Conversation to the given Queue
func (conversation Conversation) TransferParticipant(participant *Participant, queueID string) error {
	return conversation.Client.Post(
		fmt.Sprintf("/conversations/chat/%s/participants/%s/replace", conversation.ID, participant.ID),
		struct{ID string `json:"queueId"`}{ID: queueID},
		nil,
	)
}
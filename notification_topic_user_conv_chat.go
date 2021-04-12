package purecloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// UserConversationChatTopic describes a Topic about User's Presence
type UserConversationChatTopic struct {
	Name          string
	User          *User
	Conversation  *ConversationChat
	Participants  []*Participant
	CorrelationID string
	Client        *Client
}

// Match tells if the given topicName matches this topic
func (topic UserConversationChatTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".conversations.chats")
}

// GetClient gets the PureCloud Client associated with this
func (topic *UserConversationChatTopic) GetClient() *Client {
	return topic.Client
}

// TopicFor builds the topicName for the given identifiables
func (topic UserConversationChatTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.users.%s.conversations.chats", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *UserConversationChatTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Child("user_conversation_chat", "send")
	log.Debugf("User: %s, Conversation: %s (state: %s)", topic.User, topic.Conversation, topic.Conversation.State)
	topic.Client = channel.Client
	topic.User.Client = channel.Client
	topic.Conversation.Client = channel.Client

	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *UserConversationChatTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ConversationID uuid.UUID      `json:"id"`
			Name           string         `json:"name"`
			Participants   []*Participant `json:"participants"`
		} `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId,omitempty"`
		} `json:"metadata,omitempty"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	userID, err := uuid.Parse(strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".conversations.chats"))
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("id", inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.User = &User{ID: userID}
	topic.Conversation = &ConversationChat{ID: inner.EventBody.ConversationID}
	topic.Participants = inner.EventBody.Participants
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic UserConversationChatTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Conversation)
}

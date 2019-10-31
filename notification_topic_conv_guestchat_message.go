package purecloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ConversationGuestChatMessageTopic describes a Topic about User's Presence
type ConversationGuestChatMessageTopic struct {
	ID             string
	Name           string
	Conversation   *ConversationGuestChat
	Sender         *ChatMember
	Type           string     // message, typing-indicator, 
	Body           string
	BodyType       string     // standard,
	TimeStamp      time.Time
	CorrelationID  string
	Client         *Client
}

// Match tells if the given topicName matches this topic
func (topic ConversationGuestChatMessageTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.conversations.chats.") && strings.HasSuffix(topicName, ".messages")
}

// Get the PureCloud Client associated with this
func (topic *ConversationGuestChatMessageTopic) GetClient() *Client {
	return topic.Client
}

// TopicFor builds the topicName for the given identifiables
func (topic ConversationGuestChatMessageTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.conversations.chats.%s.messages", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *ConversationGuestChatMessageTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Topic("conversation_chat_message").Scope("send")
	log.Record("sender", topic.Sender).Debugf("Conversation: %s, Type: %s, Body Type: %s, Sender: %s", topic.Conversation, topic.Type, topic.BodyType, topic.Sender)
	topic.Client              = channel.Client
	topic.Conversation.Client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationGuestChatMessageTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		EventBody struct {
			ID           string                 `json:"id,omitempty"`
			Conversation *ConversationGuestChat `json:"conversation,omitempty"`
			Sender       *ChatMember            `json:"sender,omitempty"`
			Body         string                 `json:"body,omitempty"`
			BodyType     string                 `json:"bodyType,omitempty"`
			Timestamp    time.Time              `json:"timestamp,omitempty"`
		} `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId,omitempty"`
			Type          string `json:"type,omitempty"`
		} `json:"metadata,omitempty"`
		Version   string `json:"version"` // all
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	conversationID := strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.conversations.chats."), ".messages")
	topic.Name          = inner.TopicName
	topic.Type          = inner.Metadata.Type
	topic.Conversation  = &ConversationGuestChat{ID:conversationID}
	topic.Sender        = inner.EventBody.Sender
	topic.BodyType      = inner.EventBody.BodyType
	topic.Body          = inner.EventBody.Body
	topic.TimeStamp     = inner.EventBody.Timestamp
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic ConversationGuestChatMessageTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Sender)
}
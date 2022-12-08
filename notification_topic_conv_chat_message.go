package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationChatMessageTopic describes a Topic about User's Presence
type ConversationChatMessageTopic struct {
	ID            uuid.UUID
	Name          string
	Conversation  *ConversationChat
	Sender        *ChatMember
	Type          string // message, typing-indicator,
	Body          string
	BodyType      string // standard,
	TimeStamp     time.Time
	CorrelationID string
	client        *Client
}

// Match tells if the given topicName matches this topic
func (topic ConversationChatMessageTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.conversations.chats.") && strings.HasSuffix(topicName, ".messages")
}

// GetClient gets the GCloud Client associated with this
func (topic *ConversationChatMessageTopic) GetClient() *Client {
	return topic.client
}

// TopicFor builds the topicName for the given identifiables
func (topic ConversationChatMessageTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.conversations.chats.%s.messages", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *ConversationChatMessageTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Child("conversation_chat_message", "send", "sender", topic.Sender)
	log.Debugf("Conversation: %s, Type: %s, Body Type: %s, Sender: %s", topic.Conversation, topic.Type, topic.BodyType, topic.Sender)
	topic.client = channel.Client
	topic.Conversation.client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationChatMessageTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID           uuid.UUID         `json:"id,omitempty"`
			Conversation *ConversationChat `json:"conversation,omitempty"`
			Sender       *ChatMember       `json:"sender,omitempty"`
			Body         string            `json:"body,omitempty"`
			BodyType     string            `json:"bodyType,omitempty"`
			Timestamp    time.Time         `json:"timestamp,omitempty"`
		} `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId,omitempty"`
			Type          string `json:"type,omitempty"`
		} `json:"metadata,omitempty"`
		Version string `json:"version"` // all
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	conversationID, err := uuid.Parse(strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.conversations.chats."), ".messages"))
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("id", inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.Type = inner.Metadata.Type
	topic.Conversation = &ConversationChat{ID: conversationID}
	topic.Sender = inner.EventBody.Sender
	topic.BodyType = inner.EventBody.BodyType
	topic.Body = inner.EventBody.Body
	topic.TimeStamp = inner.EventBody.Timestamp
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationChatMessageTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Sender)
}

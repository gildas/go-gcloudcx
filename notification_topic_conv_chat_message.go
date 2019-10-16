package purecloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ConversationChatMessageTopic describes a Topic about User's Presence
type ConversationChatMessageTopic struct {
	Name         string
	Conversation *ConversationChat
	Sender       ChatMember
	Body         string
	BodyType     string
	TimeStamp    time.Time
}

// Match tells if the given topicName matches this topic
func (topic ConversationChatMessageTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.conversations.chats.") && strings.HasSuffix(topicName, ".messages")
}

// Send sends the current topic to the Channel's chan
func (topic ConversationChatMessageTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Scope(topic.Name)
	log.Infof("Conversation: %s, Sender: %s", topic.Conversation.ID, topic.Sender.DisplayName)
	topic.Conversation.Client = channel.Client
	channel.TopicReceived <- topic
}

// TopicFor builds the topicName for the given parameter(s)
func (topic ConversationChatMessageTopic) TopicFor(chat *ConversationChat) string {
	if chat != nil {
		return fmt.Sprintf("v2.conversations.chats.%s.messages", chat.ID)
	}
	return ""
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationChatMessageTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		EventBody struct {
			ID        string     `json:"id,omitempty"`
			Sender    ChatMember `json:"sender,omitempty"`
			Body      string     `json:"body,omitempty"`
			BodyType  string     `json:"bodyType,omitempty"`
			Timestamp time.Time  `json:"timestamp,omitempty"`
		} `json:"eventBody"`
		Metadata struct {
			Type          string `json:"type,omitempty"`
		} `json:"metadata,omitempty"`
		Version   string `json:"version"` // all
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	topic.Name         = inner.TopicName
	topic.Conversation = &ConversationChat{ID: strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.v2.conversations.chats."), ".messages")}
	topic.Sender       = inner.EventBody.Sender
	topic.BodyType     = inner.EventBody.BodyType
	topic.Body         = inner.EventBody.Body
	topic.TimeStamp    = inner.EventBody.Timestamp
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic ConversationChatMessageTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Sender.ID)
}
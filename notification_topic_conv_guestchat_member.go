package purecloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ConversationGuestChatMemberTopic describes a Topic about User's Presence
type ConversationGuestChatMemberTopic struct {
	ID             string
	Name           string
	Conversation   *ConversationGuestChat
	Member         *ChatMember
	Type           string     // member-change
	TimeStamp      time.Time
	CorrelationID  string
	Client         *Client
}

// Match tells if the given topicName matches this topic
func (topic ConversationGuestChatMemberTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.conversations.chats.") && strings.HasSuffix(topicName, ".members")
}

// Get the PureCloud Client associated with this
func (topic *ConversationGuestChatMemberTopic) GetClient() *Client {
	return topic.Client
}

// TopicFor builds the topicName for the given identifiables
func (topic ConversationGuestChatMemberTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.conversations.chats.%s.members", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *ConversationGuestChatMemberTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Topic("conversation_chat_member").Scope("send")
	log.Debugf("Conversation: %s, Type: %s, Member: %s", topic.Conversation, topic.Type, topic.Member)
	topic.Client              = channel.Client
	topic.Conversation.Client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationGuestChatMemberTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		EventBody struct {
			ID           string                 `json:"id,omitempty"`
			Conversation *ConversationGuestChat `json:"conversation,omitempty"`
			Member       *ChatMember            `json:"member,omitempty"`
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
	topic.Member        = inner.EventBody.Member
	topic.TimeStamp     = inner.EventBody.Timestamp
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic ConversationGuestChatMemberTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Member)
}
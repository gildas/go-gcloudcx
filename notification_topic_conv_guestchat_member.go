package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationGuestChatMemberTopic describes a Topic about User's Presence
type ConversationChatMemberTopic struct {
	ID             uuid.UUID
	Name           string
	ConversationID uuid.UUID
	Member         *ChatMember
	Type           string // member-change
	TimeStamp      time.Time
	CorrelationID  string
	Targets        []Identifiable
}

func init() {
	notificationTopicRegistry.Add(ConversationChatMemberTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic ConversationChatMemberTopic) GetType() string {
	return "v2.conversations.chats.{id}.members"
}

// GetTargets returns the targets of this topic
func (topic ConversationChatMemberTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic ConversationChatMemberTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationChatMemberTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationChatMemberTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID           string      `json:"id,omitempty"`
			Conversation EntityRef   `json:"conversation,omitempty"`
			Member       *ChatMember `json:"member,omitempty"`
			Timestamp    time.Time   `json:"timestamp,omitempty"`
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
	topic.Name = inner.TopicName
	topic.Type = inner.Metadata.Type
	topic.ConversationID = inner.EventBody.Conversation.ID
	topic.Member = inner.EventBody.Member
	topic.TimeStamp = inner.EventBody.Timestamp
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

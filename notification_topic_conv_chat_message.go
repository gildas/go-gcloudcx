package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationChatMessageTopic describes a Topic about User's Presence
type ConversationChatMessageTopic struct {
	ID             uuid.UUID
	Name           string
	ConversationID uuid.UUID
	Sender         *ChatMember
	Type           string // message, typing-indicator,
	Body           string
	BodyType       string // standard,
	TimeStamp      time.Time
	CorrelationID  string
	Targets        []Identifiable
}

func init() {
	notificationTopicRegistry.Add(ConversationChatMessageTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic ConversationChatMessageTopic) GetType() string {
	return "v2.conversations.chats.{id}.messages"
}

// GetTargets returns the targets of this topic
func (topic ConversationChatMessageTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic ConversationChatMessageTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// TopicNameWith builds the topicName for the given identifiables
func (topic ConversationChatMessageTopic) TopicNameWith(identifiables ...Identifiable) string {
	return topicNameWith(topic, identifiables...)
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationChatMessageTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationChatMessageTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID           uuid.UUID   `json:"id,omitempty"`
			Conversation EntityRef   `json:"conversation,omitempty"`
			Sender       *ChatMember `json:"sender,omitempty"`
			Body         string      `json:"body,omitempty"`
			BodyType     string      `json:"bodyType,omitempty"`
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
	topic.Sender = inner.EventBody.Sender
	topic.BodyType = inner.EventBody.BodyType
	topic.Body = inner.EventBody.Body
	topic.TimeStamp = inner.EventBody.Timestamp
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

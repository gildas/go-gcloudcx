package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// UserConversationChatTopic describes a Topic about User's Presence
type UserConversationChatTopic struct {
	Name           string
	User           *User
	ConversationID uuid.UUID
	Participants   []*Participant
	CorrelationID  string
	Targets        []Identifiable
}

func init() {
	notificationTopicRegistry.Add(UserConversationChatTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic UserConversationChatTopic) GetType() string {
	return "v2.users.{id}.conversations.chats"
}

// GetTargets returns the targets of this topic
func (topic UserConversationChatTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic UserConversationChatTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic UserConversationChatTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
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
	found, targets := getTargets(topic.GetType(), inner.TopicName)
	if !found || len(targets) == 0 {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("topicName", inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.User = &User{ID: targets[0].GetID()}
	topic.ConversationID = inner.EventBody.ConversationID
	topic.Participants = inner.EventBody.Participants
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

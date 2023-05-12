package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// UserPresenceTopic describes a Topic about User's Presence
type UserPresenceTopic struct {
	Name          string
	User          *User
	Presence      UserPresence
	CorrelationID string
	Targets       []Identifiable
}

func init() {
	notificationTopicRegistry.Add(UserPresenceTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic UserPresenceTopic) GetType() string {
	return "v2.users.{id}.presence"
}

// GetTargets returns the targets of this topic
func (topic UserPresenceTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic UserPresenceTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic UserPresenceTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *UserPresenceTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		Presence  UserPresence `json:"eventBody"`
		Metadata  struct {
			CorrelationID string `json:"correlationId"`
		} `json:"metadata"`
		Version string `json:"version"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	found, targets := getTargets(topic.GetType(), inner.TopicName)
	if !found || len(targets) == 0 {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("topicName", inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.Presence = inner.Presence
	topic.User = &User{ID: targets[0].GetID()}
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

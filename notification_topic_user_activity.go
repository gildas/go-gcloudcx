package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// UserActivityTopic describes a Topic about User's Activity
type UserActivityTopic struct {
	Name          string
	User          *User
	Presence      *UserPresence
	RoutingStatus *RoutingStatus
	CorrelationID string
	ActiveQueues  []*Queue
	Targets       []Identifiable
}

func init() {
	notificationTopicRegistry.Add(UserActivityTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic UserActivityTopic) GetType() string {
	return "v2.users.{id}.activity"
}

// GetTargets returns the targets of this topic
func (topic UserActivityTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic UserActivityTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic UserActivityTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *UserActivityTopic) UnmarshalJSON(payload []byte) (err error) {
	// TODO: Put this schema:
	/*
		{
		  "activeQueueIds": [
		    {}
		  ],
		  "dateActiveQueuesChanged": "string"
		}
	*/
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID                      uuid.UUID      `json:"id"`
			RoutingStatus           *RoutingStatus `json:"routingStatus"`
			Presence                *UserPresence  `json:"presence"`
			OutOfOffice             *OutOfOffice   `json:"outOfOffice"`
			ActiveQueueIDs          []string       `json:"activeQueueIds"` // TODO: Not sure about this (the doc says: "activeQueueIds": [{}])
			DateActiveQueuesChanged string         `json:"dateActiveQueuesChanged"`
		}
		Metadata struct {
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
	topic.User = &User{ID: targets[0].GetID()}
	topic.Presence = inner.EventBody.Presence
	topic.RoutingStatus = inner.EventBody.RoutingStatus
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

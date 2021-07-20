package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strings"

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
	Client        *Client
}

// Match tells if the given topicName matches this topic
func (topic UserActivityTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".activity")
}

// GetClient gets the GCloud Client associated with this
func (topic *UserActivityTopic) GetClient() *Client {
	return topic.Client
}

// TopicFor builds the topicName for the given identifiables
func (topic UserActivityTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.users.%s.activity", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *UserActivityTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Child("user_activity", "send")
	log.Debugf("User: %s, New Presence: %s", topic.User, topic.Presence)
	topic.Client = channel.Client
	topic.User.Client = channel.Client
	channel.TopicReceived <- topic
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
	userID, err := uuid.Parse(strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".activity"))
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("id",inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.User = &User{ID: userID}
	topic.Presence = inner.EventBody.Presence
	topic.RoutingStatus = inner.EventBody.RoutingStatus
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic UserActivityTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Presence)
}

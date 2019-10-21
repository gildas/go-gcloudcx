package purecloud

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/pkg/errors"
)

// UserActivityTopic describes a Topic about User's Activity
type UserActivityTopic struct {
	Name          string
	User          *User
	Presence      *UserPresence
	CorrelationID string
	Client        *Client
}

// Match tells if the given topicName matches this topic
func (topic UserActivityTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".activity")
}

// Get the PureCloud Client associated with this
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
	log := channel.Logger.Topic("user_activity").Scope("send")
	log.Debugf("User: %s, New Presence: %s", topic.User, topic.Presence)
	topic.Client      = channel.Client
	topic.User.Client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *UserActivityTopic) UnmarshalJSON(payload []byte) (err error) {
	// TODO: Put this schema:
	/*
{
  "id": "string",
  "routingStatus": {
    "status": "OFF_QUEUE|IDLE|INTERACTING|NOT_RESPONDING|COMMUNICATING",
    "startTime": "string"
  },
  "presence": {
    "presenceDefinition": {
      "id": "string",
      "systemPresence": "string"
    },
    "presenceMessage": "string",
    "modifiedDate": "string"
  },
  "outOfOffice": {
    "active": true,
    "modifiedDate": "string"
  },
  "activeQueueIds": [
    {}
  ],
  "dateActiveQueuesChanged": "string"
}
	*/
	var inner struct {
		TopicName string        `json:"topicName"`
		Presence  *UserPresence `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId"`
		}                       `json:"metadata"`
		Version   string        `json:"version"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	userID := strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".activity")
	topic.Name          = inner.TopicName
	topic.User          = &User{ID:userID}
	topic.Presence      = inner.Presence
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic UserActivityTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Presence)
}

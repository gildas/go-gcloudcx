package purecloud

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/pkg/errors"
)

// UserPresenceTopic describes a Topic about User's Presence
type UserPresenceTopic struct {
	Name     string
	User     *User
	Presence UserPresence
}

// Match tells if the given topicName matches this topic
func (topic UserPresenceTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".presence")
}

// Send sends the current topic to the Channel's chan
func (topic UserPresenceTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Scope(topic.Name)
	log.Infof("User: %s, New Presence: %s", topic.User.ID, topic.Presence.String())
	topic.User.Client = channel.Client
	channel.TopicReceived <- topic
}

// TopicFor builds the topicName for the given parameter(s)
func (topic UserPresenceTopic) TopicFor(user *User) string {
	if user != nil {
		return fmt.Sprintf("v2.users.%s.presence", user.ID)
	}
	return ""
}

// UnmarshalJSON unmarshals JSON into this
func (topic *UserPresenceTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		Presence  UserPresence `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId"`
		}                      `json:"metadata"`
		Version   string       `json:"version"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	topic.Name     = inner.TopicName
	topic.User     = &User{ID: strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".presence")}
	topic.Presence = inner.Presence
	return
}

// String gets a string version
//   implements the fmt.Stringer interface
func (topic UserPresenceTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Presence)
}
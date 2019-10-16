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
	UserID   string
	Presence UserPresence
}

func (topic UserPresenceTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".presence")
}

func (topic UserPresenceTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Scope(topic.Name)
	log.Infof("User: %s, New Presence: %s", topic.UserID, topic.Presence.String())
}

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
	topic.UserID   = strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".presence")
	topic.Presence = inner.Presence
	return
}

func (topic UserPresenceTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.Presence)
}
package purecloud

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// MetadataTopic describes a Topic about the channel itself
type MetadataTopic struct {
	Name    string
	Message string
}

func (topic MetadataTopic) Match(topicName string) bool {
	return topicName == "channel.metadata"
}

func (topic MetadataTopic) Send(channel *NotificationChannel) {
	if topic.Message == "WebSocket Heartbeat" && !channel.LogHeartbeat {
		return
	}
	log := channel.Logger.Scope(topic.Name)
	
	log.Infof("Topic Message: %s", topic.Message)
}

func (topic *MetadataTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			Message string `json:"message"`
		}                `json:"eventBody"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	topic.Name    = inner.TopicName
	topic.Message = inner.EventBody.Message
	return
}
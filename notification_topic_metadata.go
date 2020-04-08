package purecloud

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// MetadataTopic describes a Topic about the channel itself
type MetadataTopic struct {
	Name    string
	Message string
	Client  *Client
}

// Match tells if the given topicName matches this topic
func (topic MetadataTopic) Match(topicName string) bool {
	return topicName == "channel.metadata"
}

// GetClient gets the PureCloud Client associated with this
func (topic *MetadataTopic) GetClient() *Client {
	return topic.Client
}

// TopicFor builds the topicName for the given identifiables
func (topic MetadataTopic) TopicFor(identifiables ...Identifiable) string {
	return "channel.metadata"
}

// Send sends the current topic to the Channel's chan
func (topic *MetadataTopic) Send(channel *NotificationChannel) {
	if topic.Message == "WebSocket Heartbeat" && !channel.LogHeartbeat {
		return
	}
	log := channel.Logger.Scope(topic.Name)

	log.Debugf("Topic Message: %s", topic.Message)
	topic.Client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *MetadataTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			Message string `json:"message"`
		} `json:"eventBody"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	topic.Name = inner.TopicName
	topic.Message = inner.EventBody.Message
	return
}

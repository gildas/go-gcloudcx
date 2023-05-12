package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// MetadataTopic describes a Topic about the channel itself
type MetadataTopic struct {
	Name    string
	Message string
}

func init() {
	notificationTopicRegistry.Add(MetadataTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic MetadataTopic) GetType() string {
	return "channel.metadata"
}

// GetTargets returns the targets of this topic
func (topic MetadataTopic) GetTargets() []Identifiable {
	return nil
}

// With creates a new NotificationTopic with the given targets
func (topic MetadataTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	return newTopic
}

// String gets a string version
//
// implements fmt.Stringer
func (topic MetadataTopic) String() string {
	return topic.GetType()
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

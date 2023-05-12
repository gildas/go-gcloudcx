package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationVoicemailEndTopic describes a Topic about a Conversation Voicemail End
//
// See: https://developer.genesys.cloud/notificationsalerts/notifications/available-topics#v2-detail-events-conversation--id--voicemail-end
type ConversationVoicemailEndTopic struct {
	ID             uuid.UUID
	Name           string
	ConversationID uuid.UUID
	CorrelationID  string
	Targets        []Identifiable
}

func init() {
	notificationTopicRegistry.Add(ConversationVoicemailEndTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic ConversationVoicemailEndTopic) GetType() string {
	return "v2.detail.events.conversation.{id}.acd.end"
}

// GetTargets returns the targets
func (topic ConversationVoicemailEndTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic ConversationVoicemailEndTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationVoicemailEndTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationVoicemailEndTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID           uuid.UUID `json:"id,omitempty"`
			Conversation EntityRef `json:"conversation,omitempty"`
		} `json:"eventBody"`
		Metadata struct {
			CorrelationID string `json:"correlationId,omitempty"`
			Type          string `json:"type,omitempty"`
		} `json:"metadata,omitempty"`
		Version string `json:"version"` // all
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	topic.Name = inner.TopicName
	topic.ConversationID = inner.EventBody.Conversation.ID
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

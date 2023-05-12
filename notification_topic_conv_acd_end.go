package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationACDEndTopic describes a Topic about a Conversation ACD End
//
// See: https://developer.genesys.cloud/notificationsalerts/notifications/available-topics#v2-detail-events-conversation--id--acd-end
type ConversationACDEndTopic struct {
	ID             uuid.UUID
	Name           string
	Time           time.Time
	ConversationID uuid.UUID
	SessionID      uuid.UUID
	QueueID        uuid.UUID
	DivisionID     uuid.UUID
	CorrelationID  string
	DisconnectType string
	MediaType      string
	MessageType    string
	Provider       string
	Direction      string
	ANI            string
	DNIS           string
	AddressTo      string
	AddressFrom    string
	Subject        string
	ACDOutcome     string
	Participant    *Participant
	Targets        []Identifiable
}

func init() {
	notificationTopicRegistry.Add(ConversationACDEndTopic{})
}

// GetType returns the type of this topic
//
// implements core.TypeCarrier
func (topic ConversationACDEndTopic) GetType() string {
	return "v2.detail.events.conversation.{id}.acd.end"
}

// GetTargets returns the targets
func (topic ConversationACDEndTopic) GetTargets() []Identifiable {
	return topic.Targets
}

// With creates a new NotificationTopic with the given targets
func (topic ConversationACDEndTopic) With(targets ...Identifiable) NotificationTopic {
	newTopic := topic
	newTopic.Targets = targets
	return newTopic
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationACDEndTopic) String() string {
	if len(topic.Targets) == 0 {
		return topic.GetType()
	}
	return topicNameWith(topic, topic.Targets...)
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationACDEndTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID             uuid.UUID      `json:"id,omitempty"`
			ConversationID uuid.UUID      `json:"conversationId,omitempty"`
			Time           core.Timestamp `json:"eventTime"`
			SessionID      uuid.UUID
			QueueID        uuid.UUID
			DivisionID     uuid.UUID
			ParticipantID  uuid.UUID
			CorrelationID  string
			DisconnectType string
			MediaType      string
			MessageType    string
			Provider       string
			Direction      string
			ANI            string
			DNIS           string
			AddressTo      string
			AddressFrom    string
			Subject        string
			ACDOutcome     string
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
	topic.ConversationID = inner.EventBody.ConversationID
	topic.SessionID = inner.EventBody.SessionID
	topic.QueueID = inner.EventBody.QueueID
	topic.DivisionID = inner.EventBody.DivisionID
	topic.CorrelationID = inner.EventBody.CorrelationID
	topic.Time = time.Time(inner.EventBody.Time)
	topic.DisconnectType = inner.EventBody.DisconnectType
	topic.MediaType = inner.EventBody.MediaType
	topic.MessageType = inner.EventBody.MessageType
	topic.Provider = inner.EventBody.Provider
	topic.Direction = inner.EventBody.Direction
	topic.ANI = inner.EventBody.ANI
	topic.DNIS = inner.EventBody.DNIS
	topic.AddressTo = inner.EventBody.AddressTo
	topic.AddressFrom = inner.EventBody.AddressFrom
	topic.Subject = inner.EventBody.Subject
	topic.ACDOutcome = inner.EventBody.ACDOutcome
	topic.CorrelationID = inner.Metadata.CorrelationID
	topic.Participant = &Participant{ID: inner.EventBody.ParticipantID}
	return
}

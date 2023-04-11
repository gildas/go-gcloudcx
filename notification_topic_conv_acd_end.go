package gcloudcx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ConversationACDEndTopic describes a Topic about a Conversation ACD End
//
// See: https://developer.genesys.cloud/notificationsalerts/notifications/available-topics#v2-detail-events-conversation--id--acd-end
type ConversationACDEndTopic struct {
	ID            uuid.UUID
	Name          string
	Conversation  *Conversation
	CorrelationID string
	client        *Client
}

// Match tells if the given topicName matches this topic
func (topic ConversationACDEndTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.detail.events.conversation.") && strings.HasSuffix(topicName, ".acd.end")
}

// GetClient gets the GCloud Client associated with this
func (topic *ConversationACDEndTopic) GetClient() *Client {
	return topic.client
}

// TopicFor builds the topicName for the given identifiables
func (topic ConversationACDEndTopic) TopicFor(identifiables ...Identifiable) string {
	if len(identifiables) > 0 {
		return fmt.Sprintf("v2.detail.events.conversation.%s.acd.end", identifiables[0].GetID())
	}
	return ""
}

// Send sends the current topic to the Channel's chan
func (topic *ConversationACDEndTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Child("conversation_acd_end", "send")
	log.Debugf("Conversation: %s", topic.Conversation)
	topic.client = channel.Client
	topic.Conversation.client = channel.Client
	channel.TopicReceived <- topic
}

// UnmarshalJSON unmarshals JSON into this
func (topic *ConversationACDEndTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string `json:"topicName"`
		EventBody struct {
			ID           uuid.UUID     `json:"id,omitempty"`
			Conversation *Conversation `json:"conversation,omitempty"`
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
	conversationID, err := uuid.Parse(strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.conversations.chats."), ".messages"))
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentInvalid.With("id", inner.TopicName))
	}
	topic.Name = inner.TopicName
	topic.Conversation = &Conversation{ID: conversationID}
	topic.CorrelationID = inner.Metadata.CorrelationID
	return
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (topic ConversationACDEndTopic) String() string {
	// TODO: Use relevant fields
	return fmt.Sprintf("%s=%s", topic.Name, "")
}

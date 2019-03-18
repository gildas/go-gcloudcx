package conversation

import (
	"github.com/gildas/go-purecloud"
)

// Conversation contains the details of a live chat/conversation
type Conversation struct {
	ID          string `json:"id,omitifempty"`
	JWT         string `json:"jwt,omitifempty"`
	EventStream string `json:"eventStreamUri,omitifempty"`
	Member      Member `json:"member,omitifempty"`

	Client      *purecloud.Client `json:"-"`
}

// Member describes a chat guest member
type Member struct {
	ID          string            `json:"id,omitifempty"`
	State       string            `json:"state"`
	DisplayName string            `json:"displayName,omitifempty"`
	ImageURL    string            `json:"avatarImageUrl,omitifempty"`
	Custom      map[string]string `json:"customFields,omitifempty"`
}

// Target describes the target of a Chat/Conversation
type Target struct {
	Type    string `json:"targetType,omitifempty"`
	Address string `json:"targetAddress,omitifempty"`
}

// Message describes messages exchanged over a websocket
type Message struct {
	TopicName string `json:"topicName,omitifempty"`
	EventBody struct {
		ID           string       `json:"id,omitifempty"`           // typing-indicator, message
		Sender       Member       `json:"sender,omitifempty"`       // typing-indicator, message
		Body         string       `json:"body,omitifempty"`         // message
		BodyType     string       `json:"bodyType,omitifempty"`     // message
		Message      string       `json:"message,omitifempty"`      // heartbeat (channel.metadata)
		Conversation Conversation `json:"conversation,omitifempty"` // typing-indicator, member-change
		Member       Member       `json:"member,omitifempty"`       // member-change
		Timestamp    string       `json:"timestamp,omitifempty"`    // time.Time!?, all
	} `json:"eventBody,omitifempty"`
	Metadata struct {
		CorrelationID string `json:"CorrelationId,omitifempty"` // typing-indicator
		Type          string `json:"type,omitifempty"`          // typing-indicator, message, member-change
	} `json:"metadata,omitifempty"`
	Version   string `json:"version"` // all
}
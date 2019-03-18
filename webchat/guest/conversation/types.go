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
	TopicName string `json:"topicName"`
	EventBody struct {
		Message      string       `json:"message"` // if TopicName == "channel.metadata"
		Conversation Conversation `json:"conversation"`
		Member       Member       `json:"member"`
		Timestamp    string       `json:"timestamp"` // time.Time!?
	} `json:"eventBody"`
	Metadata struct {
		CorrelationID string `json:"CorrelationId"`
		Type          string `json:"type"`
	} `json:"metadata"`
	Version   string `json:"version"`
}
package conversation

import (
	"time"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
	"github.com/gorilla/websocket"
)

// Conversation contains the details of a live chat/conversation
type Conversation struct {
	ID          string `json:"id,omitempty"`
	JWT         string `json:"jwt,omitempty"`
	EventStream string `json:"eventStreamUri,omitempty"`
	Guest       Member `json:"member,omitempty"`
	SelfURI     string `json:"selfUri,omitempty"`

	Client      *purecloud.Client  `json:"-"`
	Socket      *websocket.Conn    `json:"-"`
	Members     map[string]*Member `json:"-"`
	Logger      *logger.Logger     `json:"-"`
}

// Member describes a chat guest member
type Member struct {
	ID            string            `json:"id,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	AvatarURL     string            `json:"avatarImageUrl,omitempty"`
	Role          string            `json:"role,omitempty"`
	State         string            `json:"state,omitempty"`
	JoinedAt      time.Time         `json:"joinDate,omitempty"`
	LeftAt        time.Time         `json:"leaveDate,omitempty"`
	Authenticated bool              `json:"authenticatedGuest,omitempty"`
	Custom        map[string]string `json:"customFields,omitempty"`
}

// Target describes the target of a Chat/Conversation
type Target struct {
	Type    string `json:"targetType,omitempty"`
	Address string `json:"targetAddress,omitempty"`
}

// Message describes messages exchanged over a websocket
type Message struct {
	TopicName string `json:"topicName,omitempty"`
	EventBody struct {
		ID           string       `json:"id,omitempty"`           // typing-indicator, message
		Sender       Member       `json:"sender,omitempty"`       // typing-indicator, message
		Body         string       `json:"body,omitempty"`         // message
		BodyType     string       `json:"bodyType,omitempty"`     // message
		Message      string       `json:"message,omitempty"`      // heartbeat (channel.metadata)
		Conversation Conversation `json:"conversation,omitempty"` // typing-indicator, member-change
		Member       Member       `json:"member,omitempty"`       // member-change
		Timestamp    time.Time    `json:"timestamp,omitempty"`    // all
	} `json:"eventBody,omitempty"`
	Metadata struct {
		CorrelationID string `json:"CorrelationId,omitempty"` // typing-indicator
		Type          string `json:"type,omitempty"`          // typing-indicator, message, member-change
	} `json:"metadata,omitempty"`
	Version   string `json:"version"` // all

	Logger    *logger.Logger `json:"-"`
}
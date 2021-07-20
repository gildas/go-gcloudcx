package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// AgentlessMessage sends an agentless outbound text message to a Messenger address
//
// See https://developer.genesys.cloud/api/rest/v2/conversations/#post-api-v2-conversations-messages-agentless
type AgentlessMessage struct {
	From          string             `json:"fromAddress"`
	To            string             `json:"toAddress"`
	MessengerType string             `json:"toAddressMessengerType"`
	Text          string             `json:"textBody"`
	Template      *MessagingTemplate `json:"messagingTemplate,omitempty"`
}

// MessagingTemplate describes the Template to use (WhatsApp Template, for example)
type MessagingTemplate struct {
	ResponseID string              `json:"responseId"`
	Parameters []TemplateParameter `json:"parameters"`
}

// TemplateParameter describes a template parameter
type TemplateParameter struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// AgentlessMessageResult describes the results of the send action on an AgentlessMessage
type AgentlessMessageResult struct {
	ID             string               `json:"id"`
	ConversationID uuid.UUID            `json:"conversationId"`
	From          string                `json:"fromAddress"`
	To            string                `json:"toAddress"`
	MessengerType string                `json:"messengerType"`
	Text          string                `json:"textBody"`
	Template      *MessagingTemplate    `json:"messagingTemplate,omitempty"`
	JobUser       *AddressableEntityRef `json:"user,omitempty"`
	Timestamp     time.Time             `json:"timestamp"`
	SelfURI       string                `json:"selfUri"`
}

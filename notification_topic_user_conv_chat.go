package purecloud

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/pkg/errors"
)

// UserConversationChatTopic describes a Topic about User's Presence
type UserConversationChatTopic struct {
	Name           string
	UserID         string
	ConversationID string
	Participants   []ConversationParticipant
}

func (topic UserConversationChatTopic) Match(topicName string) bool {
	return strings.HasPrefix(topicName, "v2.users.") && strings.HasSuffix(topicName, ".conversations.chats")
}

func (topic UserConversationChatTopic) Send(channel *NotificationChannel) {
	log := channel.Logger.Scope(topic.Name)
	log.Infof("User: %s, Conversation: %s", topic.UserID, topic.ConversationID)
	channel.TopicReceived <- topic
}

func (topic *UserConversationChatTopic) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		TopicName string       `json:"topicName"`
		EventBody struct {
			ConversationID string                    `json:"id"`
			Name           string                    `json:"name"`
			Participants   []ConversationParticipant `json:"participants"`

		} `json:"eventBody"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	topic.Name           = inner.TopicName
	topic.UserID         = strings.TrimSuffix(strings.TrimPrefix(inner.TopicName, "v2.users."), ".conversations.chats")
	topic.ConversationID = inner.EventBody.ConversationID
	topic.Participants   = inner.EventBody.Participants
	return
}

func (topic UserConversationChatTopic) String() string {
	return fmt.Sprintf("%s=%s", topic.Name, topic.ConversationID)
}

/*
{
  "id": "string",
  "name": "string",
  "participants": [
    {
      "id": "string",
      "name": "string",
      "address": "string",
      "startTime": "string",
      "connectedTime": "string",
      "endTime": "string",
      "startHoldTime": "string",
      "purpose": "string",
      "state": "alerting|dialing|contacting|offering|connected|disconnected|terminated|converting|uploading|transmitting|scheduled|none",
      "direction": "inbound|outbound",
      "disconnectType": "endpoint|client|system|transfer|timeout|transfer.conference|transfer.consult|transfer.forward|transfer.noanswer|transfer.notavailable|transport.failure|error|peer|other|spam|uncallable",
      "held": true,
      "wrapupRequired": true,
      "wrapupPrompt": "string",
      "user": {
        "id": "string",
        "name": "string"
      },
      "queue": "object",
      "attributes": "object",
      "errorInfo": {
        "status": 0,
        "code": "string",
        "entityId": "string",
        "entityName": "string",
        "message": "string",
        "messageWithParams": "string",
        "messageParams": "object",
        "contextId": "string",
        "details": [
          {
            "errorCode": "string",
            "fieldName": "string",
            "entityId": "string",
            "entityName": "string"
          }
        ],
        "errors": [
          {}
        ]
      },
      "script": "object",
      "wrapupTimeoutMs": 0,
      "wrapupSkipped": true,
      "alertingTimeoutMs": 0,
      "provider": "string",
      "externalContact": "object",
      "externalOrganization": "object",
      "wrapup": {
        "code": "string",
        "notes": "string",
        "tags": [
          {}
        ],
        "durationSeconds": 0,
        "endTime": "string",
        "additionalProperties": "object"
      },
      "conversationRoutingData": {
        "queue": "object",
        "language": "object",
        "priority": 0,
        "skills": [
          {}
        ],
        "scoredAgents": [
          {
            "agent": "object",
            "score": 0
          }
        ]
      },
      "peer": "string",
      "screenRecordingState": "string",
      "flaggedReason": "general",
      "journeyContext": {
        "customer": {
          "id": "string",
          "idType": "string"
        },
        "customerSession": {
          "id": "string",
          "type": "string"
        },
        "triggeringAction": {
          "id": "string",
          "actionMap": {
            "id": "string",
            "version": 0
          }
        }
      },
      "roomId": "string",
      "avatarImageUrl": "string"
    }
  ],
  "otherMediaUris": [
    {}
  ]
}
*/
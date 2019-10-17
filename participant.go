package purecloud

import (
	"time"
)

type Participant struct {
  ID                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	StartTime              time.Time               `json:"startTime"`
	ConnectedTime          time.Time               `json:"connectedTime"`
	EndTime                time.Time               `json:"endTime"`
	StartHoldTime          time.Time               `json:"startHoldTime"`
  Purpose                string                  `json:"purpose"`

  UserURI                string                  `json:"userUri"`
  UserID                 string                  `json:"userId"`
  ExternalContactID      string                  `json:"externalContactId"`
  ExternalOrganizationID string                  `json:"externalOrganizationId"`

  QueueID                string                  `json:"queueId"`
  GroupID                string                  `json:"groupId"`
  QueueName              string                  `json:"queueName"`
  ParticipantType        string                  `json:"participantType"`
  ConsultParticipantID   string                  `json:"consultParticipantId"`
  MonitoredParticipantID string                  `json:"monitoredParticipantId"`

  Address                string                  `json:"address"`
  ANI                    string                  `json:"ani"`
  ANIName                string                  `json:"aniName"`
  DNIS                   string                  `json:"dnis"`
  Locale                 string                  `json:"locale"`

  Attributes             map[string]string       `json:"attributes"`
  Calls                  []ConversationCall      `json:"calls"`
  Callbacks              []ConversationCallback  `json:"callbacks"`
  Chats                  []ConversationChat      `json:"chats"`
  CobrowseSessions       []CobrowseSession       `json:"cobrowseSession"`
  Emails                 []ConversationEmail     `json:"emails"`
  Messages               []ConversationMessage   `json:"messages"`
  ScreenShares           []ScreenShare           `json:"screenShares"`
  SocialExpressions      []SocialExpression      `json:"socialExpressions"`
  Videos                 []ConversationVideo     `json:"videos"`
  Evaluations            []Evaluation            `json:"evaluations"`

  WrapupRequired         bool                    `json:"wrapupRequired"`
  WrapupPrompt           string                  `json:"wrapupPrompt"`
  WrapupTimeout          int                     `json:"wrapupTimeoutMs` // time.Duration
  WrapupSkipped          bool                    `json:"wrapupSkipped"`
  Wrapup                 Wrapup                  `json:"wrapup"`

  AlertingTimeout        int                     `json:"alertingTimeoutMs` // time.Duration
  ScreenRecordingState   string                  `json:"screenRecordingState"`
  FlaggedReason          string                  `json:"flaggedReason"`

  RoutingData            ConversationRoutingData `json:"conversationRoutingData"`
	SelfURI                string                  `json:"selfUri"`

/*
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
*/
}

// GetID gets the identifier of this
//   implements Identifiable
func (participant Participant) GetID() string {
	return participant.ID
}
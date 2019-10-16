package purecloud

import (
	"time"
)

type ErrorInfo struct {
	Status int `json:"status"`
	Code   string `json:"code"`
	EntityID string `json:"entityId"`
	EntityName string `json:"entityName"`
	Message    string `json:"message"`
	MessageWithParams string `json:"messageWithParams"`
	MessageParams     interface{} `json:"messageParams"`
	ContextID         string      `json:"contextId"`
	Details           struct {
		ErrorCode  string `json:"errorCode"`
		Fieldname  string `json:"fieldName"`
		EntityID   string `json:"entityId"`
		EntityName string `json:"entityName"`
	} `json:"details"`
	Errors []ErrorInfo `json:"errors"`
}

type Participant struct {
  ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	StartTime              time.Time `json:"startTime"`
	ConnectedTime          time.Time `json:"connectedTime"`
	EndTime                time.Time `json:"endTime"`
	StartHoldTime          time.Time `json:"startHoldTime"`
  Purpose                string    `json:"purpose"`

  UserURI                string    `json:"userUri"`
  UserID                 string    `json:"userId"`
  ExternalContactID      string    `json:"externalContactId"`
  ExternalOrganizationID string    `json:"externalOrganizationId"`

  QueueID                string    `json:"queueId"`
  GroupID                string    `json:"groupId"`
  QueueName              string    `json:"queueName"`
  ParticipantType        string    `json:"participantType"`
  ConsultParticipantID   string    `json:"consultParticipantId"`

  Address                string    `json:"address"`
  ANI                    string    `json:"ani"`
  ANIName                string    `json:"aniName"`
  DNIS                   string    `json:"dnis"`
  Locale                 string    `json:"locale"`

  WrapupRequired         bool      `json:"wrapupRequired"`
  WrapupPrompt           string    `json:"wrapupPrompt"`
  WrapupTimeout          int       `json:"wrapupTimeoutMs` // time.Duration
  WrapupSkipped          bool      `json:"wrapupSkipped"`
  Wrapup                 Wrapup    `json:"wrapup"`



	State         string    `json:"state"` // alerting,dialing,contacting,offering,connected,disconnected,terminated,converting,uploading,transmitting,scheduled,none
	Direction     string    `json:"direction"` // inbound,outbound
	DisconnectType string   `json:"disconnectType"` // endpoint,client,system,transfer,timeout,transfer.conference,transfer.consult,transfer.forward,transfer.noanswer,transfer.notavailable,transport.failure,error,peer,other,spam,uncallable
	Held           bool     `json:"held"`
	User           struct {
		ID   string           `json:"id"`
		Name string           `json:"name"`
	}                       `json:"user"`
	Queue interface{} `json:"queue"`
	ErrorInfo ErrorInfo `json:"errorInfo"`
	Attributes interface{} `json:"attributes"`

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
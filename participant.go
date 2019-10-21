package purecloud

import (
	"strings"
	"encoding/json"
  "time"

	"github.com/pkg/errors"
)

type Participant struct {
  ID                     string                  `json:"id"`
  Name                   string                  `json:"name"`
  State                  string                  `json:"state"`
	StartTime              time.Time               `json:"startTime"`
	ConnectedTime          time.Time               `json:"connectedTime"`
	EndTime                time.Time               `json:"endTime"`
	StartHoldTime          time.Time               `json:"startHoldTime"`
  Purpose                string                  `json:"purpose"`

  User                   *User                   `json:"user"`
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
  Calls                  []*ConversationCall     `json:"calls"`
  Callbacks              []*ConversationCallback `json:"callbacks"`
  Chats                  []*ConversationChat     `json:"chats"`
  CobrowseSessions       []*CobrowseSession      `json:"cobrowseSession"`
  Emails                 []*ConversationEmail    `json:"emails"`
  Messages               []*ConversationMessage  `json:"messages"`
  ScreenShares           []*ScreenShare          `json:"screenShares"`
  SocialExpressions      []*SocialExpression     `json:"socialExpressions"`
  Videos                 []*ConversationVideo    `json:"videos"`
  Evaluations            []*Evaluation           `json:"evaluations"`

  WrapupRequired         bool                    `json:"wrapupRequired"`
  WrapupPrompt           string                  `json:"wrapupPrompt"`
  WrapupTimeout          int                     `json:"wrapupTimeoutMs"` // time.Duration
  WrapupSkipped          bool                    `json:"wrapupSkipped"`
  Wrapup                 *Wrapup                 `json:"wrapup"`

  AlertingTimeout        int                     `json:"alertingTimeoutMs"` // time.Duration
  ScreenRecordingState   string                  `json:"screenRecordingState"`
  FlaggedReason          string                  `json:"flaggedReason"`

  RoutingData            ConversationRoutingData `json:"conversationRoutingData"`
	SelfURI                string                  `json:"selfUri"`
}

func (participant Participant) IsMember(mediaType string, identifiable Identifiable) bool {
  switch strings.ToLower(mediaType) {
  case "call":
    for _, conversation := range participant.Calls {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "chat":
    for _, conversation := range participant.Chats {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "callback":
    for _, conversation := range participant.Callbacks {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "cobrowsesession", "cobrowse":
    for _, conversation := range participant.CobrowseSessions {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "email":
    for _, conversation := range participant.Emails {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "evaluation":
    for _, conversation := range participant.Evaluations {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "message":
    for _, conversation := range participant.Messages {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "screenshare":
    for _, conversation := range participant.ScreenShares {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "socialexpression":
    for _, conversation := range participant.SocialExpressions {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  case "video":
    for _, conversation := range participant.Videos {
      if identifiable.GetID() == conversation.ID {
        return true
      }
    }
  default:
    return false
  }
  return false
}

// GetID gets the identifier of this
//   implements Identifiable
func (participant Participant) GetID() string {
	return participant.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (participant Participant) String() string {
	if len(participant.Name) != 0 {
		return participant.Name
	}
	return participant.ID
}

// UnmarshalJSON unmarshals JSON into this
func (participant *Participant) UnmarshalJSON(payload []byte) (err error) {
  type surrogate Participant
  var inner struct {
    surrogate
    UserID    string `json:"userId"`
    UserURI   string `json:"userUri"`
  }

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
  }
  *participant = Participant(inner.surrogate)
  if participant.User == nil && len(inner.UserID) > 0 {
    participant.User = &User{ID: inner.UserID, SelfURI: inner.UserURI}
  }
  return
}
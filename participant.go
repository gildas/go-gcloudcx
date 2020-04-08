package purecloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gildas/go-errors"
)

// Participant describes a Chat Participant
type Participant struct {
	ID              string `json:"id"`
	SelfURI         string `json:"selfUri"`
	Name            string `json:"name"`
	ParticipantType string `json:"participantType"`
	State           string `json:"state"`
	Held            bool   `json:"held"`
	Direction       string `json:"direction"`

	StartTime      time.Time `json:"startTime"`
	ConnectedTime  time.Time `json:"connectedTime"`
	EndTime        time.Time `json:"endTime"`
	StartHoldTime  time.Time `json:"startHoldTime"`
	Purpose        string    `json:"purpose"`
	DisconnectType string    `json:"disconnectType"`

	User                   *User            `json:"user"`
	ExternalContact        *DomainEntityRef `json:"externalContact"`
	ExternalContactID      string           `json:"externalContactId"`
	ExternalOrganization   *DomainEntityRef `json:"externalOrganization"`
	ExternalOrganizationID string           `json:"externalOrganizationId"`

	Queue                  *Queue           `json:"queue"`
	QueueID                string           `json:"queueId"`
	GroupID                string           `json:"groupId"`
	QueueName              string           `json:"queueName"`
	ConsultParticipantID   string           `json:"consultParticipantId"`
	MonitoredParticipantID string           `json:"monitoredParticipantId"`
	Script                 *DomainEntityRef `json:"script"`

	Address string `json:"address"`
	ANI     string `json:"ani"`
	ANIName string `json:"aniName"`
	DNIS    string `json:"dnis"`
	Locale  string `json:"locale"`

	Attributes        map[string]string       `json:"attributes"`
	Calls             []*ConversationCall     `json:"calls"`
	Callbacks         []*ConversationCallback `json:"callbacks"`
	Chats             []*ConversationChat     `json:"chats"`
	CobrowseSessions  []*CobrowseSession      `json:"cobrowseSession"`
	Emails            []*ConversationEmail    `json:"emails"`
	Messages          []*ConversationMessage  `json:"messages"`
	ScreenShares      []*ScreenShare          `json:"screenShares"`
	SocialExpressions []*SocialExpression     `json:"socialExpressions"`
	Videos            []*ConversationVideo    `json:"videos"`
	Evaluations       []*Evaluation           `json:"evaluations"`

	WrapupRequired bool          `json:"wrapupRequired"`
	WrapupPrompt   string        `json:"wrapupPrompt"`
	WrapupTimeout  time.Duration `json:"-"`
	WrapupSkipped  bool          `json:"wrapupSkipped"`
	Wrapup         *Wrapup       `json:"wrapup"`

	AlertingTimeout      time.Duration           `json:"-"`
	ScreenRecordingState string                  `json:"screenRecordingState"`
	FlaggedReason        string                  `json:"flaggedReason"`
	Peer                 string                  `json:"peer"`
	RoutingData          ConversationRoutingData `json:"conversationRoutingData"`
	JourneyContext       *JourneyContext         `json:"journeyContext"`
	ErrorInfo            *ErrorBody              `json:"errorInfo"`
}

// IsMember tells if the Participant is a memmber of the Conversation (Identifiable)
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

// UpdateState updates the state of the Participant in target
func (participant *Participant) UpdateState(target StateUpdater, state string) error {
	return target.UpdateState(participant, state)
}

// MarshalJSON marshals this into JSON
func (participant Participant) MarshalJSON() ([]byte, error) {
	userID := ""
	userURI := ""
	if participant.User != nil {
		userID = participant.User.ID
		userURI = participant.User.SelfURI
	}
	type surrogate Participant
	data, err := json.Marshal(struct {
		surrogate
		UserID            string `json:"userId"`
		UserURI           string `json:"userUri"`
		AlertingTimeoutMs int64  `json:"alertingTimeoutMs"`
		WrapupTimeoutMs   int64  `json:"wrapupTimeoutMs"`
	}{
		surrogate:         surrogate(participant),
		UserID:            userID,
		UserURI:           userURI,
		AlertingTimeoutMs: int64(participant.AlertingTimeout.Milliseconds()),
		WrapupTimeoutMs:   int64(participant.WrapupTimeout.Milliseconds()),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (participant *Participant) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Participant
	var inner struct {
		surrogate
		UserID            string `json:"userId"`
		UserURI           string `json:"userUri"`
		AlertingTimeoutMs int64  `json:"alertingTimeoutMs"`
		WrapupTimeoutMs   int64  `json:"wrapupTimeoutMs"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*participant = Participant(inner.surrogate)
	if participant.User == nil && len(inner.UserID) > 0 {
		participant.User = &User{ID: inner.UserID, SelfURI: inner.UserURI}
	}
	participant.AlertingTimeout = time.Duration(inner.AlertingTimeoutMs) * time.Millisecond
	participant.WrapupTimeout = time.Duration(inner.WrapupTimeoutMs) * time.Millisecond
	return
}

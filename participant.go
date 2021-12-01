package gcloudcx

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// Participant describes a Chat Participant
type Participant struct {
	ID              uuid.UUID  `json:"id"`
	SelfURI         URI        `json:"selfUri"`
	Type            string     `json:"type"`
	Provider        string     `json:"provider"`
	Name            string     `json:"name"`
	ParticipantType string     `json:"participantType,omitempty"`
	State           string     `json:"state,omitempty"`
	Held            bool       `json:"held,omitempty"`
	Direction       string     `json:"direction,omitempty"`
	StartTime       time.Time  `json:"startTime,omitempty"`
	ConnectedTime   time.Time  `json:"connectedTime,omitempty"`
	EndTime         time.Time  `json:"endTime,omitempty"`
	StartHoldTime   time.Time  `json:"startHoldTime,omitempty"`
	Purpose         string     `json:"purpose"`
	DisconnectType  string     `json:"disconnectType,omitempty"`

	User                   *User            `json:"user,omitempty"`
	ExternalContact        *DomainEntityRef `json:"externalContact,omitempty"`
	ExternalContactID      string           `json:"externalContactId,omitempty"`
	ExternalOrganization   *DomainEntityRef `json:"externalOrganization,omitempty"`
	ExternalOrganizationID string           `json:"externalOrganizationId,omitempty"`

	Queue                  *Queue           `json:"queue,omitempty"`
	QueueID                string           `json:"queueId,omitempty"`
	GroupID                string           `json:"groupId,omitempty"`
	QueueName              string           `json:"queueName,omitempty"`
	ConsultParticipantID   string           `json:"consultParticipantId,omitempty"`
	MonitoredParticipantID string           `json:"monitoredParticipantId,omitempty"`
	Script                 *DomainEntityRef `json:"script,omitempty"`

	Address string  `json:"address,omitempty"`
	ANI     string  `json:"ani,omitempty"`
	ANIName string  `json:"aniName,omitempty"`
	DNIS    string  `json:"dnis,omitempty"`
	Locale  string  `json:"locale,omitempty"`
	From    Address `json:"fromAddress,omitempty"`
	To      Address `json:"toAddress,omitempty"`

	Attributes        map[string]string       `json:"attributes,omitempty"`
	Calls             []*ConversationCall     `json:"calls,omitempty"`
	Callbacks         []*ConversationCallback `json:"callbacks,omitempty"`
	Chats             []*ConversationChat     `json:"chats,omitempty"`
	CobrowseSessions  []*CobrowseSession      `json:"cobrowseSession,omitempty"`
	Emails            []*ConversationEmail    `json:"emails,omitempty"`
	Messages          []*ConversationMessage  `json:"messages,omitempty"`
	ScreenShares      []*ScreenShare          `json:"screenShares,omitempty"`
	SocialExpressions []*SocialExpression     `json:"socialExpressions,omitempty"`
	Videos            []*ConversationVideo    `json:"videos,omitempty"`
	Evaluations       []*Evaluation           `json:"evaluations,omitempty"`

	WrapupRequired bool          `json:"wrapupRequired"`
	WrapupPrompt   string        `json:"wrapupPrompt,omitempty"`
	WrapupTimeout  time.Duration `json:"-"`
	WrapupSkipped  bool          `json:"wrapupSkipped,omitempty"`
	Wrapup         *Wrapup       `json:"wrapup,omitempty"`

	AlertingTimeout      time.Duration           `json:"-"`
	ScreenRecordingState string                  `json:"screenRecordingState,omitempty"`
	FlaggedReason        string                  `json:"flaggedReason,omitempty"`
	Peer                 string                  `json:"peer,omitempty"`
	RoutingData          ConversationRoutingData `json:"conversationRoutingData,omitempty"`
	JourneyContext       *JourneyContext         `json:"journeyContext,omitempty"`
	ErrorInfo            *ErrorBody              `json:"errorInfo,omitempty"`
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
func (participant Participant) GetID() uuid.UUID {
	return participant.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (participant Participant) GetURI() URI {
	return participant.SelfURI
}

// String gets a string version
//   implements the fmt.Stringer interface
func (participant Participant) String() string {
	if len(participant.Name) != 0 {
		return participant.Name
	}
	return participant.ID.String()
}

// UpdateState updates the state of the Participant in target
func (participant *Participant) UpdateState(context context.Context, target StateUpdater, state string) error {
	return target.UpdateState(context, participant, state)
}

// MarshalJSON marshals this into JSON
func (participant Participant) MarshalJSON() ([]byte, error) {
	userID := uuid.Nil
	userURI := URI("")
	if participant.User != nil {
		userID = participant.User.ID
		userURI = participant.User.SelfURI
	}
	type surrogate Participant
	data, err := json.Marshal(struct {
		surrogate
		UserID            uuid.UUID `json:"userId"`
		UserURI           URI       `json:"userUri"`
		AlertingTimeoutMs int64     `json:"alertingTimeoutMs"`
		WrapupTimeoutMs   int64     `json:"wrapupTimeoutMs"`
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
		UserID            uuid.UUID `json:"userId"`
		UserURI           URI       `json:"userUri"`
		AlertingTimeoutMs int64     `json:"alertingTimeoutMs"`
		WrapupTimeoutMs   int64     `json:"wrapupTimeoutMs"`
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

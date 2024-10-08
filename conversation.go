package gcloudcx

import (
	"context"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Conversation contains the details of a live conversation
//
//	See: https://developer.mypurecloud.com/api/rest/v2/conversations
type Conversation struct {
	ID              uuid.UUID     `json:"id"`
	SelfURI         URI           `json:"selfUri,omitempty"`
	Name            string        `json:"name"`
	ExternalTag     string        `json:"externalTag,omitempty"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Address         string        `json:"address"`
	Participants    []Participant `json:"participants"`
	ConversationIDs []uuid.UUID   `json:"conversationIds"`
	MaxParticipants int           `json:"maxParticipants"`
	RecordingState  string        `json:"recordingState"`
	State           string        `json:"state"`
	Divisions       []struct {
		Division DomainEntityRef   `json:"division"`
		Entities []DomainEntityRef `json:"entities"`
	} `json:"divisions"`
	client *Client        `json:"-"`
	logger *logger.Logger `json:"-"`
}

// ConversationRoutingData  defines routing details of a Conversation
type ConversationRoutingData struct {
	Queue        AddressableEntityRef   `json:"queue"`
	Language     AddressableEntityRef   `json:"language"`
	Priority     int                    `json:"priority"`
	Skills       []AddressableEntityRef `json:"skills"`
	ScoredAgents []struct {
		Agent AddressableEntityRef `json:"agent"`
		Score int                  `json:"score"`
	} `json:"scoredAgents"`
}

// Segment describes a fragment of a Conversation
type Segment struct {
	Type           string    `json:"type"`
	DisconnectType string    `json:"disconnectType"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	HowEnded       string    `json:"howEnded"`
}

// DisconnectReason describes the reason of a disconnect
type DisconnectReason struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Phrase string `json:"phrase"`
}

// DialerPreview describes a Diapler Preview
type DialerPreview struct {
	ID                 uuid.UUID `json:"id"`
	ContactID          uuid.UUID `json:"contactId"`
	ContactListID      uuid.UUID `json:"contactListId"`
	CampaignID         uuid.UUID `json:"campaignId"`
	PhoneNumberColumns []struct {
		Type       string `json:"type"`
		ColumnName string `json:"columnName"`
	} `json:"phoneNumberColumns"`
}

// Voicemail describes a voicemail
type Voicemail struct {
	ID           string `json:"id"`
	UploadStatus string `json:"uploadStatus"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (conversation *Conversation) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			conversation.ID = parameter
		case *Client:
			conversation.client = parameter
		case *logger.Logger:
			conversation.logger = parameter.Child("conversation", "conversation", "id", conversation.ID)
		}
	}
	if conversation.logger == nil {
		conversation.logger = logger.Create("gclouccx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (conversation Conversation) GetID() uuid.UUID {
	return conversation.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (conversation Conversation) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/conversations/%s", ids[0])
	}
	if conversation.ID != uuid.Nil {
		return NewURI("/api/v2/conversations/%s", conversation.ID)
	}
	return URI("/api/v2/conversations/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (conversation Conversation) String() string {
	if len(conversation.Name) != 0 {
		return conversation.Name
	}
	return conversation.ID.String()
}

// Disconnect disconnect an Identifiable from this
//
// implements Disconnecter
func (conversation Conversation) Disconnect(context context.Context, identifiable Identifiable) error {
	return conversation.client.Patch(
		conversation.logger.ToContext(context),
		NewURI("/conversations/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: "disconnected"},
		nil,
	)
}

// UpdateState update the state of an identifiable in this
//
// implements StateUpdater
func (conversation Conversation) UpdateState(context context.Context, identifiable Identifiable, state string) error {
	return conversation.client.Patch(
		conversation.logger.ToContext(context),
		NewURI("/conversations/%s/participants/%s", conversation.ID, identifiable.GetID()),
		MediaParticipantRequest{State: state},
		nil,
	)
}

// GetParticipantByPurpose get the conversation's participant by its purpose
func (conversation Conversation) GetParticipantByPurpose(purpose string) (participant *Participant, found bool) {
	for _, current := range conversation.Participants {
		if current.Purpose == purpose {
			return &current, true
		}
	}
	return nil, false
}

// AssociateExternalContact associates an ExternalContact to this Conversation
func (conversation Conversation) AssociateExternalContact(context context.Context, contact *ExternalContact, communicationID uuid.UUID, mediaType string) error {
	if contact == nil {
		return errors.ArgumentMissing.With("contact")
	}
	return conversation.client.Put(
		context,
		NewURI("/externalcontacts/conversations/%s", conversation.ID),
		struct {
			MediaType         string `json:"mediaType"`
			ConversationID    string `json:"conversationId"`
			CommunicationID   string `json:"communicationId"`
			ExternalContactID string `json:"externalContactId"`
		}{
			MediaType:         strings.ToUpper(mediaType),
			ConversationID:    conversation.ID.String(),
			CommunicationID:   communicationID.String(),
			ExternalContactID: contact.ID.String(),
		},
		nil,
	)
}

package purecloud

import (
	"time"

	"github.com/gildas/go-logger"
)

// Conversation contains the details of a live conversation
//   See: https://developer.mypurecloud.com/api/rest/v2/conversations
type Conversation struct {
	ID              string          `json:"id"`
	SelfURI         string          `json:"selfUri,omitempty"`
	Name            string          `json:"name"`
	StartTime       time.Time       `json:"startTime"`
	EndTime         time.Time       `json:"endTime"`
	Address         string          `json:"address"`
	Participants    []*Participant  `json:"participants"`
	ConversationIDs []string        `json:"conversationIds"`
	MaxParticipants int             `json:"maxParticipants"`
	RecordingState  string          `json:"recordingState"`
	State           string          `json:"state"`
	Divisions       []struct {
		Division DomainEntityRef     `json:"division"`
		Entities []DomainEntityRef   `json:"entities"`
	}                               `json:"divisions"`

	Client          *Client         `json:"-"`
	Logger          *logger.Logger  `json:"-"`
}

// ConversationRoutingData  defines routing details of a Conversation
type ConversationRoutingData struct {
	Queue    AddressableEntityRef   `json:"queue"`
	Language AddressableEntityRef   `json:"language"`
	Priority int                    `json:"priority"`
	Skills   []AddressableEntityRef `json:"skills"`
	ScoredAgents []struct{
		Agent AddressableEntityRef    `json:"agent"`
		Score int                     `json:"score"`
	}                               `json:"scoredAgents"`
}

// Segment describes a fragment of a Conversation
type Segment struct {
	Type            string    `json:"type"`
	DisconnectType  string    `json:"disconnectType"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	HowEnded        string    `json:"howEnded"`
}

// DisconnectReason describes the reason of a disconnect
type DisconnectReason struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Phrase string `json:"phrase"`
}

// DialerPreview describes a Diapler Preview
type DialerPreview struct {
	ID                 string     `json:"id"`
	ContactID          string     `json:"contactId"`
	ContactListID      string     `json:"contactListId"`
	CaompaignID        string     `json:"campaignId"`
	PhoneNumberColumns []struct {
		Type       string `json:"type"`
		ColumnName string `json:"columnName"`
	}                             `json:"phoneNumberColumns"`
}

// Voicemail describes a voicemail
type Voicemail struct {
	ID           string `json:"id"`
	UploadStatus string `json:"uploadStatus"`

}

// Initialize initializes this from the given Client
//   implements Initializable
func (conversation *Conversation) Initialize(parameters ...interface{}) error {
	client, logger, err := ExtractClientAndLogger(parameters...)
	if err != nil {
		return err
	}
	if len(conversation.ID) > 0 {
		if err := conversation.Client.Get("/conversations/" + conversation.ID, &conversation); err != nil {
			return err
		}
	}
	conversation.Client = client
	conversation.Logger = logger.Topic("conversation").Scope("conversation").Record("media", "chat")
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (conversation Conversation) GetID() string {
	return conversation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (conversation Conversation) String() string {
	if len(conversation.Name) != 0 {
		return conversation.Name
	}
	return conversation.ID
}
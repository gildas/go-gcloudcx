package gcloudcx

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

type Recording struct {
	ID                           uuid.UUID                     `json:"id"`
	Name                         string                        `json:"name"`
	ConversationID               uuid.UUID                     `json:"conversationId"`
	Conversation                 *Conversation                 `json:"-"`
	SessionID                    uuid.UUID                     `json:"sessionId"`
	Path                         string                        `json:"path"`
	StartTime                    time.Time                     `json:"startTime"`
	EndTime                      time.Time                     `json:"endTime"`
	OriginalRecordingStartTime   time.Time                     `json:"originalRecordingStartTime"`
	Media                        string                        `json:"media"`        // audio, chat, messaging, email, screen
	MediaSubtype                 string                        `json:"mediaSubtype"` // Trunk, Station, Consult, Screen
	MediaSubject                 string                        `json:"mediaSubject"`
	Annotations                  []RecordingAnnotation         `json:"annotations"`
	Transcript                   []RecordingChatMessage        `json:"transcript"`
	EmailTranscript              []RecordingEmailMessage       `json:"emailTranscript"`
	MessagingTranscript          []RecordingMessagingMessage   `json:"messagingTranscript"`
	FileState                    string                        `json:"fileState"` // ARCHIVED, AVAILABLE, DELETED, RESTORED, RESTORING, UPLOADING, ERROR
	RestoreExpires               time.Time                     `json:"restoreExpirationTime"`
	MediaURIs                    map[string]*RecordingMediaURI `json:"mediaUris"`
	EstimatedTranscodeTime       time.Duration                 `json:"-"`
	ActualTranscodeTime          time.Duration                 `json:"-"`
	ArchiveMedium                string                        `json:"archiveMedium"` // CLOUDARCHIVE
	ArchiveDate                  time.Time                     `json:"archiveDate"`
	DeleteDate                   time.Time                     `json:"deleteDate"`
	ExportDate                   time.Time                     `json:"exportDate"`
	ExportedDate                 time.Time                     `json:"exportedDate"`
	OutputDuration               time.Duration                 `json:"-"`
	OutputSizeInBytes            int64                         `json:"outputSizeInBytes"`
	MaxAllowedRestorationsFxxx   int                           `json:"maxAllowedRestorations"` //TODO: Find the exact name of this field
	RemainingRestorationsAllowed int                           `json:"remainingRestorations"`  // TODO: Find the exact name of this field
	Users                        []User                        `json:"users"`
	RecordingFileRole            string                        `json:"recordingFileRole"`    // CUSTOMER_EXPERIENCE, ADHOC
	RecordingErrorStatus         string                        `json:"recordingErrorStatus"` // EMAIL_TRANSCRIPT_TOO_LARGE
	CreationTime                 time.Time                     `json:"creationTime"`
	ExternalContacts             []DomainEntityRef             `json:"externalContacts"`
	SelfURI                      string                        `json:"selfUri"`
	client                       *Client                       `json:"-"`
	logger                       *logger.Logger                `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger, *gcloudcx.Conversation
//
// implements Initializable
func (recording *Recording) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			recording.ID = parameter
		case *Client:
			recording.client = parameter
		case *Conversation:
			recording.Conversation = parameter
		case *logger.Logger:
			recording.logger = parameter.Child("recording", "recording", "id", recording.ID)
		}
	}
	if recording.logger == nil {
		recording.logger = logger.Create("gclouccx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (recording Recording) GetID() uuid.UUID {
	return recording.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (recording Recording) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 1 {
		return NewURI("/api/v2/conversations/%s/recordings/%s", ids[0], ids[1])
	}
	if recording.ID != uuid.Nil {
		return NewURI("/api/v2/conversations/%s/recordings/%s", recording.ConversationID, recording.ID)
	}
	return URI("/api/v2/conversations/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (recording Recording) String() string {
	if len(recording.Name) != 0 {
		return recording.Name
	}
	return recording.ID.String()
}

// UnmarshalJSON unmarshals the recording from JSON
//
// Implements json.Unmarshaler
func (recording *Recording) UnmarshalJSON(data []byte) error {
	type surrogate Recording
	var inner struct {
		surrogate
		ID                     core.UUID `json:"id"`
		ConversationID         core.UUID `json:"conversationId"`
		SessionID              core.UUID `json:"sessionId"`
		StartTime              core.Time `json:"startTime"`
		EndTime                core.Time `json:"endTime"`
		OriginalStartTime      core.Time `json:"originalRecordingStartTime"`
		CreationTime           core.Time `json:"creationTime"`
		RestoreExpiration      core.Time `json:"restoreExpirationTime"`
		ArchiveDate            core.Time `json:"archiveDate"`
		DeleteDate             core.Time `json:"deleteDate"`
		ExportDate             core.Time `json:"exportDate"`
		ExportedDate           core.Time `json:"exportedDate"`
		EstimatedTranscodeTime int64     `json:"estimatedTranscodeTimeMs"`
		ActualTranscodeTime    int64     `json:"actualTranscodeTimeMs"`
		OutputDuration         int64     `json:"outputDurationMs"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*recording = Recording(inner.surrogate)
	recording.ID = uuid.UUID(inner.ID)
	recording.ConversationID = uuid.UUID(inner.ConversationID)
	recording.SessionID = uuid.UUID(inner.SessionID)
	recording.StartTime = time.Time(inner.StartTime)
	recording.EndTime = time.Time(inner.EndTime)
	recording.OriginalRecordingStartTime = time.Time(inner.OriginalStartTime)
	recording.CreationTime = time.Time(inner.CreationTime)
	recording.RestoreExpires = time.Time(inner.RestoreExpiration)
	recording.ArchiveDate = time.Time(inner.ArchiveDate)
	recording.DeleteDate = time.Time(inner.DeleteDate)
	recording.ExportDate = time.Time(inner.ExportDate)
	recording.ExportedDate = time.Time(inner.ExportedDate)
	recording.EstimatedTranscodeTime = time.Duration(inner.EstimatedTranscodeTime) * time.Millisecond
	recording.ActualTranscodeTime = time.Duration(inner.ActualTranscodeTime) * time.Millisecond
	recording.OutputDuration = time.Duration(inner.OutputDuration) * time.Millisecond
	recording.FileState = strings.ToUpper(recording.FileState)
	recording.Media = strings.ToLower(recording.Media)

	return nil
}

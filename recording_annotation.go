package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type RecordingAnnotation struct {
	ID                uuid.UUID             `json:"id"`
	Name              string                `json:"name"`
	Type              string                `json:"type"`
	Description       string                `json:"description"`
	Reason            string                `json:"reason"`
	User              *User                 `json:"user,omitempty"`
	Annotations       []RecordingAnnotation `json:"annotations"`
	Location          time.Duration         `json:"location"` // Offset from start of recording
	Duration          time.Duration         `json:"-"`
	AbsoluteLocation  time.Duration         `json:"-"` // Offset from start of recording, after removing the cumulative duration of all pauses
	AbsoluteDuration  time.Duration         `json:"-"`
	RecordingLocation time.Duration         `json:"-"` // Offset from start of recording, adjusted for any recording cuts
	RecordingDuration time.Duration         `json:"-"` // Duration of annotation, adjusted for any recording cuts
	RealtimeLocation  time.Duration         `json:"-"` // Offset from start of recording, before removing the cumulative duration of all pauses before this annotation
	SelfURI           string                `json:"selfUri"`
}

// UnmarshalJSON unmarshals the annotation from JSON
//
// Implements json.Unmarshaler
func (annotation *RecordingAnnotation) UnmarshalJSON(data []byte) error {
	type surrogate RecordingAnnotation
	var inner struct {
		surrogate
		ID                core.UUID `json:"id"`
		Location          int64     `json:"location"` // Offset in milliseconds from start of recording
		Duration          int64     `json:"durationMs"`
		AbsoluteLocation  int64     `json:"absoluteLocation"` // Offset in milliseconds from start of recording, after removing the cumulative duration of all pauses
		AbsoluteDuration  int64     `json:"absoluteDurationMs"`
		RecordingLocation int64     `json:"recordingLocation"`   // Offset in milliseconds from start of recording, adjusted for any recording cuts
		RecordingDuration int64     `json:"recordingDurationMs"` // Duration of annotation in Ms, adjusted for any recording cuts
		RealtimeLocation  int64     `json:"realtimeLocation"`    // Offset in milliseconds from start of recording, before removing the cumulative duration of all pauses before this annotation
	}
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*annotation = RecordingAnnotation(inner.surrogate)
	annotation.ID = uuid.UUID(inner.ID)
	annotation.Location = time.Duration(inner.Location) * time.Millisecond
	annotation.Duration = time.Duration(inner.Duration) * time.Millisecond
	annotation.AbsoluteLocation = time.Duration(inner.AbsoluteLocation) * time.Millisecond
	annotation.AbsoluteDuration = time.Duration(inner.AbsoluteDuration) * time.Millisecond
	annotation.RecordingLocation = time.Duration(inner.RecordingLocation) * time.Millisecond
	annotation.RecordingDuration = time.Duration(inner.RecordingDuration) * time.Millisecond
	annotation.RealtimeLocation = time.Duration(inner.RealtimeLocation) * time.Millisecond
	return nil
}

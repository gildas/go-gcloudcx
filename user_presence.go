package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

// UserPresence  describes the Presence of a User
type UserPresence struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Source       string              `json:"source"`
	Primary      bool                `json:"primary"`
	Definition   *PresenceDefinition `json:"presenceDefinition"`
	Message      string              `json:"message"`
	ModifiedDate time.Time           `json:"modifiedDate"`
	SelfURI      string              `json:"selfUri"`
}

// PresenceDefinition  defines Presence
type PresenceDefinition struct {
	ID             string `json:"id"`
	SystemPresence string `json:"systemPresence"`
	SelfURI        string `json:"selfUri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (definition PresenceDefinition) GetID() string {
	return definition.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (definition PresenceDefinition) String() string {
	if len(definition.SystemPresence) > 0 {
		return definition.SystemPresence
	}
	return definition.ID
}

// GetID gets the identifier of this
//   implements Identifiable
func (presence UserPresence) GetID() string {
	return presence.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (presence UserPresence) String() string {
	if len(presence.Name) > 0 {
		return presence.Name
	}
	if presence.Definition != nil {
		return presence.Definition.String()
	}
	return presence.Message
}

// UnmarshalJSON unmarshals JSON into this
func (presence *UserPresence) UnmarshalJSON(payload []byte) (err error) {
	type surrogate UserPresence
	var inner struct {
		surrogate
		PresenceMessage string `json:"presenceMessage"` // found in the UserActivityTopic
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*presence = UserPresence(inner.surrogate)
	if len(inner.PresenceMessage) > 0 {
		presence.Message = inner.PresenceMessage
	}
	return
}

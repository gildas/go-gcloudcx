package gcloudcx

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ChatMember describes a Chat Member
type ChatMember struct {
	ID            uuid.UUID         `json:"id,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	AvatarURL     *url.URL          `json:"-"`
	Role          string            `json:"role,omitempty"`
	State         string            `json:"state,omitempty"`
	JoinedAt      time.Time         `json:"joinDate,omitempty"`
	LeftAt        time.Time         `json:"leaveDate,omitempty"`
	Authenticated bool              `json:"authenticatedGuest,omitempty"`
	Custom        map[string]string `json:"customFields,omitempty"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (member ChatMember) GetID() uuid.UUID {
	return member.ID
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (member ChatMember) String() string {
	if len(member.DisplayName) != 0 {
		return member.DisplayName
	}
	return member.ID.String()
}

// MarshalJSON marshals this into JSON
func (member ChatMember) MarshalJSON() ([]byte, error) {
	type surrogate ChatMember
	data, err := json.Marshal(struct {
		surrogate
		A *core.URL `json:"avatarImageUrl,omitempty"`
	}{
		surrogate: surrogate(member),
		A:         (*core.URL)(member.AvatarURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (member *ChatMember) UnmarshalJSON(payload []byte) (err error) {
	type surrogate ChatMember
	var inner struct {
		surrogate
		A *core.URL `json:"avatarImageUrl,omitempty"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*member = ChatMember(inner.surrogate)
	member.AvatarURL = (*url.URL)(inner.A)
	return
}

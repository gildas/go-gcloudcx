package purecloud

import (
	"time"
)

type ChatMember struct {
	ID            string            `json:"id,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	AvatarURL     string            `json:"avatarImageUrl,omitempty"`
	Role          string            `json:"role,omitempty"`
	State         string            `json:"state,omitempty"`
	JoinedAt      time.Time         `json:"joinDate,omitempty"`
	LeftAt        time.Time         `json:"leaveDate,omitempty"`
	Authenticated bool              `json:"authenticatedGuest,omitempty"`
	Custom        map[string]string `json:"customFields,omitempty"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (member ChatMember) GetID() string {
	return member.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (member ChatMember) String() string {
	if len(member.DisplayName) != 0 {
		return member.DisplayName
	}
	return member.ID
}
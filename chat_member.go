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
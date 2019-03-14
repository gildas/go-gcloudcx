package conversation

// Conversation contains the details of a live chat/conversation
type Conversation struct {
	ID          string     `json:"id"`
	JWT         string     `json:"jwt"`
	EventStream string     `json:"eventStreamUri"`
	Member      ChatMember `json:"member"`
}

// ChatMember describes a chat guest member
type ChatMember struct {
	ID          string            `json:"id,omitifempty"`
	DisplayName string            `json:"displayName,omitifempty"`
	ImageURL    string            `json:"avatarImageUrl,omitifempty"`
	Custom      map[string]string `json:"customFields,omitifempty"`
}

// Target describes the target of a Chat/Conversation
type Target struct {
	Type    string `json:"targetType"`
	Address string `json:"targetAddress"`
}
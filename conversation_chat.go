package purecloud

// ConversationChat contains the details of a live chat/conversation
type ConversationChat struct {
	ID          string `json:"id,omitempty"`
	JWT         string `json:"jwt,omitempty"`
	EventStream string `json:"eventStreamUri,omitempty"`
	Guest       *ChatMember `json:"member,omitempty"`
	Members     map[string]*ChatMember `json:"-"`
	SelfURI     string `json:"selfUri,omitempty"`

	Client      *Client  `json:"-"`
	//Socket      *websocket.Conn    `json:"-"`
	//Logger      *logger.Logger     `json:"-"`
}
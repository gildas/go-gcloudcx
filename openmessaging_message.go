package purecloud

type OpenMessage struct {
	ID              string                `json:"id"`
	Channel         *OpenMessageChannel   `json:"channel"`
	Direction       string                `json:"direction"`
	Type            string                `json:"type"` // Text, Structured, Receipt
	Text            string                `json:"text"`
	Content         []*OpenMessageContent `json:"content,omitempty"`
	RelatedMessages []*OpenMessage        `json:"relatedMessages,omitempty"`
	Metadata        map[string]string     `json:"metadata,omitempty"`
}

type OpenMessageResult struct {
	ID             string              `json:"id"`
	Channel        *OpenMessageChannel `json:"channel"`
	Type           string              `json:"type"` // Text, Structured, Receipt
	Text           string              `json:"text"`
	Direction      string              `json:"direction"`
	Content        *OpenMessageContent `json:"content,omitempty"`
	Status         string              `json:"status,omitempty"`
	Reasons        []*StatusReason     `json:"reasons,omitempty"`
	Entity         string              `json:"originatingEntity"`
	IsFinalReceipt bool                `json:"isFinalReceipt"`
	Metadata       map[string]string   `json:"metadata"`
}

type OpenMessageContent struct {
	Type string `json:"contentType"` // Attachment, Location, QuickReply, ButtonResponse, Notification, GenericTemplate, ListTemplate, Postback, Reactions, Mention
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

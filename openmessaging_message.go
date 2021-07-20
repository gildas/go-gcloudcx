package gcloudcx

type OpenMessage struct {
	ID              string                `json:"id,omitempty"`
	Channel         *OpenMessageChannel   `json:"channel"`
	Direction       string                `json:"direction"`
	Type            string                `json:"type"` // Text, Structured, Receipt
	Text            string                `json:"text"`
	Content         []*OpenMessageContent `json:"content,omitempty"`
	RelatedMessages []*OpenMessage        `json:"relatedMessages,omitempty"`
	Reasons         []*StatusReason       `json:"reasons,omitempty"`
}

type OpenMessageResult struct {
	OpenMessage
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

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

// Redact redacts sensitive data
//
// implements logger.Redactable
func (message OpenMessage) Redact() interface{} {
	redacted := message
	if message.Channel != nil {
		redacted.Channel = message.Channel.Redact().(*OpenMessageChannel)
	}
	return &redacted
}

type OpenMessageResult struct {
	OpenMessage
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (result OpenMessageResult) Redact() interface{} {
	redacted := result
	if result.Channel != nil {
		redacted.Channel = result.Channel.Redact().(*OpenMessageChannel)
	}
	return &redacted
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

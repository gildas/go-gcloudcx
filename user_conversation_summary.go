package gcloudcx

// UserConversationSummary  describes the summary of a User's conversations
type UserConversationSummary struct {
	UserID           string       `json:"userId"`
	Call             MediaSummary `json:"call"`
	Callback         MediaSummary `json:"callback"`
	Email            MediaSummary `json:"email"`
	Message          MediaSummary `json:"message"`
	Chat             MediaSummary `json:"chat"`
	SocialExpression MediaSummary `json:"socialExpression"`
	Video            MediaSummary `json:"video"`
}

// MediaSummary describes a Media summary
type MediaSummary struct {
	ContactCenter MediaSummaryDetail `json:"contactCenter"`
	Enterprise    MediaSummaryDetail `json:"enterprise"`
}

// MediaSummaryDetail describes the details about a MediaSummary
type MediaSummaryDetail struct {
	Active int `json:"active"`
	ACW    int `json:"acw"`
}

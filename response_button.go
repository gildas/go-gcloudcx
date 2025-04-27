package gcloudcx

type ButtonResponse struct {
	Type        string `json:"type"` // Button, QuickReply, DatePicker
	Text        string `json:"text"`
	Payload     string `json:"payload"`
	MessageType string `json:"messageType"` // QuickReply, Card, Carousel
}

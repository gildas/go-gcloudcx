package main

type ChatMessage struct {
	ID         string `json:"id"`
	RequestID  string `json:"reqid"`
	UserID     string `json:"userId"`
	TrackingID string `json:"trackingId"`
	Content    string `json:"text"`
	Chat       *Chat  `json:"-"`
}

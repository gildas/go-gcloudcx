package gcloudcx

import "time"

// MessageDetails  describes details of a Message in a Message Conversation
type MessageDetails struct {
	ID           string           `json:"messageId"`
	Status       string           `json:"messageStatus"`
	SegmentCount int              `json:"messageSegmentCount"`
	Time         time.Time        `json:"messageTime"`
	Media        []MessageMedia   `json:"media"`
	Stickers     []MessageSticker `json:"stickers"`
}

// MessageMedia  describes the Media of a Message
type MessageMedia struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	MediaType     string `json:"mediaType"`
	ContentLength int64  `json:"contentLengthBytes"`
}

// MessageSticker  describes a Message Sticker
type MessageSticker struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

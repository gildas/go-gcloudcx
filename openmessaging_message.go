package purecloud

import (
	"time"
)

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

type OpenMessageChannel struct {
	Platform  string           `json:"platform"` // Open
	Type      string           `json:"type"` // Private, Public
	MessageID string           `json:"messageId"`
	Time      time.Time        `json:"time"`
	To        *OpenMessageTo   `json:"to"`
	From      *OpenMessageFrom `json:"from"`

}

type OpenMessageTo struct {
	ID string `json:"id"`
}

type OpenMessageFrom struct {
	ID string `json:"id"`
	Type string `json:"idType"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Nickname  string `json:"nickname"`
	ImageURL  string `json:"image"`     // TODO should be an url.URL
}

type OpenMessageContent struct {
	Type string `json:"contentType"` // Attachment, Location, QuickReply, ButtonResponse, Notification, GenericTemplate, ListTemplate, Postback, Reactions, Mention
}

type StatusReason struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
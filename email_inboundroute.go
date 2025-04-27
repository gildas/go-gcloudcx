package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type EmailInboundRoute struct {
	ID                   uuid.UUID         `json:"id"`
	Name                 string            `json:"name"`
	Pattern              string            `json:"pattern"`
	Queue                DomainEntityRef   `json:"queue"`
	Priority             int               `json:"priority"`
	Skills               []DomainEntityRef `json:"skills"`
	Language             DomainEntityRef   `json:"language"`
	FromName             string            `json:"fromName"`
	FromEmail            string            `json:"fromEmail"`
	Flow                 DomainEntityRef   `json:"flow"`
	ReplyEmailAddress    QueueEmailAddress `json:"replyEmailAddress"`
	AutoBCC              []EmailAddress    `json:"autoBcc"`
	SpamFlow             DomainEntityRef   `json:"spamFlow"`
	Signature            EmailSignature    `json:"signature"`
	HistoryInclusion     string            `json:"historyInclusion"` // Include, Exclude, Optional
	AllowMultipleActions bool              `json:"allowMultipleActions"`
	SelfURI              string            `json:"selfUri"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (route EmailInboundRoute) GetID() uuid.UUID {
	return route.ID
}

// GetURI gets the URI of this
//
//	implements Addressable
func (route EmailInboundRoute) GetURI() URI {
	return URI(route.SelfURI)
}

// GetName gets the name of this
//
//	implements Named
func (route EmailInboundRoute) GetName() string {
	return route.Name
}

// UnmarshalJSON unmarshals the email inbound route from JSON
//
// Implements json.Unmarshaler
func (route *EmailInboundRoute) UnmarshalJSON(data []byte) error {
	type surrogate EmailInboundRoute
	var inner struct {
		surrogate
		ID core.UUID `json:"id"`
	}
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*route = EmailInboundRoute(inner.surrogate)
	route.ID = uuid.UUID(inner.ID)
	return nil
}

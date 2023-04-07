package gcloudcx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

type RoutingMessageRecipient struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	MessengerType string         `json:"messengerType"`
	Flow          *Flow          `json:"flow"`
	DateCreated   time.Time      `json:"dateCreated,omitempty"`
	CreatedBy     *User          `json:"createdBy,omitempty"`
	DateModified  time.Time      `json:"dateModified,omitempty"`
	ModifiedBy    *User          `json:"modifiedBy,omitempty"`
	client        *Client        `json:"-"`
	logger        *logger.Logger `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (recipient *RoutingMessageRecipient) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case *Client:
			recipient.client = parameter
		case *logger.Logger:
			recipient.logger = parameter.Child("routingmessagerecipient", "routingmessagerecipient", "id", recipient.ID)
		}
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (recipient RoutingMessageRecipient) GetID() uuid.UUID {
	return recipient.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (recipient RoutingMessageRecipient) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/routing/message/recipients/%s", ids[0])
	}
	if recipient.ID != uuid.Nil {
		return NewURI("/api/v2/routing/message/recipients/%s", recipient.ID)
	}
	return URI("/api/v2/routing/message/recipients/")
}

// UpdateFlow updates the flow of this recipient
func (recipient *RoutingMessageRecipient) UpdateFlow(context context.Context, flow *Flow) error {
	// Should we disconnect the current flow first?
	type FlowRef struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	recipient.Flow = flow
	return recipient.client.Put(
		context,
		recipient.GetURI(),
		struct {
			ID            string  `json:"id"`
			Name          string  `json:"name"`
			MessengerType string  `json:"messengerType"`
			Flow          FlowRef `json:"flow"`
		}{
			ID:            recipient.ID.String(),
			Name:          recipient.Name,
			MessengerType: recipient.MessengerType,
			Flow:          FlowRef{ID: recipient.Flow.ID.String(), Name: recipient.Flow.Name},
		},
		nil,
	)
}

// DeleteFlow removes the flow of this recipient
func (recipient *RoutingMessageRecipient) DeleteFlow(context context.Context) error {
	recipient.Flow = nil
	return recipient.client.Put(context, recipient.GetURI(), recipient, nil)
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (recipient RoutingMessageRecipient) MarshalJSON() ([]byte, error) {
	type surrogate RoutingMessageRecipient
	data, err := json.Marshal(&struct {
		surrogate
		SelfURI URI `json:"selfUri"`
	}{
		surrogate: surrogate(recipient),
		SelfURI:   recipient.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

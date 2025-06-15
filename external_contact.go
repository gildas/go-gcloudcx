package gcloudcx

import (
	"context"
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// ExternalContact represents an external contact
//
// See: https://developer.genesys.cloud/devapps/api-explorer#post-api-v2-externalcontacts-contacts
type ExternalContact struct {
	ID      uuid.UUID      `json:"id"`
	SelfURI URI            `json:"selfUri,omitempty"`
	Name    string         `json:"name"`
	client  *Client        `json:"-"`
	logger  *logger.Logger `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (contact *ExternalContact) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			contact.ID = parameter
		case *Client:
			contact.client = parameter
		case *logger.Logger:
			contact.logger = parameter.Child("conversation", "conversation", "id", contact.ID)
		}
	}
	if contact.logger == nil {
		contact.logger = logger.Create("gclouccx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (contact ExternalContact) GetID() uuid.UUID {
	return contact.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (contact ExternalContact) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/externalcontacts/contacts/%s", ids[0])
	}
	if contact.ID != uuid.Nil {
		return NewURI("/api/v2/externalcontacts/contacts/%s", contact.ID)
	}
	return URI("/api/v2/externalcontacts/contacts/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (contact ExternalContact) String() string {
	if len(contact.Name) != 0 {
		return contact.Name
	}
	return contact.ID.String()
}

// SearchExternalContact search for an external contact by one of its attributes
func (client *Client) SearchExternalContact(context context.Context, value string) (contact *ExternalContact, correlationID string, err error) {
	entities, correlationID, err := client.FetchEntities(context, NewURI("/externalcontacts/contacts").WithQuery(Query{"q": value}))
	if err != nil {
		return nil, correlationID, err
	}
	if len(entities) == 0 {
		return nil, correlationID, errors.NotFound.With("value", value)
	}
	if err = json.Unmarshal(entities[0], &contact); err != nil {
		return nil, correlationID, err
	}
	return
}

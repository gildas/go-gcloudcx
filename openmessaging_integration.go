package purecloud

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

type OpenMessagingIntegration struct {
	ID               uuid.UUID             `json:"id"`
	Name             string                `json:"name"`
	WebhookURL       *url.URL              `json:"-"`
	WebhookToken     string                `json:"outboundNotificationWebhookSignatureSecretToken"`
	Recipient        *DomainEntityRef      `json:"recipient,omitempty"`
	SupportedContent *AddressableEntityRef `json:"supportedContent,omitempty"`
	DateCreated      time.Time             `json:"dateCreated,omitempty"`
	CreatedBy        *DomainEntityRef      `json:"createdBy,omitempty"`
	DateModified     time.Time             `json:"dateModified,omitempty"`
	ModifiedBy       *DomainEntityRef      `json:"modifiedBy,omitempty"`
	CreateStatus     string                `json:"createStatus,omitempty"`
	CreateError      *ErrorBody            `json:"createError,omitempty"`
	SelfURI          URI                   `json:"selfUri,omitempty"`
	Client           *Client               `json:"-"`
	Logger           *logger.Logger        `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   if the parameters contain a uuid.UUID, the corresponding integration is fetched
func (integration *OpenMessagingIntegration) Initialize(parameters ...interface{}) error {
	client, logger, id, err := parseParameters(parameters...)
	if err != nil {
		return err
	}
	if id != uuid.Nil {
		if err := client.Get(NewURI("/conversations/messaging/integrations/open/%s", id), &integration); err != nil {
			return err
		}
	}
	integration.Client = client
	integration.Logger = logger
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (integration OpenMessagingIntegration) GetID() uuid.UUID {
	return integration.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (integration OpenMessagingIntegration) String() string {
	if len(integration.Name) > 0 {
		return integration.Name
	}
	return integration.ID.String()
}

// FetchOpenMessagingIntegrations Fetches all OpenMessagingIntegration object
func FetchOpenMessagingIntegrations(parameters ...interface{}) ([]*OpenMessagingIntegration, error) {
	client, logger, _, err := parseParameters(parameters...)
	if err != nil {
		return nil, err
	}
	response := struct{
		Integrations []*OpenMessagingIntegration `json:"entities"`
		PageSize     int                         `json:"pageSize"`
		PageNumber   int                         `json:"pageNumber"`
		PageCount    int                         `json:"pageCount"`
		PageTotal    int                         `json:"total"`
		FirstURI     string                      `json:"firstUri"`
		SelfURI      string                      `json:"selfUri"`
		LastURI      string                      `json:"lastUri"`
	}{}
	if err = client.Get("/conversations/messaging/integrations/open", &response); err != nil {
		return nil, err
	}
	logger.Record("response", response).Infof("Got a response")
	for _, integration := range response.Integrations {
		integration.Client = client
		integration.Logger = logger.Child("openmessagingintegration", "openmessagingintegration", "openmessgingintegration", integration.ID)
	}
	return response.Integrations, nil
}

func (integration *OpenMessagingIntegration) Create(name string, webhookURL *url.URL, token string) error {
	response := &OpenMessagingIntegration{}
	err := integration.Client.Post(
		"/conversations/messaging/integrations/open",
		struct {
			Name    string `json:"name"`
			Webhook string `json:"outboundNotificationWebhookUrl"`
			Token   string `json:"outboundNotificationWebhookSignatureSecretToken"`
		}{
			Name:    name,
			Webhook: webhookURL.String(),
			Token:   token,
		},
		&response,
	)
	if err != nil {
		return err
	}
	integration.Logger.Record("response", response).Debugf("Created integration %#v", response)
	return nil
}

func (integration *OpenMessagingIntegration) Delete() error {
	return integration.Client.Delete(NewURI("/conversations/messaging/integrations/open/%s", integration.ID), nil)
}

// MarshalJSON marshals this into JSON
func (integration OpenMessagingIntegration) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessagingIntegration
	data, err := json.Marshal(struct {
		surrogate
		W *core.URL `json:"outboundNotificationWebhookUrl"`
	}{
		surrogate: surrogate(integration),
		W:         (*core.URL)(integration.WebhookURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (integration *OpenMessagingIntegration) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessagingIntegration
	var inner struct {
		surrogate
		W *core.URL `json:"outboundNotificationWebhookUrl"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*integration = OpenMessagingIntegration(inner.surrogate)
	integration.WebhookURL = (*url.URL)(inner.W)
	return
}

// TODO: There is also a PATCH method... we might want to provide some func
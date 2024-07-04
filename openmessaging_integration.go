package gcloudcx

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// OpenMessagingIntegration  describes an GCloud OpenMessaging Integration
//
// See https://developer.genesys.cloud/api/digital/openmessaging
type OpenMessagingIntegration struct {
	ID               uuid.UUID             `json:"id"`
	Name             string                `json:"name"`
	WebhookURL       *url.URL              `json:"-"`
	WebhookToken     string                `json:"outboundNotificationWebhookSignatureSecretToken"`
	WebhookHeaders   map[string]string     `json:"webhookHeaders,omitempty"`
	Recipient        *DomainEntityRef      `json:"recipient,omitempty"`
	SupportedContent *AddressableEntityRef `json:"supportedContent,omitempty"`
	DateCreated      time.Time             `json:"dateCreated,omitempty"`
	CreatedBy        *DomainEntityRef      `json:"createdBy,omitempty"`
	DateModified     time.Time             `json:"dateModified,omitempty"`
	ModifiedBy       *DomainEntityRef      `json:"modifiedBy,omitempty"`
	CreateStatus     string                `json:"createStatus,omitempty"` // Initiated, Completed, Error
	CreateError      *ErrorBody            `json:"createError,omitempty"`
	Status           string                `json:"status,omitempty"` // Active, Inactive
	Client           *Client               `json:"-"`
	logger           *logger.Logger        `json:"-"`
}

// IsCreated tells if this OpenMessagingIntegration has been created successfully
func (integration OpenMessagingIntegration) IsCreated() bool {
	return integration.CreateStatus == "Completed"
}

// IsCreated tells if this OpenMessagingIntegration has not been created successfully
func (integration OpenMessagingIntegration) IsError() bool {
	return integration.CreateStatus == "Error"
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (integration *OpenMessagingIntegration) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			integration.ID = parameter
		case *Client:
			integration.Client = parameter
		case *logger.Logger:
			integration.logger = parameter.Child("integration", "integration", "id", integration.ID)
		}
	}
	if integration.logger == nil {
		integration.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (integration OpenMessagingIntegration) GetID() uuid.UUID {
	return integration.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (integration OpenMessagingIntegration) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/conversations/messaging/integrations/open/%s", ids[0])
	}
	if integration.ID != uuid.Nil {
		return NewURI("/api/v2/conversations/messaging/integrations/open/%s", integration.ID)
	}
	return URI("/api/v2/conversations/messaging/integrations/open/")
}

// Create creates a new OpenMessaging Integration
func (client *Client) CreateOpenMessagingIntegration(context context.Context, name string, webhookURL *url.URL, token string, headers map[string]string) (*OpenMessagingIntegration, error) {
	integration := OpenMessagingIntegration{}
	err := client.Post(
		context,
		"/conversations/messaging/integrations/open",
		struct {
			Name    string            `json:"name"`
			Webhook string            `json:"outboundNotificationWebhookUrl"`
			Token   string            `json:"outboundNotificationWebhookSignatureSecretToken"`
			Headers map[string]string `json:"webhookHeaders,omitempty"`
		}{
			Name:    name,
			Webhook: webhookURL.String(),
			Token:   token,
			Headers: headers,
		},
		&integration,
	)
	if err != nil {
		return nil, err
	}
	integration.Client = client
	integration.logger = client.Logger.Child("openmessagingintegration", "openmessagingintegration", "id", integration.ID)
	return &integration, nil
}

// Delete deletes an OpenMessaging Integration
//
// If the integration was not created, nothing is done
func (integration *OpenMessagingIntegration) Delete(context context.Context) error {
	if integration.ID == uuid.Nil {
		return nil
	}
	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID)
	return integration.Client.Delete(
		log.ToContext(context),
		NewURI("/conversations/messaging/integrations/open/%s", integration.ID),
		nil,
	)
}

func (integration *OpenMessagingIntegration) Refresh(ctx context.Context) error {
	var value OpenMessagingIntegration
	if err := integration.Client.Get(ctx, integration.GetURI(), &value); err != nil {
		return err
	}
	integration.Name = value.Name
	integration.CreateStatus = value.CreateStatus
	integration.CreateError = value.CreateError
	integration.WebhookURL = value.WebhookURL
	integration.WebhookToken = value.WebhookToken
	integration.WebhookHeaders = value.WebhookHeaders
	integration.Recipient = value.Recipient
	integration.SupportedContent = value.SupportedContent
	integration.DateModified = value.DateModified
	integration.ModifiedBy = value.ModifiedBy
	return nil
}

// Update updates an OpenMessaging Integration
//
// If the integration was not created, an error is return without reaching GENESYS Cloud
func (integration *OpenMessagingIntegration) Update(context context.Context, name string, webhookURL *url.URL, token string) error {
	if integration.ID == uuid.Nil {
		return errors.ArgumentMissing.With("ID")
	}
	if webhookURL == nil {
		return errors.ArgumentMissing.With("webhookURL")
	}
	response := &OpenMessagingIntegration{}
	err := integration.Client.Patch(
		integration.logger.ToContext(context),
		NewURI("/conversations/messaging/integrations/open/%s", integration.ID),
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
		return errors.CreationFailed.Wrap(err)
	}
	integration.logger.Record("response", response).Debugf("Updated integration")
	return nil
}

// GetRoutingMessageRecipient fetches the RoutingMessageRecipient for this OpenMessagingIntegration
func (integration *OpenMessagingIntegration) GetRoutingMessageRecipient(context context.Context) (*RoutingMessageRecipient, error) {
	if integration == nil || integration.ID == uuid.Nil {
		return nil, errors.ArgumentMissing.With("ID")
	}
	if !integration.IsCreated() {
		return nil, errors.CreationFailed.With("integration", integration.ID)
	}
	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID)
	return Fetch[RoutingMessageRecipient](log.ToContext(context), integration.Client, integration)
}

// SendInboundTextMessage sends an Open Message text message from the middleware to GENESYS Cloud
//
// See https://developer.genesys.cloud/api/digital/openmessaging/inboundMessages#send-an-inbound-open-message
func (integration *OpenMessagingIntegration) SendInboundTextMessage(context context.Context, message OpenMessageText) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	if len(message.Channel.ID) == 0 {
		return "", errors.ArgumentMissing.With("channel.ID")
	}
	if len(message.Channel.MessageID) == 0 {
		return "", errors.ArgumentMissing.With("channel.MessageID")
	}
	message.Channel.Platform = "Open"
	message.Channel.Type = "Private"
	message.Channel.Time = time.Now().UTC()
	message.Channel.To = &OpenMessageTo{ID: integration.ID.String()}
	if err := message.Channel.Validate(); err != nil {
		return "", err
	}
	if len(message.Text) == 0 && len(message.Content) == 0 {
		return "", errors.ArgumentMissing.With("text")
	}
	message.Direction = "Inbound"
	// TODO: attributes and metadata should be of a new type Metadata that containd a map and a []string for keysToRedact

	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID, "message", message.GetID())
	result := OpenMessageText{}
	err = integration.Client.Post(
		log.ToContext(context),
		NewURI("/conversations/messages/%s/inbound/open/message", integration.ID),
		message,
		&result,
	)
	return result.ID, err
}

// SendInboundReceipt sends a receipt from the middleware to GENESYS Cloud
//
// Valid status values are: Delivered, Failed.
//
// Genesys Cloud will return a receipt from this request. If the returned receipt has a Failed status, the return error contains the reason(s) for the failure.
//
// See https://developer.genesys.cloud/commdigital/digital/openmessaging/inboundReceiptMessages
func (integration *OpenMessagingIntegration) SendInboundReceipt(context context.Context, receipt OpenMessageReceipt) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	if len(receipt.ID) == 0 {
		// if the messageID was provided in the Channel, we need to move it to the receipt
		receipt.ID = receipt.Channel.MessageID
		if len(receipt.ID) == 0 {
			return "", errors.ArgumentMissing.With("ID")
		}
	}
	receipt.Channel.MessageID = ""
	receipt.Direction = "Outbound"
	if len(receipt.Channel.ID) == 0 {
		return "", errors.ArgumentMissing.With("channel.ID")
	}
	receipt.Channel.Platform = "Open"
	receipt.Channel.Type = "Private"
	receipt.Channel.Time = time.Now().UTC()
	receipt.Channel.From = &OpenMessageFrom{ID: integration.ID.String(), Type: "email"}
	if err := receipt.Channel.Validate(); err != nil {
		return "", err
	}

	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID, "receipt", receipt.GetID())
	result := OpenMessageReceipt{}
	err = integration.Client.Post(
		log.ToContext(context),
		NewURI("/conversations/messages/%s/inbound/open/receipt", integration.ID),
		receipt,
		&result,
	)
	if err != nil {
		return "", err
	}
	if result.IsFailed() {
		return "", result.AsError()
	}
	return result.ID, err
}

// SendInboundEvent sends an event from the middleware to GENESYS Cloud
//
// See https://developer.genesys.cloud/commdigital/digital/openmessaging/inboundEventMessages
func (integration *OpenMessagingIntegration) SendInboundEvents(context context.Context, events OpenMessageEvents) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	events.Channel.MessageID = ""
	events.Channel.Platform = "Open"
	events.Channel.Type = "Private"
	events.Channel.Time = time.Now().UTC()
	events.Channel.To = &OpenMessageTo{ID: integration.ID.String()}
	if err := events.Channel.Validate(); err != nil {
		return "", err
	}
	result := OpenMessageEvents{}
	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID, "message", events.GetID())
	err = integration.Client.Post(
		log.ToContext(context),
		NewURI("/conversations/messages/%s/inbound/open/event", integration.ID),
		events,
		&result,
	)
	return result.ID, err
}

// SendOutboundMessage sends a message from GENESYS Cloud to the middleware
//
// The message can be only text as it is sent bia the AgentLess Message API.
//
// # This is mainly for debugging purposes
//
// See https://developer.genesys.cloud/api/digital/openmessaging/outboundMessages#send-an-agentless-outbound-text-message
func (integration *OpenMessagingIntegration) SendOutboundMessage(context context.Context, destination, text string) (*AgentlessMessageResult, error) {
	if integration.ID == uuid.Nil {
		return nil, errors.ArgumentMissing.With("ID")
	}
	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID)
	result := &AgentlessMessageResult{}
	err := integration.Client.Post(
		log.ToContext(context),
		"/conversations/messages/agentless",
		AgentlessMessage{
			From:          integration.ID.String(),
			To:            destination,
			MessengerType: "Open",
			Text:          text,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetMessageData gets the details ofa message
func (integration *OpenMessagingIntegration) GetMessageData(context context.Context, message OpenMessage) (*OpenMessageData, error) {
	if integration.ID == uuid.Nil {
		return nil, errors.ArgumentMissing.With("ID")
	}
	if len(message.GetID()) == 0 {
		return nil, errors.ArgumentMissing.With("messageID")
	}
	log := logger.Must(logger.FromContext(context, integration.logger)).Child("integration", "getmessagedata", "integration", integration.ID, "message", message.GetID())
	data := &OpenMessageData{}
	err := integration.Client.Get(
		log.ToContext(context),
		NewURI("/conversations/messages/%s/details", message.GetID()),
		data,
	)
	if err != nil {
		return nil, err
	}
	data.Conversation.client = integration.Client
	data.Conversation.logger = integration.logger.Child("conversation", "conversation", "id", data.Conversation.ID)
	return data, nil
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (integration OpenMessagingIntegration) String() string {
	if len(integration.Name) > 0 {
		return integration.Name
	}
	return integration.ID.String()
}

// MarshalJSON marshals this into JSON
func (integration OpenMessagingIntegration) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessagingIntegration
	data, err := json.Marshal(struct {
		surrogate
		WebhookURL *core.URL `json:"outboundNotificationWebhookUrl"`
		SelfURI    URI       `json:"selfUri"`
	}{
		surrogate:  surrogate(integration),
		WebhookURL: (*core.URL)(integration.WebhookURL),
		SelfURI:    integration.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (integration *OpenMessagingIntegration) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessagingIntegration
	var inner struct {
		surrogate
		WebhookURL *core.URL `json:"outboundNotificationWebhookUrl"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*integration = OpenMessagingIntegration(inner.surrogate)
	integration.WebhookURL = (*url.URL)(inner.WebhookURL)
	return
}

// TODO: There is also a PATCH method... we might want to provide some func

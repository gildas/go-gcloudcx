package gcloudcx

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
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
	SelfURI          URI                   `json:"selfUri,omitempty"`
	client           *Client               `json:"-"`
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

// Fetch fetches an OpenMessaging Integration
//
// implements Fetchable
func (integration *OpenMessagingIntegration) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, integration, parameters...)

	if id != uuid.Nil {
		if err := client.Get(ctx, NewURI("/conversations/messaging/integrations/open/%s", id), &integration); err != nil {
			return err
		}
		integration.logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(ctx, selfURI, &integration); err != nil {
			return err
		}
		integration.logger = log.Record("id", integration.ID)
	} else if len(name) > 0 {
		response := struct {
			Integrations []*OpenMessagingIntegration `json:"entities"`
			PageSize     int                         `json:"pageSize"`
			PageNumber   int                         `json:"pageNumber"`
			PageCount    int                         `json:"pageCount"`
			PageTotal    int                         `json:"total"`
			FirstURI     string                      `json:"firstUri"`
			SelfURI      string                      `json:"selfUri"`
			LastURI      string                      `json:"lastUri"`
		}{}
		if err := client.Get(ctx, "/conversations/messaging/integrations/open", &response); err != nil {
			return err
		}
		nameLowercase := strings.ToLower(name)
		for _, item := range response.Integrations {
			if strings.Compare(strings.ToLower(item.Name), nameLowercase) == 0 {
				*integration = *item
				break
			}
		}
		if integration == nil || integration.ID == uuid.Nil {
			return errors.NotFound.With("name", name)
		}
		integration.logger = log.Record("id", integration.ID)
	} else {
		return errors.ArgumentMissing.With("idOrName")
	}
	integration.client = client
	return nil
}

// FetchOpenMessagingIntegrations Fetches all OpenMessagingIntegration object
func (client *Client) FetchOpenMessagingIntegrations(ctx context.Context, parameters ...interface{}) ([]*OpenMessagingIntegration, error) {
	_, _, _, log := client.ParseParameters(ctx, nil, parameters...)
	entities := struct {
		Integrations []*OpenMessagingIntegration `json:"entities"`
		PageSize     int                         `json:"pageSize"`
		PageNumber   int                         `json:"pageNumber"`
		PageCount    int                         `json:"pageCount"`
		PageTotal    int                         `json:"total"`
		FirstURI     string                      `json:"firstUri"`
		SelfURI      string                      `json:"selfUri"`
		LastURI      string                      `json:"lastUri"`
	}{}
	if err := client.Get(ctx, "/conversations/messaging/integrations/open", &entities); err != nil {
		return nil, err
	}
	log.Record("response", entities).Infof("Got a response")
	for _, integration := range entities.Integrations {
		integration.client = client
		integration.logger = log.Child("openmessagingintegration", "openmessagingintegration", "id", integration.ID)
	}
	// TODO: fetch all pages!!!
	return entities.Integrations, nil
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
	integration.client = client
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
	return integration.client.Delete(
		integration.logger.ToContext(context),
		NewURI("/conversations/messaging/integrations/open/%s", integration.ID),
		nil,
	)
}

func (integration *OpenMessagingIntegration) Refresh(ctx context.Context) error {
	var value OpenMessagingIntegration
	if err := integration.client.Get(ctx, integration.GetURI(), &value); err != nil {
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
	response := &OpenMessagingIntegration{}
	err := integration.client.Patch(
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
	integration.logger.Record("response", response).Debugf("Updated integration %#v", response)
	return nil
}

// SendInboundMessage sends a text message from the middleware to GENESYS Cloud
//
// See https://developer.genesys.cloud/api/digital/openmessaging/inboundMessages#send-an-inbound-open-message
func (integration *OpenMessagingIntegration) SendInboundMessage(context context.Context, from *OpenMessageFrom, messageID, text string, attributes map[string]string, metadata map[string]string) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	if len(messageID) == 0 {
		return "", errors.ArgumentMissing.With("messageID")
	}
	result := OpenMessageText{}
	err = integration.client.Post(
		integration.logger.ToContext(context),
		"/conversations/messages/inbound/open",
		&OpenMessageText{
			Direction: "Inbound",
			Channel: NewOpenMessageChannel(
				messageID,
				&OpenMessageTo{ID: integration.ID.String()},
				from,
			).WithAttributes(attributes),
			Text:     text,
			Metadata: metadata,
		},
		&result,
	)
	return result.ID, err
}

// SendInboundMessageWithAttachment sends a text message with an attachemnt from the middleware to GENESYS Cloud
//
// See https://developer.genesys.cloud/api/digital/openmessaging/inboundMessages#inbound-message-with-attached-photo
// See https://developer.genesys.cloud/api/rest/v2/conversations/#post-api-v2-conversations-messages-inbound-open
func (integration *OpenMessagingIntegration) SendInboundMessageWithAttachment(context context.Context, from *OpenMessageFrom, messageID, text string, attachment *OpenMessageAttachment, attributes map[string]string, metadata map[string]string) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	if len(messageID) == 0 {
		return "", errors.ArgumentMissing.With("messageID")
	}
	if attachment.URL == nil {
		return "", errors.ArgumentMissing.With("url")
	}

	result := OpenMessageText{}
	err = integration.client.Post(
		integration.logger.ToContext(context),
		"/conversations/messages/inbound/open",
		&OpenMessageText{
			Direction: "Inbound",
			Channel: NewOpenMessageChannel(
				messageID,
				&OpenMessageTo{ID: integration.ID.String()},
				from,
			).WithAttributes(attributes),
			Text: text,
			Content: []*OpenMessageContent{
				{
					Type:       "Attachment",
					Attachment: attachment,
				},
			},
			Metadata: metadata,
		},
		&result,
	)
	return result.ID, err
}

// SendInboundReceipt sends a receipt from the middleware to GENESYS Cloud
//
// Valid status values are: Delivered, Failed, Published, Read, Removed, and Sent.
//
// See https://developer.genesys.cloud/api/digital/openmessaging/inboundMessages#send-an-inbound-open-message
func (integration *OpenMessagingIntegration) SendInboundReceipt(context context.Context, from *OpenMessageFrom, messageID, status string, reasons []StatusReason, attributes map[string]string, metadata map[string]string) (id string, err error) {
	if integration.ID == uuid.Nil {
		return "", errors.ArgumentMissing.With("ID")
	}
	if len(messageID) == 0 {
		return "", errors.ArgumentMissing.With("messageID")
	}
	result := OpenMessageText{}
	err = integration.client.Post(
		integration.logger.ToContext(context),
		"/conversations/messages/inbound/open",
		&OpenMessageReceipt{
			Direction: "Inbound",
			Channel: NewOpenMessageChannel(
				messageID,
				&OpenMessageTo{ID: integration.ID.String()},
				from,
			).WithAttributes(attributes),
			Status:   status,
			Reasons:  reasons,
			Metadata: metadata,
		},
		&result,
	)
	return result.ID, err
}

// SendOutboundMessage sends a message from GENESYS Cloud to the middleware
//
// The message can be only text as it is sent bia the AgentLess Message API.
//
// This is mainly for debugging purposes
//
// See https://developer.genesys.cloud/api/digital/openmessaging/outboundMessages#send-an-agentless-outbound-text-message
func (integration *OpenMessagingIntegration) SendOutboundMessage(context context.Context, destination, text string) (*AgentlessMessageResult, error) {
	if integration.ID == uuid.Nil {
		return nil, errors.ArgumentMissing.With("ID")
	}
	result := &AgentlessMessageResult{}
	err := integration.client.Post(
		integration.logger.ToContext(context),
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
	data := &OpenMessageData{}
	err := integration.client.Get(
		integration.logger.ToContext(context),
		NewURI("/conversations/messages/%s/details", message.GetID()),
		data,
	)
	if err != nil {
		return nil, err
	}
	data.Conversation.client = integration.client
	data.Conversation.logger = integration.logger.Child("conversation", "conversation", "id", data.Conversation.ID)
	return data, nil
}

// String gets a string version
//
//   implements the fmt.Stringer interface
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

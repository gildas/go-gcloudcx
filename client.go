package gcloudcx

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Client is the primary object to use Gcloud
type Client struct {
	Region         string         `json:"region"`
	DeploymentID   uuid.UUID      `json:"deploymentId"`
	Organization   *Organization  `json:"-"`
	API            *url.URL       `json:"apiUrl,omitempty"`
	LoginURL       *url.URL       `json:"loginUrl,omitempty"`
	Proxy          *url.URL       `json:"proxyUrl,omitempty"`
	Grant          Authorizable   `json:"-"`
	RequestTimeout time.Duration  `json:"requestTimout"`
	Logger         *logger.Logger `json:"-"`
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Context        context.Context
	Region         string
	OrganizationID uuid.UUID
	DeploymentID   uuid.UUID
	Proxy          *url.URL
	Grant          Authorizable
	RequestTimeout time.Duration
	Logger         *logger.Logger
}

// NewClient creates a new Gcloud Client
func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	if options.RequestTimeout < 2*time.Second {
		options.RequestTimeout = 10 * time.Second
	}
	if log, err := logger.FromContext(options.Context); err == nil && options.Logger == nil {
		options.Logger = log
	}
	client := Client{
		Proxy:          options.Proxy,
		DeploymentID:   options.DeploymentID,
		Organization:   &Organization{ID: options.OrganizationID},
		Grant:          options.Grant,
		RequestTimeout: options.RequestTimeout,
	}
	return client.SetLogger(options.Logger).SetRegion(options.Region)
}

// SetLogger sets the logger
func (client *Client) SetLogger(log *logger.Logger) *Client {
	client.Logger = logger.CreateIfNil(log, "gcloudcx").Child("gcloudcx", "gcloudcx")
	return client
}

// SetRegion sets the region and its main API
func (client *Client) SetRegion(region string) *Client {
	client.Region = region
	client.API, _ = url.Parse(fmt.Sprintf("https://api.%s", region))
	client.LoginURL, _ = url.Parse(fmt.Sprintf("https://login.%s", region))
	return client
}

// SetAuthorizationGrant sets the Authorization Grant
func (client *Client) SetAuthorizationGrant(grant Authorizable) *Client {
	client.Grant = grant
	return client
}

// GetLogger gets the logger from the given Context
//
// If the Context is nil or does not contain a logger, it returns the default logger
func (client Client) GetLogger(context context.Context) *logger.Logger {
	if context != nil {
		if log, err := logger.FromContext(context); err == nil {
			return log
		}
	}
	return client.Logger
}

// IsAuthorized tells if the client has an Authorization Token
// It migt be expired and the app should login again as needed
func (client *Client) IsAuthorized() bool {
	return client.Grant.AccessToken().IsValid()
}

// CheckScopes checks if the current client allows/denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (client *Client) CheckScopes(context context.Context, scopes ...string) (permitted []string, denied []string, correlationID string, err error) {
	return client.CheckScopesWithID(context, client.Grant, scopes...)
}

// CheckScopesWithID checks if the given grant allows/denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (client *Client) CheckScopesWithID(context context.Context, id core.Identifiable, scopes ...string) (permitted []string, denied []string, correlationID string, err error) {
	if id.GetID() == uuid.Nil {
		return nil, nil, "", errors.ArgumentMissing.With("id")
	}
	subject, correlationID, err := Fetch[AuthorizationSubject](context, client, id)
	if err != nil {
		return []string{}, scopes, correlationID, err
	}
	permitted, denied = subject.CheckScopes(scopes...)
	return permitted, denied, correlationID, nil
}

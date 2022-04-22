package gcloudcx

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gildas/go-core"
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

// Fetch fetches an object from the Genesys Cloud API
//
// The object must implement the Fetchable interface
//
// Objects can be fetched by their ID:
//
//    client.Fetch(context, &User{ID: uuid.UUID})
//
//    client.Fetch(context, &User{}, uuid.UUID)
//
// or by their name:
//
//    client.Fetch(context, &User{}, "user-name")
//
// or by their URI:
//
//    client.Fetch(context, &User{}, URI("/api/v2/users/user-id"))
//
//    client.Fetch(context, &User{URI: "/api/v2/users/user-id"})
func (client *Client) Fetch(ctx context.Context, object Fetchable, parameters ...interface{}) error {
	if _, err := logger.FromContext(ctx); err != nil {
		ctx = client.Logger.ToContext(ctx)
	}
	return object.Fetch(ctx, client, parameters...)
}

// ParseParameters parses the parameters to get an id, name or URI
//
// the id can be a string or a uuid.UUID, or coming from the object
//
// the uri can be a URI, or coming from the object
//
// a logger.Logger is also returned, either from the context or the client
func (client *Client) ParseParameters(ctx context.Context, object interface{}, parameters ...interface{}) (uuid.UUID, string, URI, *logger.Logger) {
	var (
		id   uuid.UUID = uuid.Nil
		name string
		uri  URI
	)

	for _, parameter := range parameters {
		switch parameter := parameter.(type) {
		case uuid.UUID:
			id = parameter
		case string:
			name = parameter
		case URI:
			uri = parameter
		}
	}
	if identifiable, ok := object.(Identifiable); id == uuid.Nil && ok {
		id = identifiable.GetID()
	}
	if addressable, ok := object.(Addressable); len(uri) == 0 && ok {
		uri = addressable.GetURI()
	}
	log, err := logger.FromContext(ctx)
	if err != nil {
		log = client.Logger
	}
	if typed, ok := object.(core.TypeCarrier); ok {
		log = log.Child(typed.GetType(), typed.GetType())
	}
	if id != uuid.Nil {
		log = log.Record("id", id.String())
	}
	return id, name, uri, log
}

// CheckScopes checks if the current client allows/denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (client *Client) CheckScopes(context context.Context, scopes ...string) (permitted []string, denied []string) {
	return client.CheckScopesWithID(context, client.Grant, scopes...)
}

// CheckScopesWithID checks if the given grant allows/denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (client *Client) CheckScopesWithID(context context.Context, id core.Identifiable, scopes ...string) (permitted []string, denied []string) {
	var subject AuthorizationSubject

	if err := client.Fetch(context, &subject, id); err != nil {
		return []string{}, scopes
	}
	return subject.CheckScopes(scopes...)
}

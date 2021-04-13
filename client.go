package purecloud

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region             string             `json:"region"`
	DeploymentID       uuid.UUID          `json:"deploymentId"`
	Organization       *Organization      `json:"-"`
	API                *url.URL           `json:"apiUrl,omitempty"`
	LoginURL           *url.URL           `json:"loginUrl,omitempty"`
	Proxy              *url.URL           `json:"proxyUrl,omitempty"`
	AuthorizationGrant AuthorizationGrant `json:"auth"`
	RequestTimeout     time.Duration      `json:"requestTimout"`
	Logger             *logger.Logger     `json:"-"`
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID uuid.UUID
	DeploymentID   uuid.UUID
	Proxy          *url.URL
	RequestTimeout time.Duration
	Logger         *logger.Logger
}

// NewClient creates a new PureCloud Client
func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	if options.RequestTimeout < 2 * time.Second {
		options.RequestTimeout = 10 * time.Second
	}
	client := Client{
		Proxy:          options.Proxy,
		DeploymentID:   options.DeploymentID,
		Organization:   &Organization{ID: options.OrganizationID},
		RequestTimeout: options.RequestTimeout,
	}
	return client.SetLogger(options.Logger).SetRegion(options.Region)
}

// SetLogger sets the logger
func (client *Client) SetLogger(log *logger.Logger) *Client {
	client.Logger = logger.CreateIfNil(log, "PureCloud").Child("purecloud", "purecloud")
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
func (client *Client) SetAuthorizationGrant(grant AuthorizationGrant) *Client {
	client.AuthorizationGrant = grant
	return client
}

// IsAuthorized tells if the client has an Authorization Token
// It migt be expired and the app should login again as needed
func (client *Client) IsAuthorized() bool {
	return client.AuthorizationGrant.AccessToken().IsValid()
}

// Fetch fetches an initializable object
func (client *Client) Fetch(object Initializable) error {
	return object.Initialize(client)
}

package purecloud

import (
	"fmt"
	"net/url"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region             string             `json:"region"`
	DeploymentID       string             `json:"deploymentId"`
	Organization       *Organization      `json:"-"`
	API                *url.URL           `json:"apiUrl,omitempty"`
	Proxy              *url.URL           `json:"proxyUrl,omitempty"`
	AuthorizationGrant AuthorizationGrant `json:"auth"`
	Logger             *logger.Logger     `json:"-"`
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	Proxy          *url.URL
	Logger         *logger.Logger
}

// New creates a new PureCloud Client
func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	client := Client{
		Proxy:        options.Proxy,
		DeploymentID: options.DeploymentID,
		Organization: &Organization{ID: options.OrganizationID},
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
	var err error

	client.Region = region
	client.API, err = url.Parse(fmt.Sprintf("https://api.%s", region))
	if err != nil {
		client.API, _ = url.Parse("https://api.mypurecloud.com")
	}
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
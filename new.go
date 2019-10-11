package purecloud

import (
	"fmt"
	"net/url"

	"github.com/gildas/go-logger"
)

// See: https://developer.mypurecloud.com/api/webchat/agentchat.html
// See: https://developer.mypurecloud.com/api/tutorials/
// See: https://developer.mypurecloud.com/api/rest/authorization/

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	Proxy          *url.URL
	Logger         *logger.Logger
}

// New creates a new PureCloud Client
func New(options ClientOptions) *Client {
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	client := Client{
		Proxy:        options.Proxy,
		Organization: &Organization{},
		DeploymentID: options.DeploymentID,
		Logger:       options.Logger,
	}
	return client.SetLogger(options.Logger).SetRegion(options.Region)
}

// SetLogger sets the logger
func (client *Client) SetLogger(log *logger.Logger) (*Client) {
	client.Logger = logger.CreateIfNil(log, "PureCloud").Topic("purecloud").Scope("purecloud")
	return client
}

// SetRegion sets the region and its main API
func (client *Client) SetRegion(region string) (*Client) {
	var err error

	client.Region = region
	client.API, err = url.Parse(fmt.Sprintf("https://api.%s/api/v2/", region))
	if err != nil {
		client.API, _ = url.Parse("https://api.mypurecloud.com/api/v2/")
	}
	return client
}

// SetAuthorizationGrant sets the Authorization Grant
func (client *Client) SetAuthorizationGrant(grant AuthorizationGrant) (*Client) {
	client.AuthorizationGrant = grant
	return client
}

// IsAuthorized tells if the client has an Authorization Token
// It migt be expired and the app should login again as needed
func (client *Client) IsAuthorized() bool {
	return client.AuthorizationGrant.AccessToken().IsValid()
}
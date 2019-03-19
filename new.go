package purecloud

import (
	"fmt"
	"net/url"

	"github.com/gildas/go-logger"
)

// New creates a new PureCloud Client
func New(options ClientOptions) *Client {
	if options.Logger == nil {
		options.Logger = logger.Create("Purecloud")
	}
	options.Logger = options.Logger.Record("topic", "purecloud").Record("scope", "purecloud").Child().(*logger.Logger)
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	return &Client{
		Region:        options.Region,
		API:           getAPI(options.Region),
		Organization:  &Organization{},
		Authorization: &Authorization{GrantType: "ClientCredentials"},
		DeploymentID:  options.DeploymentID,
		Logger:        options.Logger,
	}
}

// SetRegion sets the region and its main API
func (client *Client) SetRegion(region string) {
	client.Region = region
	client.API    = getAPI(region)
}

func getAPI(region string) *url.URL {
	if api, err := url.Parse(fmt.Sprintf("https://api.%s/api/v2/", region)); err == nil {
		return api
	}
	api, _ := url.Parse("https://api.mypurecloud.com/api/v2/")
	return api
}
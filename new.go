package purecloud

import (
	"fmt"
	"net/url"

	logger "bitbucket.org/gildas_cherruel/go-logger"
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
	apiURL, err := url.Parse(fmt.Sprintf("https://api.%s/api/v2/", options.Region))
	if err != nil {
		apiURL, _ = url.Parse("https://api.mypurecloud.com/api/v2/")
	}
	return &Client{
		Region:       options.Region,
		API:          apiURL,
		Organization: &Organization{},
		DeploymentID: options.DeploymentID,
		Logger:       options.Logger,
	}
}
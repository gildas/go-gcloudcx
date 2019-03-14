package purecloud

import (
	"fmt"
	"net/url"
	"strings"

	logger "bitbucket.org/gildas_cherruel/go-logger"
)

type responseLogin struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   uint32 `json:"expires_in,omitempty"`
	Error       string `json:"error,omitempty"`
}

// New creates a new PureCloud Client
func New(options ClientOptions) *Client {
	if options.Logger == nil {
		options.Logger = logger.Create("Purecloud")
	}
	options.Logger = options.Logger.Record("topic", "purecloud").Record("scope", "purecloud").Child().(*logger.Logger)
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	apiURL, err := url.Parse(fmt.Sprintf("https://api.%s/api/v2", options.Region))
	if err != nil {
		apiURL, _ = url.Parse("https://api.mypurecloud.com/api/v2")
	}
	return &Client{
		Region:         options.Region,
		API:            apiURL,
		OrganizationID: options.OrganizationID,
		DeploymentID:   options.DeploymentID,
		Logger:         options.Logger,
	}
}

// Login logs in a Client to PureCloud
func (client *Client) Login(authorization Authorization) (err error) {
	log := client.Logger.Record("scope", "login").Child().(*logger.Logger)

	switch strings.ToLower(authorization.GrantType) {
	case "clientcredentials":
		log.Debugf("Login type: %s", authorization.GrantType)

		// sanitize the options
		if len(authorization.ClientID) == 0 {
			return fmt.Errorf("Missing Argument ClientID")
		}
		if len(authorization.Secret) == 0 {
			return fmt.Errorf("Missing Argument Secret")
		}

		// TODO: Should we encrypt this?!?
		client.Authorization = authorization

		return client.authorize()
	default:
		return fmt.Errorf("Invalid GrantType: %s", authorization.GrantType)
	}
	return nil
}

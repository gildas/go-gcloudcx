package purecloud

import (
	"fmt"
	"strings"

	"github.com/gildas/go-logger"
)

type responseLogin struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   uint32 `json:"expires_in,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Login logs in a Client to PureCloud
func (client *Client) Login(authorization *Authorization) (err error) {
	log := client.Logger.Record("scope", "login").Child().(*logger.Logger)

	log.Debugf("Login type: %s, region: %s", authorization.GrantType, client.Region)
	switch strings.ToLower(authorization.GrantType) {
	case "clientcredentials":
		// sanitize the options
		if len(authorization.ClientID) == 0 { return fmt.Errorf("Missing Argument ClientID") }
		if len(authorization.Secret)   == 0 { return fmt.Errorf("Missing Argument Secret") }

		client.Authorization.GrantType = "ClientCredentials"
		client.Authorization.ClientID  = authorization.ClientID
		client.Authorization.Secret    = authorization.Secret

		if err = client.authorize(); err != nil { return err }
		if client.Organization, err = client.GetMyOrganization(); err != nil { return err }
	default:
		return fmt.Errorf("Invalid GrantType: %s", authorization.GrantType)
	}
	return nil
}

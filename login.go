package purecloud

import (
	"fmt"
	"strings"
)

// Login logs in a Client to PureCloud
func (client *Client) Login(authorization *Authorization) (err error) {
	log := client.Logger.Scope("login").Child()

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

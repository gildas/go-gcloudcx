package purecloud

import (
	"fmt"
	"time"

	"github.com/gildas/go-request"
)

// ClientCredentialsGrant implements PureCloud's Client Credentials Grants
//   See: https://developer.mypurecloud.com/api/rest/authorization/use-client-credentials.html
type ClientCredentialsGrant struct {
	ClientID string
	Secret   string
	Token    AccessToken
}

// Authorize this Grant with PureCloud
func (grant *ClientCredentialsGrant) Authorize(client *Client) (err error) {
	log := client.Logger.Child(nil, "authorize", "grant", "client_credentials")

	log.Infof("Authenticating with %s using Client Credentials grant", client.Region)

	// Validates the Grant
	if len(grant.ClientID) == 0 { return fmt.Errorf("Missing Argument ClientID") }
	if len(grant.Secret)   == 0 { return fmt.Errorf("Missing Argument Secret") }

	// Resets the token before authenticating
	grant.Token.Reset()
	response := struct {
		AccessToken string `json:"access_token,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   uint32 `json:"expires_in,omitempty"`
		Error       string `json:"error,omitempty"`
	}{}

	err = client.SendRequest(
		"https://login." + client.Region + "/oauth/token",
		&request.Options{
			Authorization: request.BasicAuthorization(grant.ClientID, grant.Secret),
			Payload: map[string]string{
				"grant_type": "client_credentials",
			},
		},
		&response,
	)
	if err != nil { return err }

	// Saves the token
	grant.Token.Type      = response.TokenType
	grant.Token.Token     = response.AccessToken
	grant.Token.ExpiresOn = time.Now().Add(time.Duration(int64(response.ExpiresIn)))

	client.Organization, _ = client.GetMyOrganization()

	return
}

// AccessToken gives the access Token carried by this Grant
func (grant *ClientCredentialsGrant) AccessToken() *AccessToken {
	return &grant.Token
}
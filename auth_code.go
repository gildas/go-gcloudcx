package purecloud

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gildas/go-request"
)

// AuthorizationCodeGrant implements PureCloud's Client Authorization Code Grants
//   See: https://developer.mypurecloud.com/api/rest/authorization/use-authorization-code.html
type AuthorizationCodeGrant struct {
	ClientID    string
	Secret      string
	Code        string
	RedirectURL *url.URL
	Token       AccessToken
}

// Authorize this Grant with PureCloud
func (grant *AuthorizationCodeGrant) Authorize(client *Client) (err error) {
	log := client.Logger.Child(nil, "authorize", "grant", "authorization_code")

	log.Infof("Authenticating with %s using Authorization Code grant", client.Region)

	// Validates the Grant
	if len(grant.ClientID) == 0 {
		return fmt.Errorf("Missing Argument ClientID")
	}
	if len(grant.Secret) == 0 {
		return fmt.Errorf("Missing Argument Secret")
	}
	if len(grant.Code) == 0 {
		return fmt.Errorf("Missing Argument Code")
	}

	// Resets the token before authenticating
	grant.Token.Reset()
	response := struct {
		AccessToken string `json:"access_token,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   uint32 `json:"expires_in,omitempty"`
		Error       string `json:"error,omitempty"`
	}{}

	err = client.SendRequest(
		"https://login."+client.Region+"/oauth/token",
		&request.Options{
			Authorization: request.BasicAuthorization(grant.ClientID, grant.Secret),
			Payload: map[string]string{
				"grant_type":   "authorization_code",
				"code":         grant.Code,
				"redirect_uri": grant.RedirectURL.String(),
			},
		},
		&response,
	)
	if err != nil {
		return err
	}

	// Saves the token
	grant.Token.Type = response.TokenType
	grant.Token.Token = response.AccessToken
	grant.Token.ExpiresOn = time.Now().Add(time.Duration(int64(response.ExpiresIn)))

	return
}

// AccessToken gives the access Token carried by this Grant
func (grant *AuthorizationCodeGrant) AccessToken() *AccessToken {
	return &grant.Token
}

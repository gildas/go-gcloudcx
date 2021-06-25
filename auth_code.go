package purecloud

import (
	"net/url"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
	"github.com/google/uuid"
)

// AuthorizationCodeGrant implements PureCloud's Client Authorization Code Grants
//   See: https://developer.mypurecloud.com/api/rest/authorization/use-authorization-code.html
type AuthorizationCodeGrant struct {
	ClientID     uuid.UUID
	Secret       string
	Code         string
	RedirectURL  *url.URL
	Token        AccessToken
	CustomData   interface{}
	TokenUpdated chan UpdatedAccessToken
}

// GetID gets the client Identifier
//
// Implements core.Identifiable
func (grant *AuthorizationCodeGrant) GetID() uuid.UUID {
	return grant.ClientID
}

// Authorize this Grant with PureCloud
func (grant *AuthorizationCodeGrant) Authorize(client *Client) (err error) {
	log := client.Logger.Child(nil, "authorize", "grant", "authorization_code")

	log.Infof("Authenticating with %s using Authorization Code grant", client.Region)

	// Validates the Grant
	if grant.ClientID == uuid.Nil {
		return errors.ArgumentMissing.With("ClientID").WithStack()
	}
	if len(grant.Secret) == 0 {
		return errors.ArgumentMissing.With("Secret").WithStack()
	}
	if len(grant.Code) == 0 {
		return errors.ArgumentMissing.With("Code").WithStack()
	}

	// Resets the token before authenticating
	grant.Token.Reset()
	response := struct {
		AccessToken string `json:"access_token,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   int64  `json:"expires_in,omitempty"`
		Error       string `json:"error,omitempty"`
	}{}

	err = client.SendRequest(
		NewURI("%s/oauth/token", client.LoginURL),
		&request.Options{
			Authorization: request.BasicAuthorization(grant.ClientID.String(), grant.Secret),
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
	grant.Token.ExpiresOn = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

	log.Debugf("New %s token expires on %s", grant.Token.Type, grant.Token.ExpiresOn)
	if grant.TokenUpdated != nil {
		log.Debugf("Sending new token to TokenUpdated chan")
		grant.TokenUpdated <- UpdatedAccessToken{
			AccessToken: grant.Token,
			CustomData:  grant.CustomData,
		}
	}
	return
}

// AccessToken gives the access Token carried by this Grant
func (grant *AuthorizationCodeGrant) AccessToken() *AccessToken {
	return &grant.Token
}

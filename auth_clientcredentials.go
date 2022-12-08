package gcloudcx

import (
	"context"
	"net/http"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
	"github.com/google/uuid"
)

// ClientCredentialsGrant implements GCloud's Client Credentials Grants
//
// When the Token is updated, the new token is sent to the TokenUpdated chan along with the CustomData
//
//	See: https://developer.mypurecloud.com/api/rest/authorization/use-client-credentials.html
type ClientCredentialsGrant struct {
	ClientID     uuid.UUID
	Secret       string
	Token        AccessToken
	CustomData   interface{}
	TokenUpdated chan UpdatedAccessToken
}

// GetID gets the client Identifier
//
// Implements core.Identifiable
func (grant *ClientCredentialsGrant) GetID() uuid.UUID {
	return grant.ClientID
}

// Authorize this Grant with GCloud CX
//
// Implements Authorizable
func (grant *ClientCredentialsGrant) Authorize(context context.Context, client *Client) (err error) {
	log := client.GetLogger(context).Child(nil, "authorize", "grant", "client_credentials")

	log.Infof("Authenticating with %s using Client Credentials grant", client.Region)

	// Validates the Grant
	if grant.ClientID == uuid.Nil {
		return errors.ArgumentMissing.With("ClientID")
	}
	if len(grant.Secret) == 0 {
		return errors.ArgumentMissing.With("Secret")
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
		context,
		NewURI("%s/oauth/token", client.LoginURL),
		&request.Options{
			Method:        http.MethodPost,
			Authorization: request.BasicAuthorization(grant.ClientID.String(), grant.Secret),
			Payload: map[string]string{
				"grant_type": "client_credentials",
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
	client.Organization, _ = client.GetMyOrganization(context)

	return
}

// AccessToken gives the access Token carried by this Grant
//
// Implements Authorizable
func (grant *ClientCredentialsGrant) AccessToken() *AccessToken {
	return &grant.Token
}

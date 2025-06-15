package gcloudcx

import (
	"context"

	"github.com/gildas/go-core"
	"github.com/gildas/go-request"
	"github.com/google/uuid"
)

// TokenGrant can be used when you already have a token
// and want to use it to access GCloud CX
//
// If the token is expired, this grant will stop working
type TokenGrant struct {
	Token AccessToken
}

// GetID gets the client Identifier
//
// Implements core.Identifiable
func (grant *TokenGrant) GetID() uuid.UUID {
	return grant.Token.ID
}

// Authorize this Grant with GCloud CX
//
// Implements Authorizable
func (grant *TokenGrant) Authorize(context context.Context, client *Client) (correlationID string, err error) {
	log := client.GetLogger(context).Child("client", "authorize", "grant", "token")

	log.Debugf("Authenticating with %s using Token grant", client.Region)

	var response struct {
		Organization     Organization `json:"organization"`
		HomeOrganization Organization `json:"homeOrganization"`
		AuthorizedScopes []string     `json:"authorizedScope"` // there is no plural from the API
		OAuthClient      struct {
			ID           core.UUID `json:"id"`
			Name         string    `json:"name"`
			Organization struct {
				ID string `json:"id"` // This is not always a UUID
			} `json:"organization"`
		} `json:"OAuthClient"`
	}

	correlationID, err = client.SendRequest(
		context,
		NewURI("/tokens/me"),
		&request.Options{
			Authorization: request.BearerAuthorization(grant.Token.Token),
		},
		&response,
	)
	if err != nil {
		log.Errorf("Failed to authenticate with %s using Token grant", client.Region, err)
		return correlationID, err
	}

	log.Infof("Authenticated with %s using Token grant", client.Region)
	client.Organization = &response.Organization
	grant.Token.Type = "Bearer"
	grant.Token.ID = uuid.UUID(response.OAuthClient.ID)
	grant.Token.AuthorizedScopes = response.AuthorizedScopes
	return
}

// AccessToken gives the access Token carried by this Grant
//
// Implements Authorizable
func (grant *TokenGrant) AccessToken() *AccessToken {
	return &grant.Token
}

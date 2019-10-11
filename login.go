package purecloud

import (
	"net/url"
	"fmt"
	"time"
)

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	ClientID     string                 `json:"clientId"`
	Secret       string                 `json:"clientSecret"`
	RedirectURI  *url.URL               `json:"redirectUri"`
	TokenType    string                 `json:"tokenType"`
	Token        string                 `json:"token"`
	TokenExpires time.Time              `json:"tokenExpires"`
}

// Login logs in a Client to PureCloud
//   Uses the credentials stored in the Client
func (client *Client) Login() error {
	return client.LoginWithAuthorizationGrant(client.AuthorizationGrant)
}

// LoginWithAuthorizationGrant logs in a Client to PureCloud with given authorization Grant
func (client *Client) LoginWithAuthorizationGrant(authorizationGrant AuthorizationGrant) (err error) {
	if authorizationGrant == nil {
		return fmt.Errorf("Authorization Grant cannot be nil")
	}
	if err = authorizationGrant.Authorize(client); err != nil {
		return err
	}
	client.Organization, err = client.GetMyOrganization()
	return
		/*
	case AuthorizationCodeGrant:
		// sanitize the options
		if len(authorization.ClientID) == 0   { return fmt.Errorf("Missing Argument ClientID") }
		if len(authorization.Secret)   == 0   { return fmt.Errorf("Missing Argument Secret") }
		if authorization.RedirectURI   == nil { return fmt.Errorf("Missing Argument RedirectURI") }
		*/
}

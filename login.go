package purecloud

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	GrantType    AuthorizationGrantType `json:"grantType"`
	ClientID     string                 `json:"clientId"`
	Secret       string                 `json:"clientSecret"`
	TokenType    string                 `json:"tokenType"`
	Token        string                 `json:"token"`
	TokenExpires time.Time              `json:"tokenExpires"`
}

// AuthorizationGrantType defines the GrantType that can be used during the login process
type AuthorizationGrantType int

const (
	// ClientCredentialsGrant is used to login PureCloud with Client Credentials
	ClientCredentialsGrant = iota
)

func (grant AuthorizationGrantType) String() string {
	return [...]string{"client_credentials"}[grant]
}

// Login logs in a Client to PureCloud
//   Uses the credentials stored in the Client
func (client *Client) Login() error {
	return client.LoginWithCredentials(client.Authorization)
}

// LoginWithCredentials logs in a Client to PureCloud with given credentials
func (client *Client) LoginWithCredentials(authorization *Authorization) (err error) {
	log := client.Logger.Scope("login").Child()

	if authorization == nil {
		authorization = client.Authorization
	}
	log.Debugf("Login type: %s, region: %s", authorization.GrantType.String(), client.Region)
	switch authorization.GrantType {
	case ClientCredentialsGrant:
		// sanitize the options
		if len(authorization.ClientID) == 0 { return fmt.Errorf("Missing Argument ClientID") }
		if len(authorization.Secret)   == 0 { return fmt.Errorf("Missing Argument Secret") }

		response := struct {
			AccessToken string `json:"access_token,omitempty"`
			TokenType   string `json:"token_type,omitempty"`
			ExpiresIn   uint32 `json:"expires_in,omitempty"`
			Error       string `json:"error,omitempty"`
		}{}
		err = client.request(
			http.MethodPost,
			"https://login." + client.Region + "/oauth/token",
			[]byte("grant_type=" + authorization.GrantType.String()),
			&response,
			RequestOptions{
				ContentType:   "application/x-www-form-urlencoded",
				Authorization: "Basic " + base64.StdEncoding.EncodeToString([]byte(authorization.ClientID + ":" + authorization.Secret)),
			},
		)
		if err != nil { return }

		// Saves auth stuff and response
		client.Authorization.GrantType    = authorization.GrantType
		client.Authorization.ClientID     = authorization.ClientID
		client.Authorization.Secret       = authorization.Secret
		client.Authorization.TokenType    = response.TokenType
		client.Authorization.Token        = response.AccessToken
		client.Authorization.TokenExpires = time.Now().Add(time.Duration(int64(response.ExpiresIn)))

		client.Organization, err = client.GetMyOrganization()
	default:
		return fmt.Errorf("Invalid GrantType: %s", authorization.GrantType)
	}
	return
}

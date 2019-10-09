package purecloud

import (
	"context"
	"github.com/gildas/go-core"
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

		err := client.sendRequest(
			http.MethodPost,
			"https://login." + client.Region + "/oauth/token",
			&response,
		)

		url, err := client.parseURL("https://login." + client.Region + "/oauth/token")
		if err != nil {
			return APIError{ Code: "url.parse", Message: err.Error() }
		}
		res, err := core.SendRequest(context.Background(), &core.RequestOptions{
			Method:     http.MethodPost,
			URL:        url,
			Proxy:      client.Proxy,
			UserAgent:  APP + " " + VERSION,
			Headers:    map[string]string {
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(authorization.ClientID + ":" + authorization.Secret)),
			},
			Content: core.ContentReader{
				Type: "application/x-www-form-urlencoded",
			},
			Parameters: map[string]string{
				"grant_type": authorization.GrantType.String(),
			},
			Logger: log,
		}, &response)

		if err != nil {
			log.Record("err", err).Errorf("Core SendRequest error", err)
			if res != nil {
				log.Infof("Reading error from res")
				apiError := APIError{}
				err = res.UnmarshalContentJSON(&apiError)
				if err != nil { return err }
				return apiError
			}
			return err
		}

		/*
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
		if err != nil { return err }
		*/

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

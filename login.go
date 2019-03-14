package purecloud

import (
	"net/url"
	"bitbucket.org/gildas_cherruel/go-logger"
	"time"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"
	"strings"
	"net/http"
	"fmt"
)

type responseLogin struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   uint32 `json:"expires_in,omitempty"`
	Error       string `json:"error,omitempty"`
}

// New creates a new PureCloud Client
func New(options ClientOptions) *Client {
	if options.Logger == nil {
		options.Logger = logger.Create("Purecloud")
	}
	options.Logger = options.Logger.Record("topic", "purecloud").Record("scope", "purecloud").Child().(*logger.Logger)
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	apiURL, err := url.Parse(fmt.Sprintf("https://api.%s/api/v2", options.Region))
	if err != nil {
		apiURL, _ = url.Parse("https://api.mypurecloud.com/api/v2")
	}
	return &Client{
		Region:         options.Region,
		API:            *apiURL,
		OrganizationID: options.OrganizationID,
		DeploymentID:   options.DeploymentID,
		Logger:         options.Logger,
	}
}

// Login logs in a Client to PureCloud
func (client *Client) Login(options LoginOptions) (err error) {
	log := client.Logger.Record("scope", "login").Child().(*logger.Logger)

	switch (strings.ToLower(options.GrantType)) {
	case "clientcredentials":
		log.Debugf("Login type: %s", options.GrantType)

		// sanitize the options
		if len(options.ClientID) == 0 {
			return fmt.Errorf("Missing Argument ClientID")
		}
		if len(options.Secret) == 0 {
			return fmt.Errorf("Missing Argument Secret")
		}

		var request  *http.Request
		var response *http.Response
		loginURL, _ := url.Parse("https://login." + client.Region + "/oauth/token")

		if request, err = http.NewRequest("POST", loginURL.String(), strings.NewReader("grant_type=client_credentials")); err != nil {
			return
		}
		request.Header.Set("User-Agent", APP + " " + VERSION)
		request.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(options.ClientID + ":" + options.Secret)))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Accept", "application/json")

		log.Debugf("Sending POST request to %s", request.URL.String())
		if response, err = http.DefaultClient.Do(request); err != nil {
			return
		}
		defer response.Body.Close()

		log.Debugf("Response: %d - %s", response.StatusCode, response.Status)
		if response.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(response.Body)
			return fmt.Errorf("Error during login: %d - %s, %s", response.StatusCode, response.Status, string(body))
		}

		body := &responseLogin{}
		if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
			return
		}

		client.Token  = Token {
			Type:  body.TokenType,
			Token: body.AccessToken,
			Expires: time.Now().Add(time.Duration(int64(body.ExpiresIn))),
		}
	default:
		return fmt.Errorf("Invalid GrantType: %s", options.GrantType)
	}
	return nil
}
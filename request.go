package purecloud

import (
	"encoding/base64"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	logger "bitbucket.org/gildas_cherruel/go-logger"
)

type requestOptions struct {
	ContentType   string
	Authorization string
}

type responseAuth struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   uint32 `json:"expires_in,omitempty"`
	Error       string `json:"error,omitempty"`
}
// post sends a POST HTTP Request to PureCloud and gets the result
func (client *Client) post(path string, payload []byte, data interface{}) error {
	return client.request(http.MethodPost, path, payload, data, requestOptions{})
}

// get sends a GET HTTP Request to PureCloud and gets the result
func (client *Client) get(path string, payload []byte, data interface{}) error {
	return client.request(http.MethodGet, path, payload, data, requestOptions{})
}

// authorize sends a client credentials authentication request to PureCloud
func (client *Client) authorize() error {
	auth := &responseAuth{}
	if err := client.request(
		http.MethodPost,
		"https://login." + client.Region + "/oauth/token",
		[]byte("grant_type=client_credentials"),
		auth,
		requestOptions{
			ContentType:   "application/x-www-form-urlencoded",
			Authorization: "Basic " + base64.StdEncoding.EncodeToString([]byte(client.Authorization.ClientID + ":" + client.Authorization.Secret)),
		},
	); err != nil {
		return err
	}
	client.Token = Token{
		Type:    auth.TokenType,
		Token:   auth.AccessToken,
		Expires: time.Now().Add(time.Duration(int64(auth.ExpiresIn))),
	}
	return nil
}

// request sends an HTTP Request to PureCloud and gets the result
func (client *Client) request(method, path string, payload []byte, data interface{}, options requestOptions) error {
	log := client.Logger.Record("scope", "request."+method).Child().(*logger.Logger)

	url, err := client.parseURL(path)
	if err != nil {
		return err
	}

	// Creating a new HTTP request with the payload
	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Setting common Headers
	req.Header.Set("User-Agent", APP+" "+VERSION)
	if len(options.Authorization) > 0 {
		req.Header.Set("Authorization", options.Authorization)
	} else {
		if len(client.Token.Token) == 0 {
			if err := client.authorize(); err != nil {
				return err
			}
		}
		req.Header.Set("Authorization", client.Token.Type + " " + client.Token.Token)
	}
	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Sending the Request
	httpClient := http.DefaultClient
	if client.Proxy != nil {
		httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(client.Proxy)}
	}
	start := time.Now()
	log.Debugf("Sending %s request to %s", method, req.URL.String())
	res, err := httpClient.Do(req)
	duration := time.Since(start)
	log = log.Record("duration", duration.Seconds()).Child().(*logger.Logger)

	if err != nil {
		log.Errorf("Failed sending %s reqquest to %s", method, req.URL.String())
		return err
	}
	defer res.Body.Close()

	// TODO: Process redirections (3xx)
	if res.StatusCode == 401 && len(options.Authorization) == 0 { // Typically we need to acquire our token again
		if err := client.authorize(); err != nil {
			return err
		}
		return client.request(method, path, payload, data, options)
	}

	if res.StatusCode >= 400 {
		log.Errorf("Error while sending request \nstatus: %d %s, \nHeaders: %#v, Content-Length: %d", res.Status, res.StatusCode, res.Header, res.ContentLength)
		if res.ContentLength > 0 {
			body, _ := ioutil.ReadAll(res.Body)
			return fmt.Errorf("HTTP Request failed %d %s, %s", res.StatusCode, res.Status, body)
		}
		return fmt.Errorf("HTTP Request failed %d %s", res.StatusCode, res.Status)
	}

	log.Debugf("Successfully sent request in %s \nstatus: %d %s, \nHeaders: %#v, \nContent-Length: %d", duration, res.StatusCode, res.Status, res.Header, res.ContentLength)

	if data != nil {
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return err
		}
	}
	return nil
}

// parseURL parses a given Path into a useable URL with PureCloud
func (client *Client) parseURL(path string) (*url.URL, error) {
	if strings.HasPrefix(path, "http") {
		return url.Parse(path)
	}
	return client.API.Parse(path)
}

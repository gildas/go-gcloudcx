package purecloud

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestOptions contains options for requests
type RequestOptions struct {
	ContentType   string
	Authorization string
}

// Post sends a POST HTTP Request to PureCloud and gets the result
func (client *Client) Post(path string, payload []byte, data interface{}, options ...RequestOptions) error {
	return client.request(http.MethodPost, path, payload, data, options...)
}

// Get sends a GET HTTP Request to PureCloud and gets the result
func (client *Client) Get(path string, payload []byte, data interface{}, options ...RequestOptions) error {
	return client.request(http.MethodGet, path, payload, data, options...)
}

// Delete sends a DELETE HTTP Request to PureCloud and gets the result
func (client *Client) Delete(path string, payload []byte, data interface{}, options ...RequestOptions) error {
	return client.request(http.MethodDelete, path, payload, data, options...)
}

// request sends an HTTP Request to PureCloud and gets the result
func (client *Client) request(method, path string, payload []byte, data interface{}, options ...RequestOptions) error {
	log := client.Logger.Scope("request."+method).Child()

	url, err := client.parseURL(path)
	if err != nil {
		return APIError{ Code: "url.parse", Message: err.Error() }
	}

	// Creating a new HTTP request with the payload
	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(payload))
	if err != nil {
		return APIError{ Code: "http.request.create", Message: err.Error() }
	}

	// Grabbing latest options
	authorization := ""
	contentType   := "application/json"
	for _, option := range options {
		if len(option.Authorization) > 0 {
			authorization = option.Authorization
		}
		if len(option.ContentType) > 0 {
			contentType = option.ContentType
		}
	}

	// Setting common Headers
	req.Header.Set("User-Agent", APP+" "+VERSION)
	if len(authorization) > 0 {
		req.Header.Set("Authorization", authorization)
	} else {
		if len(client.Authorization.Token) == 0 {
			if err := client.Login(); err != nil { return err }
		}
		req.Header.Set("Authorization", client.Authorization.TokenType + " " + client.Authorization.Token)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	// Sending the Request
	httpClient := http.DefaultClient
	if client.Proxy != nil {
		log.Debugf("Setting HTTP Proxt to: %s", client.Proxy)
		httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(client.Proxy)}
	}
	start := time.Now()
	log.Debugf("Sending %s request to %s", method, req.URL.String())
	log.Tracef("Request Headers: %#v", req.Header)
	if len(payload) > 512 {
		log.Tracef("Data: %s", string(payload[:511]))
	} else {
		log.Tracef("Data: %s", string(payload))
	}
	res, err := httpClient.Do(req)
	duration := time.Since(start)
	log = log.Record("duration", duration.Seconds()).Child()

	if err != nil {
		return APIError{
			Code:              "http.request.send",
			Message:           "Failed sending %s request to %s: %s",
			MessageParams:     map[string]string{ "method": method, "url": req.URL.String(), "error": err.Error() },
			MessageWithParams: fmt.Sprintf("Failed sending %s request to %s: %s", method, req.URL.String(), err),
		}
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) // read the body no matter what

	// TODO: Process redirections (3xx)
	// TODO: Handle retry-after (https://developer.mypurecloud.com/forum/t/new-rate-limit-header-retry-after/4777)
	if res.StatusCode == 401 && len(authorization) == 0 { // Typically we need to acquire our token again
		if err := client.Login(); err != nil { return err }
		return client.request(method, path, payload, data, options...)
	}

	if res.StatusCode >= 400 {
		log.Errorf("Error while sending request \nstatus: %s, \nHeaders: %#v, Content-Length: %d, \nBody: %s", res.Status, res.Header, res.ContentLength, string(body))
		apiError := APIError{}
		if err := json.Unmarshal(body, &apiError); err != nil {
			apiError = APIError{ Code: "json.parse", Message: err.Error() }
		}
		if apiError.Status == 0         { apiError.Status = res.StatusCode }
		if len(apiError.ContextID) == 0 { apiError.ContextID = res.Header.Get("ININ-Correlation-Id") }
		return apiError
	}

	log.Debugf("Successfully sent request in %s \nstatus: %s, \nHeaders: %#v, \nContent-Length: %d", duration, res.Status, res.Header, res.ContentLength)

	if data != nil {
		if err := json.Unmarshal(body, &data); err != nil {
			return APIError{
				Status:    res.StatusCode,
				Code:      "json.parse",
				Message:   err.Error(),
				ContextID: res.Header.Get("ININ-Correlation-Id"),
			}
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

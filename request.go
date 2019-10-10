package purecloud

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gildas/go-core"
)

// Post sends a POST HTTP Request to PureCloud and gets the results
func (client *Client) Post(path string, payload, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{Method: http.MethodPost, Payload: payload}, results)
}

// Get sends a GET HTTP Request to PureCloud and gets the results
func (client *Client) Get(path string, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{}, results)
}

// Delete sends a DELETE HTTP Request to PureCloud and gets the results
func (client *Client) Delete(path string, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{Method: http.MethodDelete}, results)
}

// SendRequest sends a REST request to PureCloud via core.SendRequest
func (client *Client) SendRequest(path string, options *core.RequestOptions, results interface{}) (err error) {
	log := client.Logger.Scope("request." + options.Method)

	if options == nil { options = &core.RequestOptions{} }
	if strings.HasPrefix(path, "http") {
		options.URL, err = url.Parse(path)
	} else {
		options.URL, err = client.API.Parse(strings.TrimPrefix(path, "/"))
	}
	if err != nil {
		return APIError{ Code: "url.parse", Message: err.Error() }
	}
	if len(options.Authorization) == 0 && len(client.Authorization.Token) > 0 {
		options.Authorization = client.Authorization.TokenType + " " + client.Authorization.Token
	}

	options.Proxy     = client.Proxy
	options.UserAgent = APP + " " + VERSION
	options.Logger    = log

	res, err := core.SendRequest(options, results)
	if err != nil {
		log.Record("err", err).Errorf("Core SendRequest error", err)
		if res != nil {
			log.Infof("Reading error from res")
			apiError := APIError{}
			err = res.UnmarshalContentJSON(&apiError)
			if err != nil { return err }
			return apiError
		}
		return err // Make a nice APIError
	}
	return nil
}
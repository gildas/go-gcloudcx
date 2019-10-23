package purecloud

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gildas/go-core"
	"github.com/pkg/errors"
)

// Post sends a POST HTTP Request to PureCloud and gets the results
func (client *Client) Post(path string, payload, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{Method: http.MethodPost, Payload: payload}, results)
}

// Post sends a PATCH HTTP Request to PureCloud and gets the results
func (client *Client) Patch(path string, payload, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{Method: http.MethodPatch, Payload: payload}, results)
}

// Post sends an UPDATE HTTP Request to PureCloud and gets the results
func (client *Client) Put(path string, payload, results interface{}) error {
	return client.SendRequest(path, &core.RequestOptions{Method: http.MethodPut, Payload: payload}, results)
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
	if options == nil { options = &core.RequestOptions{} }
	if strings.HasPrefix(path, "http") {
		options.URL, err = url.Parse(path)
	} else {
		options.URL, err = client.API.Parse(strings.TrimPrefix(path, "/"))
	}
	if err != nil {
		return APIError{ Code: "url.parse", Message: err.Error() }
	}
	if len(options.Authorization) == 0 {
		if client.IsAuthorized() {
			options.Authorization = client.AuthorizationGrant.AccessToken().String()
		} else {
			if err = client.Login(); err != nil {
				return errors.WithStack(err)
			}
			if !client.IsAuthorized() {
				return errors.Errorf("Not Authorized Yet")
			}
		}
	}

	options.Proxy     = client.Proxy
	options.UserAgent = APP + " " + VERSION
	options.Logger    = client.Logger
	options.ResponseBodyLogSize = 4096

	res, err := core.SendRequest(options, results)
	if err != nil {
		if res != nil {
			apiError := APIError{}
			if jsonerr := res.UnmarshalContentJSON(&apiError); jsonerr != nil {
				return errors.Wrap(err, "Failed to extract an error from the response")
			}
			if err, ok := errors.Cause(err).(core.RequestError); ok {
				apiError.Status = err.StatusCode
			}
			return errors.WithStack(apiError)
		}
		return err // Make a nice APIError
	}
	return nil
}
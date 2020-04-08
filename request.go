package purecloud

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
)

// Post sends a POST HTTP Request to PureCloud and gets the results
func (client *Client) Post(path string, payload, results interface{}) error {
	return client.SendRequest(path, &request.Options{Method: http.MethodPost, Payload: payload}, results)
}

// Patch sends a PATCH HTTP Request to PureCloud and gets the results
func (client *Client) Patch(path string, payload, results interface{}) error {
	return client.SendRequest(path, &request.Options{Method: http.MethodPatch, Payload: payload}, results)
}

// Put sends an UPDATE HTTP Request to PureCloud and gets the results
func (client *Client) Put(path string, payload, results interface{}) error {
	return client.SendRequest(path, &request.Options{Method: http.MethodPut, Payload: payload}, results)
}

// Get sends a GET HTTP Request to PureCloud and gets the results
func (client *Client) Get(path string, results interface{}) error {
	return client.SendRequest(path, &request.Options{}, results)
}

// Delete sends a DELETE HTTP Request to PureCloud and gets the results
func (client *Client) Delete(path string, results interface{}) error {
	return client.SendRequest(path, &request.Options{Method: http.MethodDelete}, results)
}

// SendRequest sends a REST request to PureCloud
func (client *Client) SendRequest(path string, options *request.Options, results interface{}) (err error) {
	if options == nil {
		options = &request.Options{}
	}
	if strings.HasPrefix(path, "http") {
		options.URL, err = url.Parse(path)
	} else if client.API == nil {
		return errors.ArgumentMissing.With("Client API").WithStack()
	} else if !strings.HasPrefix(path, "/api") {
		options.URL, err = client.API.Parse("/api/v2/" + strings.TrimPrefix(path, "/"))
	} else {
		options.URL, err = client.API.Parse(path)
	}
	if err != nil {
		return errors.WithStack(APIError{Code: "url.parse", Message: err.Error()})
	}
	if len(options.Authorization) == 0 {
		if client.IsAuthorized() {
			options.Authorization = client.AuthorizationGrant.AccessToken().String()
		} else {
			if err = client.Login(); err != nil {
				return errors.WithStack(err)
			}
			if !client.IsAuthorized() {
				return errors.HTTPUnauthorized.WithMessage("Not Authorized Yet")
			}
			options.Authorization = client.AuthorizationGrant.AccessToken().String()
		}
	}

	options.Proxy = client.Proxy
	options.UserAgent = APP + " " + VERSION
	options.Logger = client.Logger
	options.ResponseBodyLogSize = 4096

	res, err := request.Send(options, results)
	if err != nil {
		if errors.Is(err, errors.HTTPUnauthorized) && len(client.AuthorizationGrant.AccessToken().String()) > 0 {
			// This means our token most probably expired, we should try again without it
			client.Logger.Infof("Authorization Token is expired, we need to authenticate again")
			options.Authorization = ""
			client.AuthorizationGrant.AccessToken().Reset()
			return client.SendRequest(path, options, results)
		}
		var details *errors.Error
		if errors.As(err, &details) {
			apiError := APIError{}
			if jsonerr := res.UnmarshalContentJSON(&apiError); jsonerr != nil {
				return errors.Wrap(err, "Failed to extract an error from the response")
			}
			apiError.Status = details.Code
			apiError.Code = details.ID
			if strings.HasPrefix(apiError.Message, "authentication failed") {
				apiError.Status = errors.HTTPUnauthorized.Code
				apiError.Code = errors.HTTPUnauthorized.ID
			}
			return errors.WithStack(apiError)
		}
		return err
	}
	return nil
}

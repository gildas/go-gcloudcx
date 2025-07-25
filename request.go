package gcloudcx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
)

// Post sends a POST HTTP Request to GCloud and gets the results
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) Post(context context.Context, path URI, payload, results interface{}) (correlationID string, err error) {
	return client.SendRequest(context, path, &request.Options{Method: http.MethodPost, Payload: payload}, results)
}

// Patch sends a PATCH HTTP Request to GCloud and gets the results
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) Patch(context context.Context, path URI, payload, results interface{}) (correlationID string, err error) {
	return client.SendRequest(context, path, &request.Options{Method: http.MethodPatch, Payload: payload}, results)
}

// Put sends an UPDATE HTTP Request to GCloud and gets the results
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) Put(context context.Context, path URI, payload, results interface{}) (correlationID string, err error) {
	return client.SendRequest(context, path, &request.Options{Method: http.MethodPut, Payload: payload}, results)
}

// Get sends a GET HTTP Request to GCloud and gets the results
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) Get(context context.Context, path URI, results interface{}) (correlationID string, err error) {
	return client.SendRequest(context, path, &request.Options{Method: http.MethodGet}, results)
}

// Delete sends a DELETE HTTP Request to GCloud and gets the results
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) Delete(context context.Context, path URI, results interface{}) (correlationID string, err error) {
	return client.SendRequest(context, path, &request.Options{Method: http.MethodDelete}, results)
}

// SendRequest sends a REST request to GCloud
//
// The Genesys Cloud Correlation ID is returned if available
func (client *Client) SendRequest(context context.Context, uri URI, options *request.Options, results interface{}) (correlationID string, err error) {
	log := client.GetLogger(context).Child(nil, "request")
	if options == nil {
		options = &request.Options{}
	}
	log = log.Record("method", options.Method)
	if uri.HasProtocol() {
		options.URL, err = uri.URL()
		log = log.Record("api", uri.String())
	} else if client.API == nil {
		return "", errors.ArgumentMissing.With("Client API")
	} else if !uri.HasPrefix("/api") {
		options.URL, err = client.API.Parse(NewURI("/api/v2/").Join(uri).String())
		log = log.Record("api", path.Join("/api/v2/", uri.String()))
	} else {
		options.URL, err = client.API.Parse(uri.String())
		log = log.Record("api", uri.String())
	}
	if err != nil {
		return "", errors.WithStack(APIError{Code: "url.parse", Message: err.Error()})
	}
	if len(options.Authorization) == 0 {
		if client.IsAuthorized() {
			options.Authorization = client.Grant.AccessToken().String()
		} else {
			correlationID, err = client.Login(context)
			if err != nil {
				return correlationID, errors.WithStack(err)
			}
			if !client.IsAuthorized() {
				return correlationID, errors.HTTPUnauthorized.WithStack()
			}
			options.Authorization = client.Grant.AccessToken().String()
		}
	}
	if options.Timeout == 0 {
		options.Timeout = client.RequestTimeout
	}
	if len(options.RetryableStatusCodes) == 0 {
		options.RetryableStatusCodes = []int{http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
	}

	options.Context = context
	options.Proxy = client.Proxy
	options.UserAgent = APP + " " + VERSION
	options.Logger = log
	options.ResponseBodyLogSize = 32768
	options.InterAttemptUseRetryAfter = true

	log.Record("payload", options.Payload).Debugf("Sending request to %s", options.URL)
	if value := context.Value("logPayload"); value != nil {
		if logPayload, ok := value.(bool); ok && logPayload {
			if payload, err := json.Marshal(options.Payload); err == nil {
				log.Scope("payload").Debugf("Payload: %s", string(payload))
			}
		}
	}
	start := time.Now()
	res, err := request.Send(options, results)
	duration := time.Since(start)
	log = log.Record("duration", duration)
	if res != nil {
		correlationID = res.Headers.Get("Genesys-Correlation-Id") // The new way
		if len(correlationID) == 0 {
			correlationID = res.Headers.Get("Inin-Correlation-Id") //The old way, back in the Interactiove Intelligence days
		}
		log = log.Record("genesys-correlation", correlationID)
	}
	if err != nil {
		urlError := &url.Error{}
		if errors.As(err, &urlError) {
			log.Errorf("URL Error", urlError)
			return correlationID, err
		}
		if errors.Is(err, errors.HTTPStatusTooManyRequests) {
			log.Errorf("Too many requests, retrying in %s seconds", res.Headers.Get("Retry-After"))
			retryAfter := time.Duration(core.Atoi(res.Headers.Get("Retry-After"), 0))*time.Second + 1*time.Second // just to stay on the safe side, add 1 second
			log.Debugf("Retry-After from headers (+1s safety net): %s", retryAfter)
			log.Infof("Waiting for %s before trying again", retryAfter)
			time.Sleep(retryAfter)
			log.Infof("Retrying request")
			return client.SendRequest(context, uri, options, results)
		}
		if errors.Is(err, errors.HTTPUnauthorized) && len(client.Grant.AccessToken().String()) > 0 {
			// This means our token most probably expired, we should try again without it
			log.Infof("Authorization Token is expired, we need to authenticate again")
			options.Authorization = ""
			client.Grant.AccessToken().Reset()
			return client.SendRequest(context, uri, options, results)
		}
		if errors.Is(err, errors.HTTPBadRequest) {
			log.Record("Request payload", options.Payload).Errorf("Bad Request from remote: %s", err.Error())
		}
		if errors.Is(err, errors.JSONUnmarshalError) {
			log.Errorf("Failed to unmarshal the response: %s", err.Error())
			log.Infof("Response payload: %s", res.Data)
			return correlationID, errors.Join(JSONUnmarshalError.SetCorrelationID(correlationID).WithParams(fmt.Sprintf("%T", results), map[string]string{"data": string(res.Data)}), err)
		}
		log.Errorf("Response payload: %s", res.Data)
		var simpleError struct {
			Error       string `json:"error"`
			Description string `json:"description"`
		}
		if jsonerr := res.UnmarshalContentJSON(&simpleError); jsonerr == nil && len(simpleError.Error) > 0 {
			return correlationID, APIError{Status: 500, Code: "generic", Message: simpleError.Error, MessageParams: map[string]string{"description": simpleError.Description}, CorrelationID: correlationID}.WithStack()
		}
		var details *errors.Error
		if errors.As(err, &details) {
			apiError := APIError{}
			if res != nil {
				if jsonerr := res.UnmarshalContentJSON(&apiError); jsonerr != nil {
					return correlationID, errors.Wrap(err, "Failed to extract an error from the response")
				}
				apiError.CorrelationID = correlationID
				return correlationID, apiError.WithStack()
			}
			// Sometimes we do not get a response with a Gcloud error, but a generic error
			apiError.Status = details.Code
			apiError.Code = details.ID
			apiError.CorrelationID = correlationID
			if strings.HasPrefix(apiError.Message, "authentication failed") {
				apiError.Status = errors.HTTPUnauthorized.Code
				apiError.Code = errors.HTTPUnauthorized.ID
			}
			return correlationID, apiError.WithStack()
		}
		return correlationID, errors.WithStack(err)
	}
	log.Debugf("Successfuly sent request in %s, correlation ID: %s", duration, correlationID)
	return correlationID, nil
}

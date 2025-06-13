package gcloudcx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Note: The declaration of ClientSuite is in client_test.go

func (suite *ClientSuite) TestCanSendGetRequest() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Get(context.Background(), "/path/to/resource", &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendPostRequest() {
	server := CreateTestServer(http.MethodPost, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Post(context.Background(), "/path/to/resource", struct{}{}, &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendPatchRequest() {
	server := CreateTestServer(http.MethodPatch, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Patch(context.Background(), "/path/to/resource", struct{}{}, &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendPutRequest() {
	server := CreateTestServer(http.MethodPut, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Put(context.Background(), "/path/to/resource", struct{}{}, &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendDeleteRequest() {
	server := CreateTestServer(http.MethodDelete, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Delete(context.Background(), "/path/to/resource", &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendRequestWithFullyQualifiedURL() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	serverURL := core.Must(url.Parse(server.URL))
	requestURL := core.Must(serverURL.Parse("/api/v2/path/to/resource"))
	correlationID, err := client.Get(context.Background(), gcloudcx.NewURI("%s", requestURL.String()), &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestCanSendRequestWithAPIPrefix() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.Get(context.Background(), "/api/v2/path/to/resource", &stuff)
	suite.Require().Nilf(err, "Failed to send GET Request: Error %s", err)
	suite.Assert().NotEmpty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestShouldNotSendRequestWithInvalidProtocol() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.SendRequest(context.Background(), "invalid://acme.com", nil, &stuff)
	suite.Require().NotNil(err, "Should not send request withan invalid URL")
	suite.Assert().Empty(correlationID, "Correlation ID should not be empty")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
	suite.Logger.Errorf("Expected error", err)
	var details *url.Error
	suite.Require().True(errors.As(err, &details), "err should contain an url.Error")
	suite.Assert().Equal(`unsupported protocol scheme "invalid"`, details.Unwrap().Error())
}

func (suite *ClientSuite) TestShouldNotSendRequestWithInvalidURL() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	stuff := struct{}{}
	correlationID, err := client.SendRequest(context.Background(), "http://wrong hostname.com", nil, &stuff)
	suite.Require().NotNil(err, "Should not send request withan invalid URL")
	suite.Logger.Errorf("Expected error", err)
	suite.Assert().Contains(err.Error(), "invalid character")
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
}

func (suite *ClientSuite) TestShouldNotSendRequestWithNoAPI() {
	server := CreateTestServer(http.MethodGet, "/api/v2/path/to/resource", suite.T())
	defer server.Close()

	client := CreateTestClient(server.URL, suite.Logger)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
	client.API = nil
	stuff := struct{}{}
	correlationID, err := client.SendRequest(context.Background(), "/path/to/resource", nil, &stuff)
	suite.Require().NotNil(err, "Should not send request without an API URL")
	suite.Logger.Errorf("Expected error", err)
	suite.Logger.Infof("Received Correlation ID: %s", correlationID)
	suite.Assert().True(errors.Is(err, errors.ArgumentMissing))
	var details errors.Error
	suite.Require().True(errors.As(err, &details), "err should contain an errors.Error")
	suite.Assert().Equal("Client API", details.What)
}

// Tool Stuff

func CreateTestServer(expectedMethod, expectedURL string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, "/api/v2/path/to/resource", r.URL.String())
		w.Header().Add("Inin-Correlation-Id", "12345")
		core.RespondWithJSON(w, http.StatusOK, struct{}{})
	}))
}

func CreateTestClient(serverURL string, log *logger.Logger) *gcloudcx.Client {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: "mypurecloud.com",
		Logger: log,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.New(),
		Secret:   "s3cr3t",
		Token: gcloudcx.AccessToken{
			Type:      "bearer",
			Token:     "F@k3T0k3nV@lu3",
			ExpiresOn: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	})
	if client != nil {
		client.Logger.Infof("Redirecting Client to Test Server at %s", serverURL)
		client.API = core.Must(url.Parse(serverURL))
		client.LoginURL = client.API
	}
	return client
}

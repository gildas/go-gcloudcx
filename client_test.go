package gcloudcx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *ClientSuite) SetupSuite() {
	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
			SourceInfo:  true,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))
}

func (suite *ClientSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *ClientSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *ClientSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *ClientSuite) TestCanInitialize() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: "mypurecloud.com",
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.New(),
		Secret:   "s3cr3t",
	})
	suite.Require().NotNil(client, "gcloudcx Client is nil")
}

func (suite *ClientSuite) TestCanInitializeWithoutOptions() {
	client := gcloudcx.NewClient(nil)
	suite.Require().NotNil(client, "GCloudCX Client is nil")
}

func (suite *ClientSuite) TestClientNotFoundErrorShouldBeBadRequest() {
	payload := `{"error":"invalid_client","description":"client not found","error_description":"client not found"}`
	var apiError gcloudcx.APIError

	err := json.Unmarshal([]byte(payload), &apiError)
	suite.Require().NoError(err, "Unmarshalling should have succeeded")
	suite.Logger.Errorf("Expected Error", apiError)
	suite.Assert().ErrorIs(apiError, gcloudcx.BadCredentialsError)
	suite.Assert().NotErrorIs(apiError, errors.RuntimeError)
	suite.Assert().Equal(gcloudcx.BadCredentialsError.Status, apiError.Status)
}

func (suite *ClientSuite) TestAuthFailedErrorShouldBeBadRequest() {
	payload := `{"error":"invalid_client","description":"authentication failed","error_description":"authentication failed"}`
	var apiError gcloudcx.APIError

	err := json.Unmarshal([]byte(payload), &apiError)
	suite.Require().NoError(err, "Unmarshalling should have succeeded")
	suite.Logger.Errorf("Expected Error", apiError)
	suite.Assert().ErrorIs(apiError, gcloudcx.BadCredentialsError)
	suite.Assert().NotErrorIs(apiError, errors.RuntimeError)
	suite.Assert().Equal(gcloudcx.BadCredentialsError.Status, apiError.Status)
}

func (suite *ClientSuite) TestCanLoginWithClientCredentials() {
	clientID := uuid.New()

	if value := core.GetEnvAsString("PURECLOUD_CLIENTID", ""); len(value) > 0 {
		clientID = uuid.MustParse(value)
	}
	
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"),
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "s3cr3t"),
	})
	err := client.Login(context.Background())
	suite.Require().NoError(err, "Login should have succeeded")
	suite.Require().NotEmpty(client.Grant.AccessToken(), "Access Token should not be empty")
}

func (suite *ClientSuite) TestShouldFailLoginWithInvalidClientCredentialsSecret() {
	clientID := uuid.New()

	if value := core.GetEnvAsString("PURECLOUD_CLIENTID", ""); len(value) > 0 {
		clientID = uuid.MustParse(value)
	}
	
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"),
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   "s3cr3t",
	})
	err := client.Login(context.Background())
	suite.Require().Error(err, "Login should have failed")
	suite.Logger.Errorf("Expected Error", err)
	suite.Assert().NotErrorIs(err, errors.RuntimeError)
	suite.Assert().ErrorIs(err, gcloudcx.BadCredentialsError)

	var apiError gcloudcx.APIError
	suite.Require().ErrorAs(err, &apiError, "Error should be an APIError")
	suite.Require().NotEmpty(apiError.MessageParams, "Error should have some parameters")
	suite.Assert().Equal("authentication failed", apiError.MessageParams["description"])
}

func (suite *ClientSuite) TestShouldFailLoginWithInvalidClientCredentialsClientID() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"),
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.New(), // The chances of this being a valid Client ID are very low
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "s3cr3t"),
	})
	err := client.Login(context.Background())
	suite.Require().Error(err, "Login should have failed")
	suite.Logger.Errorf("Expected Error", err)
	suite.Assert().NotErrorIs(err, errors.RuntimeError)
	suite.Assert().ErrorIs(err, gcloudcx.BadCredentialsError)

	var apiError gcloudcx.APIError
	suite.Require().ErrorAs(err, &apiError, "Error should be an APIError")
	suite.Require().NotEmpty(apiError.MessageParams, "Error should have some parameters")
	suite.Assert().Equal("client not found", apiError.MessageParams["description"])
}

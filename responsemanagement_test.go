package gcloudcx_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type ResponseManagementSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	LibraryName  string
	LibraryID    uuid.UUID
	ResponseName string
	ResponseID   uuid.UUID
	Client       *gcloudcx.Client
}

func TestResponseManagementSuite(t *testing.T) {
	suite.Run(t, new(ResponseManagementSuite))
}

func (suite *ResponseManagementSuite) TestCanFetchLibraryByID(){
	library := gcloudcx.ResponseManagementLibrary{}
	err := suite.Client.Fetch(context.Background(), &library, suite.LibraryID)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().Nilf(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.LibraryID, library.GetID(), "Library ID is not the same")
	suite.Assert().Equal(suite.LibraryName, library.String(), "Library Name is not the same")
	suite.Logger.Record("library", library).Infof("Library Details")
}

func (suite *ResponseManagementSuite) TestCanFetchLibraryByName(){
	library := gcloudcx.ResponseManagementLibrary{}
	err := suite.Client.Fetch(context.Background(), &library, suite.LibraryName)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().Nil(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.LibraryID, library.GetID(), "Library ID is not the same")
	suite.Assert().Equal(suite.LibraryName, library.String(), "Library Name is not the same")
	suite.Logger.Record("library", library).Infof("Library Details")
}

func (suite *ResponseManagementSuite) TestCanFetchResponseByID(){
	response := gcloudcx.ResponseManagementResponse{}
	err := suite.Client.Fetch(context.Background(), &response, suite.ResponseID)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().Nil(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.ResponseID, response.GetID(), "Client's Organization ID is not the same")
	suite.Assert().Equal(suite.ResponseName, response.String(), "Client's Organization Name is not the same")
	suite.Logger.Record("response", response).Infof("Response Details")
}

func (suite *ResponseManagementSuite) TestCanFetchResponseByName(){
	response := gcloudcx.ResponseManagementResponse{}
	err := suite.Client.Fetch(context.Background(), &response, suite.ResponseName)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().Nil(err, "Failed to fetch Response Management Response, Error: %s", err)
	suite.Assert().Equal(suite.ResponseID, response.GetID(), "Response ID is not the same")
	suite.Assert().Equal(suite.ResponseName, response.String(), "Response Name is not the same")
	suite.Logger.Record("response", response).Infof("Response Details")
}

func (suite *ResponseManagementSuite) TestShouldFailFetchingLibraryWithUnknownName() {
	library := gcloudcx.ResponseManagementLibrary{}
	err := suite.Client.Fetch(context.Background(), &library, "unknown library")
	suite.Require().NotNil(err, "Should have failed to fetch Response Management Library")
	suite.Logger.Errorf("Expected Failure", err)
	suite.Assert().ErrorIs(err, errors.NotFound, "Should have failed to fetch Response Management Library")
}

func (suite *ResponseManagementSuite) TestShouldFailFetchingResponseWithUnknownName() {
	response := gcloudcx.ResponseManagementLibrary{}
	err := suite.Client.Fetch(context.Background(), &response, "unknown response")
	suite.Require().NotNil(err, "Should have failed to fetch Response Management Response")
	suite.Logger.Errorf("Expected Failure", err)
	suite.Assert().ErrorIs(err, errors.NotFound, "Should have failed to fetch Response Management Response")
}

// Suite Tools

func (suite *ResponseManagementSuite) SetupSuite() {
	var err error
	var value string

	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	region := core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com")
	
	value = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_CLIENTID is not set")

	clientID, err := uuid.Parse(value)
	suite.Require().Nil(err, "PURECLOUD_CLIENTID is not a valid UUID")

	secret := core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	suite.Require().NotEmpty(secret, "PURECLOUD_CLIENTSECRET is not set")

	value = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_DEPLOYMENTID is not set")

	deploymentID, err := uuid.Parse(value)
	suite.Require().Nil(err, "PURECLOUD_DEPLOYMENTID is not a valid UUID")

	value = core.GetEnvAsString("RESPONSE_MANAGEMENT_LIBRARY_ID", "")
	suite.Require().NotEmpty(value, "RESPONSE_MANAGEMENT_LIBRARY_ID is not set in your environment")

	suite.LibraryID, err = uuid.Parse(value)
	suite.Require().Nil(err, "RESPONSE_MANAGEMENT_LIBRARY_ID is not a valid UUID")

	suite.LibraryName = core.GetEnvAsString("RESPONSE_MANAGEMENT_LIBRARY_NAME", "")
	suite.Require().NotEmpty(suite.LibraryName, "RESPONSE_MANAGEMENT_LIBRARY_NAME is not set in your environment")

	value = core.GetEnvAsString("RESPONSE_MANAGEMENT_RESPONSE_ID", "")
	suite.Require().NotEmpty(value, "RESPONSE_MANAGEMENT_RESPONSE_ID is not set in your environment")

	suite.ResponseID, err = uuid.Parse(value)
	suite.Require().Nil(err, "RESPONSE_MANAGEMENT_RESPONSE_ID is not a valid UUID")

	suite.ResponseName = core.GetEnvAsString("RESPONSE_MANAGEMENT_RESPONSE_NAME", "")
	suite.Require().NotEmpty(suite.ResponseName, "RESPONSE_MANAGEMENT_RESPONSE_NAME is not set in your environment")

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *ResponseManagementSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *ResponseManagementSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()

	// Reuse tokens as much as we can
	if !suite.Client.IsAuthorized() {
		suite.Logger.Infof("Client is not logged in...")
		err := suite.Client.Login(context.Background())
		suite.Require().Nil(err, "Failed to login")
		suite.Logger.Infof("Client is now logged in...")
	} else {
		suite.Logger.Infof("Client is already logged in...")
	}
}

func (suite *ResponseManagementSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}
//go:build integration
// +build integration

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

type LoginSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *gcloudcx.Client
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *LoginSuite) SetupSuite() {
	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:         fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:   true,
			SourceInfo:   true,
			FilterLevels: logger.NewLevelSet(logger.TRACE),
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		deploymentID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""))
	)
	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *LoginSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *LoginSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))

	suite.Start = time.Now()
}

func (suite *LoginSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *LoginSuite) TestCanLogin() {
	_, err := suite.Client.SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", "")),
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	}).Login(context.Background())
	suite.Assert().NoError(err, "Failed to login")
}

func (suite *LoginSuite) TestFailsLoginWithInvalidClientID() {
	_, err := suite.Client.LoginWithAuthorizationGrant(context.Background(), &gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.New(), // that UUID should not be anywhere in GCloud
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	})
	suite.Assert().Error(err, "Should have failed login in")

	var apierr gcloudcx.APIError
	ok := errors.As(err, &apierr)
	suite.Require().Truef(ok, "Error is not a gcloudcx.APIError, error: %+v", err)
	suite.Logger.Record("apierr", apierr).Errorf("API Error", err)
	suite.Assert().Equal(errors.HTTPBadRequest.Code, apierr.Status)
	suite.Assert().Equal("client not found: invalid_client", apierr.Error())
}

func (suite *LoginSuite) TestFailsLoginWithInvalidSecret() {
	_, err := suite.Client.LoginWithAuthorizationGrant(context.Background(), &gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", "")),
		Secret:   "WRONGSECRET",
	})
	suite.Assert().Error(err, "Should have failed login in")

	var apierr gcloudcx.APIError
	ok := errors.As(err, &apierr)
	suite.Require().Truef(ok, "Error is not a gcloudcx.APIError, error: %+v", err)
	suite.Logger.Record("apierr", apierr).Errorf("API Error", err)
	suite.Assert().Equal(errors.HTTPUnauthorized.Code, apierr.Status)
	suite.Assert().Equal("authentication failed: invalid_client", apierr.Error())
}

func (suite *LoginSuite) TestCanLoginWithClientCredentialsGrant() {
	_, err := suite.Client.LoginWithAuthorizationGrant(context.Background(), &gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", "")),
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	})
	suite.Assert().NoError(err, "Failed to login")
}

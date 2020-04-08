package purecloud_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/stretchr/testify/suite"

	purecloud "github.com/gildas/go-purecloud"
)

type LoginSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *purecloud.Client
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestCanLogin() {
	err := suite.Client.SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	}).Login()
	suite.Assert().Nil(err, "Failed to login")
}

func (suite *LoginSuite) TestFailsLoginWithInvalidGrant() {
	err := suite.Client.LoginWithAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: "DEADID",
		Secret:   "WRONGSECRET",
	})
	suite.Assert().NotNil(err, "Should have failed login in")

	var apierr purecloud.APIError
	ok := errors.As(err, &apierr)
	suite.Require().Truef(ok, "Error is not a purecloud.APIError, error: %+v", err)
	suite.Logger.Record("apierr", apierr).Errorf("API Error", err)
	suite.Assert().Equal(errors.HTTPUnauthorized.Code, apierr.Status)
	suite.Assert().Equal("authentication failed: invalid_client", apierr.Error())
}

func (suite *LoginSuite) TestCanLoginWithClientCredentialsGrant() {
	err := suite.Client.LoginWithAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	})
	suite.Assert().Nil(err, "Failed to login")
}

// Suite Tools

func (suite *LoginSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		deploymentID = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	)
	suite.Client = purecloud.NewClient(&purecloud.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	})
	suite.Require().NotNil(suite.Client, "PureCloud Client is nil")
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

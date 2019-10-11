package purecloud_test

import (
	"github.com/pkg/errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
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
	err := suite.Client.Login()
	suite.Assert().Nil(err, "Failed to login")
}

func (suite *LoginSuite) TestCanLoginWithGrant() {
	grant := &purecloud.ClientCredentialsGrant{
		ClientID:  core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		Secret:    core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	}
	err := suite.Client.LoginWithAuthorizationGrant(grant)
	suite.Assert().Nil(err, "Failed to login")
}

func (suite *LoginSuite) TestFailsLoginWithInvalidGrant() {
	grant := &purecloud.ClientCredentialsGrant{
		ClientID:  "DEADID",
		Secret:    "WRONGSECRET",
	}
	err := suite.Client.LoginWithAuthorizationGrant(grant)
	suite.Assert().NotNil(err, "Should have failed login in")

	apierr, ok := errors.Cause(err).(purecloud.APIError)
	suite.Require().True(ok, "Error is not a purecloud.APIError")
	suite.Logger.Record("apierr", apierr).Errorf("API Error", err)
	suite.Assert().Equal(400, apierr.Status)
	suite.Assert().Equal("bad.credentials", apierr.Code)
}

// Suite Tools

func (suite *LoginSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	logFolder := filepath.Join(".", "log")
	os.MkdirAll(logFolder, os.ModePerm)
	suite.Logger = logger.CreateWithDestination("test", fmt.Sprintf("file://%s/test-%s.log", logFolder, strings.ToLower(suite.Name)))
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID     = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		deploymentID = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	)
	suite.Client = purecloud.New(purecloud.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		ClientID:     clientID,
		ClientSecret: secret,
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
}

func (suite *LoginSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))

	suite.Start = time.Now()
}

func (suite *LoginSuite) AfterTest(suiteName, testName string) {
	duration := time.Now().Sub(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

package gcloudcx_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *gcloudcx.Client
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

// Suite Tools

func (suite *AuthSuite) SetupSuite() {
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

	var (
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID     = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		deploymentID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""))
	)

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

func (suite *AuthSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *AuthSuite) BeforeTest(suiteName, testName string) {
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

func (suite *AuthSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// Tests

func (suite *AuthSuite) TestCanCreateAuthScopeFromString() {
	var scope gcloudcx.AuthorizationScope
	
	suite.Assert().Equal("domain:entity:action", scope.With("domain", "entity", "action").String())
	suite.Assert().Equal("domain:entity:action", scope.With("domain:entity", "action").String())
	suite.Assert().Equal("domain:entity:action", scope.With("domain:entity:action").String())

	suite.Assert().Equal("domain:entity:*", scope.With("domain", "entity").String())
	suite.Assert().Equal("domain:entity:*", scope.With("domain:entity").String())

	suite.Assert().Equal("domain:*:*", scope.With("domain").String())

	suite.Assert().Equal("domain:entity:action", scope.With("domain:entity:action:extra").String())
	suite.Assert().Equal("domain:entity:action", scope.With("domain", "entity", "action", "extra").String())
}

func (suite *AuthSuite) TestCanUnmarshalAuthorizationSubject() {
	subject := gcloudcx.AuthorizationSubject{}
	err := LoadObject("authorization-subject.json", &subject)
	suite.Require().NoError(err, "Failed to load authorization subject, Error: %s", err)
}

func (suite *AuthSuite) TestCanUnmarshalAuthorizationSubjectWithDivisions() {
	subject := gcloudcx.AuthorizationSubject{}
	err := LoadObject("authorization-subject-with-divisions.json", &subject)
	suite.Require().NoError(err, "Failed to load authorization subject, Error: %s", err)
}

func (suite *AuthSuite) TestCanCheckScopes() {
	subject := gcloudcx.AuthorizationSubject{}
	err := LoadObject("authorization-subject.json", &subject)
	suite.Require().NoError(err, "Failed to load authorization subject, Error: %s", err)
	permitted, denied := subject.CheckScopes("routing:language:assign", "messaging:message", "processing:space:deploy")
	suite.Assert().Len(permitted, 2)
	suite.Assert().Len(denied, 1)
	suite.Assert().Contains(permitted, "routing:language:assign")
	suite.Assert().Contains(permitted, "messaging:message")
	suite.Assert().Contains(denied, "processing:space:deploy")
}
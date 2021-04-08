package purecloud_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	purecloud "github.com/gildas/go-purecloud"
)

type OpenMessagingSuite struct {
	suite.Suite
	Name string
	Logger *logger.Logger
	Start  time.Time

	Client *purecloud.Client
}

func TestOpenMessagingSuite(t *testing.T) {
	suite.Run(t, new(OpenMessagingSuite))
}

func (suite *OpenMessagingSuite) TestCanUnmarshal() {
	integration := purecloud.OpenMessagingIntegration{}
	err := Load("openmessagingintegration.json", &integration)
	if err != nil {
		suite.Logger.Errorf("Failed to Unmarshal", err)
	}
	suite.Require().Nil(err, "Failed to unmarshal OpenMessagingIntegration. %s", err)
	suite.Logger.Record("integration", integration).Infof("Got a integration")
	suite.Assert().NotEmpty(integration.ID)
	suite.Assert().NotEmpty(integration.CreatedBy.ID)
	suite.Assert().NotEmpty(integration.CreatedBy.SelfURI, "CreatedBy SelfURI should not be empty")
	suite.Assert().Equal(2021, integration.DateCreated.Year())
	suite.Assert().Equal(time.Month(4), integration.DateCreated.Month())
	suite.Assert().Equal(8, integration.DateCreated.Day())
	suite.Assert().NotEmpty(integration.ModifiedBy.ID)
	suite.Assert().NotEmpty(integration.ModifiedBy.SelfURI, "ModifiedBy SelfURI should not be empty")
	suite.Assert().Equal(2021, integration.DateModified.Year())
	suite.Assert().Equal(time.Month(4), integration.DateModified.Month())
	suite.Assert().Equal(8, integration.DateModified.Day())
	suite.Assert().Equal("TEST-GO-PURECLOUD", integration.Name)
	suite.Assert().Equal("DEADBEEF", integration.WebhookToken)
	suite.Require().NotNil(integration.WebhookURL)
	suite.Assert().Equal("https://www.acme.com/purecloud", integration.WebhookURL.String())
}

// Suite Tools

func (suite *OpenMessagingSuite) SetupSuite() {
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
		clientID     = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		token        = core.GetEnvAsString("PURECLOUD_CLIENTTOKEN", "")
		deploymentID = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	)
	suite.Client = purecloud.NewClient(&purecloud.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
		Token: purecloud.AccessToken{
			Type: "bearer",
			Token: token,
		},
	})
	suite.Require().NotNil(suite.Client, "PureCloud Client is nil")
}

func (suite *OpenMessagingSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *OpenMessagingSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))

	suite.Start = time.Now()

	// Reuse tokens as much as we can
	if !suite.Client.IsAuthorized() {
		suite.Logger.Infof("Client is not logged in...")
		err := suite.Client.Login()
		suite.Require().Nil(err, "Failed to login")
		suite.Logger.Infof("Client is now logged in...")
	} else {
		suite.Logger.Infof("Client is already logged in...")
	}
}

func (suite *OpenMessagingSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

package purecloud_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
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
	suite.Assert().Equal("34071108-1569-4cb0-9137-a326b8a9e815", integration.ID.String())
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

func (suite *OpenMessagingSuite) TestCanMarshal() {
	integration := purecloud.OpenMessagingIntegration{
		ID:               uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"),
		Name:             "TEST-GO-PURECLOUD",
		WebhookURL:       core.Must(url.Parse("https://www.acme.com/purecloud")).(*url.URL),
		WebhookToken:     "DEADBEEF",
		SupportedContent: &purecloud.AddressableEntityRef{
			ID:      uuid.MustParse("832066dd-6030-46b1-baeb-b89b681c6636"),
			SelfURI: "/api/v2/conversations/messaging/supportedcontent/832066dd-6030-46b1-baeb-b89b681c6636",
		},
		DateCreated:      time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		CreatedBy:        &purecloud.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		DateModified:     time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		ModifiedBy:       &purecloud.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		CreateStatus:     "Initiated",
		SelfURI:          "/api/v2/conversations/messaging/integrations/open/34071108-1569-4cb0-9137-a326b8a9e815",
	}

	data, err := json.Marshal(integration)
	suite.Require().Nil(err, "Failed to marshal OpenMessagingIntegration. %s", err)
	expected, err := LoadFile("openmessagingintegration.json")
	suite.Require().Nil(err, "Failed to Load Data. %s", err)
	suite.Assert().JSONEq(string(expected), string(data))
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
		clientID     = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		token        = core.GetEnvAsString("PURECLOUD_CLIENTTOKEN", "")
		deploymentID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""))
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

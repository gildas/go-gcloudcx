package gcloudcx_test

import (
	"context"
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

	"github.com/gildas/go-gcloudcx"
)

type OpenMessagingSuite struct {
	suite.Suite
	Name string
	Logger *logger.Logger
	Start  time.Time

	IntegrationID uuid.UUID
	Client        *gcloudcx.Client
}

func TestOpenMessagingSuite(t *testing.T) {
	suite.Run(t, new(OpenMessagingSuite))
}

func (suite *OpenMessagingSuite) TestCanFetchIntegrations() {
	integrations, err := suite.Client.FetchOpenMessagingIntegrations(context.Background())
	suite.Require().Nil(err, "Failed to fetch OpenMessaging Integrations")
	suite.T().Logf("Found %d integrations", len(integrations))
	if len(integrations) > 0 {
		for _, integration := range integrations {
			suite.Logger.Record("integration", integration).Infof("Got a integration")
			suite.Assert().NotEmpty(integration.ID)
			suite.Assert().NotEmpty(integration.Name)
			suite.Assert().NotNil(integration.WebhookURL, "WebhookURL should not be nil (%s)", integration.Name)
			suite.T().Logf("  Integration: %s (%s)", integration.Name, integration.ID)
		}
	}
}

func (suite *OpenMessagingSuite) TestCanCreateIntegration() {
	name := "TEST-GO-PURECLOUD"
	webhookURL, _ := url.Parse("https://www.genesys.com/gcloudcx")
	webhookToken := "DEADBEEF"
	integration, err := suite.Client.CreateOpenMessagingIntegration(context.Background(), name, webhookURL, webhookToken, nil)
	suite.Require().Nil(err, "Failed to create integration")
	suite.Logger.Record("integration", integration).Infof("Created a integration")
	for {
		if integration.IsCreated() {
			break
		}
		suite.Logger.Warnf("Integration %s is still in status: %s, waiting a bit", integration.ID, integration.CreateStatus)
		time.Sleep(time.Second)
		err = integration.Refresh(context.Background())
		suite.Require().Nil(err, "Failed to refresh integration")
	}
	suite.IntegrationID = integration.ID
}

func (suite *OpenMessagingSuite) TestCanDeleteIntegration() {
	suite.Require().NotNil(suite.IntegrationID, "IntegrationID should not be nil (TestCanCreateIntegration should run before this test)")
	integration := gcloudcx.OpenMessagingIntegration{}
	err := suite.Client.Fetch(context.Background(), &integration, suite.IntegrationID)
	suite.Require().Nilf(err, "Failed to fetch integration %s, Error: %s", suite.IntegrationID, err)
	suite.Logger.Record("integration", integration).Infof("Got a integration")
	suite.Require().True(integration.IsCreated(), "Integration should be created")
	err = integration.Delete(context.Background())
	suite.Require().Nilf(err, "Failed to delete integration %s, Error: %s", suite.IntegrationID, err)
	err = suite.Client.Fetch(context.Background(), &integration, suite.IntegrationID)
	suite.Require().NotNil(err, "Integration should not exist anymore")
	suite.IntegrationID = uuid.Nil
}

func (suite *OpenMessagingSuite) TestCanUnmarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{}
	err := LoadObject("openmessagingintegration.json", &integration)
	if err != nil {
		suite.Logger.Errorf("Failed to Unmarshal", err)
	}
	suite.Require().Nilf(err, "Failed to unmarshal OpenMessagingIntegration. %s", err)
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
	suite.Assert().Equal("https://www.acme.com/gcloudcx", integration.WebhookURL.String())
}

func (suite *OpenMessagingSuite) TestCanMarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{
		ID:               uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"),
		Name:             "TEST-GO-PURECLOUD",
		WebhookURL:       core.Must(url.Parse("https://www.acme.com/gcloudcx")).(*url.URL),
		WebhookToken:     "DEADBEEF",
		SupportedContent: &gcloudcx.AddressableEntityRef{
			ID:      uuid.MustParse("832066dd-6030-46b1-baeb-b89b681c6636"),
			SelfURI: "/api/v2/conversations/messaging/supportedcontent/832066dd-6030-46b1-baeb-b89b681c6636",
		},
		DateCreated:      time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		CreatedBy:        &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		DateModified:     time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		ModifiedBy:       &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		CreateStatus:     "Initiated",
		SelfURI:          "/api/v2/conversations/messaging/integrations/open/34071108-1569-4cb0-9137-a326b8a9e815",
	}

	data, err := json.Marshal(integration)
	suite.Require().Nilf(err, "Failed to marshal OpenMessagingIntegration. %s", err)
	expected, err := LoadFile("openmessagingintegration.json")
	suite.Require().Nilf(err, "Failed to Load Data. %s", err)
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalIntegrationWithInvalidJSON() {
	var err error

	integration := gcloudcx.OpenMessagingIntegration{}
	err = json.Unmarshal([]byte(`{"Name": 15}`), &integration)
	suite.Assert().NotNil(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageChannel() {
	channel := gcloudcx.OpenMessageChannel{}
	err := LoadObject("openmessaging-channel.json", &channel)
	suite.Require().Nilf(err, "Failed to unmarshal OpenMessageChannel. %s", err)
	suite.Assert().Equal("Open", channel.Platform)
	suite.Assert().Equal("Private", channel.Type)
	suite.Assert().Equal("gmAy9zNkhf4ermFvHH9mB5", channel.MessageID)
	suite.Assert().Equal(time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC), channel.Time)
	suite.Assert().Equal("edce4efa-4abf-468b-ada7-cd6d35e7bbaf", channel.To.ID)
	suite.Assert().Equal("Email", channel.From.Type)
	suite.Assert().Equal("abcdef12345", channel.From.ID)
	suite.Assert().Equal("Bob", channel.From.Firstname)
	suite.Assert().Equal("Minion", channel.From.Lastname)
	suite.Assert().Equal("Bobby", channel.From.Nickname)
}

func (suite *OpenMessagingSuite) TestCanMarshalOpenMessageChannel() {
	channel := gcloudcx.NewOpenMessageChannel(
		"gmAy9zNkhf4ermFvHH9mB5",
		&gcloudcx.OpenMessageTo{ ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
		&gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		},
	)
	channel.Time = time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC)

	data, err := json.Marshal(channel)
	suite.Require().Nilf(err, "Failed to marshal OpenMessageChannel. %s", err)
	suite.Require().NotNil(data, "Marshaled data should not be nil")
	expected, err := LoadFile("openmessaging-channel.json")
	suite.Require().Nilf(err, "Failed to Load Data. %s", err)
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessage() {
	message := gcloudcx.OpenMessageText{
		ID:        "12345678",
		Direction: "inbound",
		Text:      "text message",
	}
	message.Channel = gcloudcx.NewOpenMessageChannel(
		"gmAy9zNkhf4ermFvHH9mB5",
		&gcloudcx.OpenMessageTo{ ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
		&gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		},
	)
	message.Channel.Time = time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC)

	suite.Logger.Record("message", message).Infof("message")
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessageChannel() {
	channel := gcloudcx.NewOpenMessageChannel(
		"gmAy9zNkhf4ermFvHH9mB5",
		&gcloudcx.OpenMessageTo{ ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
		&gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		},
	)
	channel.Time = time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC)
	suite.Logger.Record("channel", channel).Infof("channel")
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessageFrom() {
	from := gcloudcx.OpenMessageFrom{
		ID:        "abcdef12345",
		Type:      "Email",
		Firstname: "Bob",
		Lastname:  "Minion",
		Nickname:  "Bobby",
	}
	suite.Logger.Record("from", from).Infof("from")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalChannelWithInvalidJSON() {
	var err error

	channel := gcloudcx.OpenMessageChannel{}
	err = json.Unmarshal([]byte(`{"Platform": 2}`), &channel)
	suite.Assert().NotNil(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalFromWithInvalidJSON() {
	var err error

	from := gcloudcx.OpenMessageFrom{}
	err = json.Unmarshal([]byte(`{"idType": 3}`), &from)
	suite.Assert().NotNil(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalMessageWithInvalidJSON() {
	_, err := gcloudcx.UnmarshalOpenMessage([]byte(`{"Direction": 6}`))
	suite.Assert().NotNil(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestCanStringifyIntegration() {
	id := uuid.New()
	integration := gcloudcx.OpenMessagingIntegration{
		ID:   id,
		Name: "Hello",
	}
	suite.Assert().Equal("Hello", integration.String())
	integration.Name = ""
	suite.Assert().Equal(id.String(), integration.String())
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
	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
		Token: gcloudcx.AccessToken{
			Type: "bearer",
			Token: token,
		},
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
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
		err := suite.Client.Login(context.Background())
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

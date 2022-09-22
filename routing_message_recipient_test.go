package gcloudcx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type RoutingMessageRecipientSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	IntegrationID uuid.UUID
	Client        *gcloudcx.Client
}

func TestRoutingMessageRecipientSuite(t *testing.T) {
	suite.Run(t, new(RoutingMessageRecipientSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *RoutingMessageRecipientSuite) SetupSuite() {
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
			Type:  "bearer",
			Token: token,
		},
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *RoutingMessageRecipientSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *RoutingMessageRecipientSuite) BeforeTest(suiteName, testName string) {
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

func (suite *RoutingMessageRecipientSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *RoutingMessageRecipientSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *RoutingMessageRecipientSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *RoutingMessageRecipientSuite) TestCanMarshal() {
	expected := suite.LoadTestData("routing-message-recipient.json")
	recipient := gcloudcx.RoutingMessageRecipient{
		ID:            uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"),
		Name:          "TEST-GO-PURECLOUD",
		MessengerType: "open",
		Flow: &gcloudcx.Flow{
			ID:   uuid.MustParse("900fa1cb-427b-4ae3-9439-079ac3f07d56"),
			Name: "Gildas-TestOpenMessaging",
		},
		DateCreated:  time.Date(2021, 4, 8, 3, 12, 7, 888000000, time.UTC),
		CreatedBy:    &gcloudcx.User{ID: uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e")},
		DateModified: time.Date(2021, 4, 8, 3, 12, 7, 888000000, time.UTC),
		ModifiedBy:   &gcloudcx.User{ID: uuid.MustParse("2229bd78-a6e4-412f-b789-ef70f447e5db")},
	}
	payload, err := json.Marshal(recipient)
	suite.Require().NoError(err, "Failed to Marshal")
	suite.Require().JSONEq(string(expected), string(payload), "Payload does not match")
}

func (suite *RoutingMessageRecipientSuite) TestCanUnmarshal() {
	var recipient gcloudcx.RoutingMessageRecipient
	err := suite.UnmarshalData("routing-message-recipient.json", &recipient)
	suite.Require().NoErrorf(err, "Failed to Unmarshal Data. %s", err)
	suite.Require().NotNil(recipient, "Recipient is nil")
	suite.Require().Implements((*core.Identifiable)(nil), recipient)
	suite.Assert().Equal("34071108-1569-4cb0-9137-a326b8a9e815", recipient.ID.String())
	suite.Assert().Equal("TEST-GO-PURECLOUD", recipient.Name)
	suite.Assert().Equal("open", recipient.MessengerType)
	suite.Assert().Equal("3e23b1b3-325f-4fbd-8fe0-e88416850c0e", recipient.CreatedBy.ID.String())
	suite.Assert().Equal("2021-04-08T03:12:07Z", recipient.DateCreated.Format(time.RFC3339))
	suite.Assert().Equal("2229bd78-a6e4-412f-b789-ef70f447e5db", recipient.ModifiedBy.ID.String())
	suite.Assert().Equal("2021-04-08T03:12:07Z", recipient.DateModified.Format(time.RFC3339))
	suite.Require().Implements((*gcloudcx.Addressable)(nil), recipient)
	suite.Assert().Equal(gcloudcx.NewURI("/api/v2/routing/message/recipients/%s", recipient.GetID()), recipient.GetURI())
}

func (suite *RoutingMessageRecipientSuite) TestCanFetchByID() {
	id := uuid.MustParse("2be1fcc8-4f7d-406d-accc-43be454e2f14")
	recipient, err := gcloudcx.Fetch[gcloudcx.RoutingMessageRecipient](context.Background(), suite.Client, id)
	suite.Require().NoErrorf(err, "Failed to fetch Routing Message Recipient %s. %s", id, err)
	suite.Assert().Equal(id, recipient.GetID())
	suite.Assert().Equal("GILDAS-OpenMessaging Integration Test-Viber", recipient.Name)
	suite.Assert().Equal("open", recipient.MessengerType)
	suite.Require().NotNil(recipient.Flow, "Recipient should have a Flow")
	suite.Assert().Equal("Gildas-TestOpenMessaging", recipient.Flow.Name)
}

func (suite *RoutingMessageRecipientSuite) TestCanFetchByName() {
	name := "GILDAS-OpenMessaging Integration Test-Viber"
	match := func(recipient gcloudcx.RoutingMessageRecipient) bool {
		return recipient.Name == name
	}
	recipient, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	suite.Require().NoErrorf(err, "Failed to fetch Routing Message Recipient %s. %s", name, err)
	suite.Assert().Equal(uuid.MustParse("2be1fcc8-4f7d-406d-accc-43be454e2f14"), recipient.GetID())
	suite.Assert().Equal(name, recipient.Name)
	suite.Assert().Equal("open", recipient.MessengerType)
	suite.Require().NotNil(recipient.Flow, "Recipient should have a Flow")
	suite.Assert().Equal("Gildas-TestOpenMessaging", recipient.Flow.Name)
}

func (suite *RoutingMessageRecipientSuite) TestCanFetchAll() {
	recipients, err := gcloudcx.FetchAll[gcloudcx.RoutingMessageRecipient](context.Background(), suite.Client, gcloudcx.Query{"messengerType": "open"})
	suite.Require().NoError(err, "Failed to fetch Routing Message Recipients")
	suite.Require().NotEmpty(recipients, "No Routing Message Recipients")
	suite.Logger.Infof("Found %d Routing Message Recipients", len(recipients))
	suite.Assert().Greater(len(recipients), 25, "Not enough Routing Message Recipients")
	for _, recipient := range recipients {
		suite.Logger.Record("recipient", recipient).Infof("Got a Routing Message Recipient")
		suite.Assert().NotEmpty(recipient.ID)
		suite.Assert().NotEmpty(recipient.Name)
	}
}

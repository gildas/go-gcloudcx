package gcloudcx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
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

type OpenMessagingSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	IntegrationID   uuid.UUID
	IntegrationName string
	Client          *gcloudcx.Client
}

func TestOpenMessagingSuite(t *testing.T) {
	suite.Run(t, new(OpenMessagingSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *OpenMessagingSuite) SetupSuite() {
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
	suite.IntegrationName = core.GetEnvAsString("PURECLOUD_INTEGRATION_NAME", "TEST-GO-PURECLOUD")
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
		correlationID, err := suite.Client.Login(context.Background())
		suite.Require().NoError(err, "Failed to login")
		suite.Logger.Infof("Correlation: %s", correlationID)
		suite.Logger.Infof("Client is now logged in...")
	} else {
		suite.Logger.Infof("Client is already logged in...")
	}
}

func (suite *OpenMessagingSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *OpenMessagingSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *OpenMessagingSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

func (suite *OpenMessagingSuite) LogLineEqual(line string, records map[string]string) {
	rex_records := make(map[string]*regexp.Regexp)
	for key, value := range records {
		rex_records[key] = regexp.MustCompile(value)
	}

	properties := map[string]interface{}{}
	err := json.Unmarshal([]byte(line), &properties)
	suite.Require().NoError(err, "Could not unmarshal line, error: %s", err)

	for key, rex := range rex_records {
		suite.Assert().Contains(properties, key, "The line does not contain the key %s", key)
		if value, found := properties[key]; found {
			var stringvalue string
			switch actual := value.(type) {
			case string:
				stringvalue = actual
			case int, int8, int16, int32, int64:
				stringvalue = strconv.FormatInt(value.(int64), 10)
			case uint, uint8, uint16, uint32, uint64:
				stringvalue = strconv.FormatUint(value.(uint64), 10)
			case float32, float64:
				stringvalue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			case fmt.Stringer:
				stringvalue = actual.String()
			case map[string]interface{}:
				stringvalue = fmt.Sprintf("%v", value)
			default:
				suite.Failf(fmt.Sprintf("The value of the key %s cannot be casted to string", key), "Type: %s", reflect.TypeOf(value))
			}
			suite.Assert().Truef(rex.MatchString(stringvalue), `Key "%s": the value %v does not match the regex /%s/`, key, value, rex)
		}
	}

	for key := range properties {
		suite.Assert().Contains(rex_records, key, "The line contains the extra key %s", key)
	}
}

func CaptureStdout(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	defer func(stdout *os.File) {
		os.Stdout = stdout
	}(os.Stdout)
	os.Stdout = writer

	f()
	writer.Close()

	output := bytes.Buffer{}
	_, _ = io.Copy(&output, reader)
	return output.String()
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *OpenMessagingSuite) TestCan00CreateIntegration() {
	webhookURL, _ := url.Parse("https://www.genesys.com/gcloudcx")
	webhookToken := "DEADBEEF"
	integration, correlationID, err := suite.Client.CreateOpenMessagingIntegration(context.Background(), suite.IntegrationName, webhookURL, webhookToken, nil)
	suite.Require().NoError(err, "Failed to create integration")
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Logger.Record("integration", integration).Infof("Created a integration")
	for {
		if integration.IsCreated() {
			break
		}
		suite.Logger.Warnf("Integration %s is still in status: %s, waiting a bit", integration.ID, integration.CreateStatus)
		time.Sleep(time.Second)
		correlationID, err = integration.Refresh(context.Background())
		suite.Require().NoError(err, "Failed to refresh integration")
		suite.Logger.Infof("Correlation: %s", correlationID)
	}
	suite.IntegrationID = integration.ID
}

func (suite *OpenMessagingSuite) TestCanFetchByID() {
	integration, correlationID, err := gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().NoErrorf(err, "Failed to fetch Open Messaging Integration %s. %s", suite.IntegrationID, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Assert().Equal(suite.IntegrationID, integration.ID)
	suite.Assert().Equal(suite.IntegrationName, integration.Name)
}

func (suite *OpenMessagingSuite) TestCanFetchByName() {
	match := func(integration gcloudcx.OpenMessagingIntegration) bool {
		return integration.Name == suite.IntegrationName
	}
	integration, correlationID, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	suite.Require().NoErrorf(err, "Failed to fetch Open Messaging Integration %s. %s", suite.IntegrationName, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Assert().Equal(suite.IntegrationID, integration.ID)
	suite.Assert().Equal(suite.IntegrationName, integration.Name)
}

func (suite *OpenMessagingSuite) TestCanFetchIntegrations() {
	integrations, correlationID, err := gcloudcx.FetchAll[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client)
	suite.Require().NoError(err, "Failed to fetch OpenMessaging Integrations")
	suite.Logger.Infof("Correlation: %s", correlationID)
	if len(integrations) > 0 {
		for _, integration := range integrations {
			suite.Logger.Record("integration", integration).Infof("Got a integration")
			suite.Assert().NotEmpty(integration.ID)
			suite.Assert().NotEmpty(integration.Name)
			suite.Assert().NotNil(integration.WebhookURL, "WebhookURL should not be nil (%s)", integration.Name)
		}
	}
}

func (suite *OpenMessagingSuite) TestCanZZDeleteIntegration() {
	suite.Require().NotNil(suite.IntegrationID, "IntegrationID should not be nil (TestCanCreateIntegration should run before this test)")
	integration, correlationID, err := gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().NoErrorf(err, "Failed to fetch integration %s, Error: %s", suite.IntegrationID, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Logger.Record("integration", integration).Infof("Got a integration")
	suite.Require().True(integration.IsCreated(), "Integration should be created")
	correlationID, err = integration.Delete(context.Background())
	suite.Require().NoErrorf(err, "Failed to delete integration %s, Error: %s", suite.IntegrationID, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	_, correlationID, err = gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().Error(err, "Integration should not exist anymore")
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Assert().ErrorIsf(err, gcloudcx.NotFoundError, "Expected NotFoundError, got %s", err)
	suite.Assert().Truef(errors.Is(err, gcloudcx.NotFoundError), "Expected NotFoundError, got %s", err)
	details := gcloudcx.NotFoundError.Clone()
	suite.Require().ErrorAsf(err, &details, "Expected NotFoundError but got %s", err)
	suite.IntegrationID = uuid.Nil
}

func (suite *OpenMessagingSuite) TestCanMarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{
		ID:           uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"),
		Name:         "TEST-GO-PURECLOUD",
		WebhookURL:   core.Must(url.Parse("https://www.acme.com/gcloudcx")),
		WebhookToken: "DEADBEEF",
		SupportedContent: &gcloudcx.OpenMessagingSupportedContent{
			ID:      "webMessagingDefault",
			SelfURI: "/api/v2/conversations/messaging/supportedcontent/webMessagingDefault",
		},
		MessagingSetting: &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("f61e2044-d338-4582-a2ed-0fbda7139aeb"),
			SelfURI: "/api/v2/conversations/messaging/settings/f61e2044-d338-4582-a2ed-0fbda7139aeb",
		},
		Recipient: &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("6a55dd3c-e148-481f-abc4-249db2939fa3"),
			SelfURI: "/api/v2/routing/message/recipients/6a55dd3c-e148-481f-abc4-249db2939fa3",
		},
		Status:      "active",
		DateCreated: time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		CreatedBy: &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		DateModified: time.Date(2021, time.April, 8, 3, 12, 7, 888000000, time.UTC),
		ModifiedBy: &gcloudcx.DomainEntityRef{
			ID:      uuid.MustParse("3e23b1b3-325f-4fbd-8fe0-e88416850c0e"),
			SelfURI: "/api/v2/users/3e23b1b3-325f-4fbd-8fe0-e88416850c0e",
		},
		CreateStatus: "Completed",
	}

	data, err := json.Marshal(integration)
	suite.Require().NoErrorf(err, "Failed to marshal OpenMessagingIntegration. %s", err)
	expected := suite.LoadTestData("openmessaging-integration.json")
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *OpenMessagingSuite) TestCanMarshalOpenMessageChannel() {
	channel := gcloudcx.OpenMessageChannel{
		Platform:  "Open",
		Type:      "Private",
		MessageID: "gmAy9zNkhf4ermFvHH9mB5",
		To:        &gcloudcx.OpenMessageTo{ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
		From: &gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		},
		Time: time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC),
	}

	data, err := json.Marshal(channel)
	suite.Require().NoErrorf(err, "Failed to marshal OpenMessageChannel. %s", err)
	suite.Require().NotNil(data, "Marshaled data should not be nil")
	expected := suite.LoadTestData("openmessaging-channel.json")
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *OpenMessagingSuite) TestCanMarshalTypingEvent() {
	event := gcloudcx.OpenMessageEvents{
		ID: "c327c2078ca056db130c55ce648d9fa2",
		Channel: gcloudcx.OpenMessageChannel{
			ID:       uuid.MustParse("73cb7fb7-c2db-4996-88f3-0a83a4fea1da"),
			Platform: "Open",
			Type:     "Private",
			From: &gcloudcx.OpenMessageFrom{
				ID:        "abcdef12345@socialmedia",
				Type:      "Email",
				Firstname: "Bob",
				Lastname:  "Minion",
				Nickname:  "Bobby",
			},
			To: &gcloudcx.OpenMessageTo{
				ID: "73cb7fb7-c2db-4996-88f3-0a83a4fea1da",
			},
			Time: time.Date(2021, 2, 1, 15, 4, 5, 0, time.UTC),
		},
		Direction: "Inbound",
		Events: []gcloudcx.OpenMessageEvent{
			gcloudcx.OpenMessageTypingEvent{IsTyping: true},
		},
		ConversationID: uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"),
	}
	payload, err := json.Marshal(event)
	suite.Require().NoErrorf(err, "Failed to marshal OpenMessageEvents. %s", err)
	expected := suite.LoadTestData("openmessaging-inbound-event-typing.json")
	suite.Require().JSONEq(string(expected), string(payload))
}

func (suite *OpenMessagingSuite) TestCanUnmarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{}
	err := suite.UnmarshalData("openmessaging-integration.json", &integration)
	if err != nil {
		suite.Logger.Errorf("Failed to Unmarshal", err)
	}
	suite.Require().NoErrorf(err, "Failed to unmarshal OpenMessagingIntegration. %s", err)
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
	suite.Require().NotNil(integration.SupportedContent)
	suite.Assert().Equal("webMessagingDefault", integration.SupportedContent.ID)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageChannel() {
	channel := gcloudcx.OpenMessageChannel{}
	err := suite.UnmarshalData("openmessaging-channel.json", &channel)
	suite.Require().NoErrorf(err, "Failed to unmarshal OpenMessageChannel. %s", err)
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

func (suite *OpenMessagingSuite) TestCanUnmarshalInboundTypingEvent() {
	payload := suite.LoadTestData("openmessaging-inbound-event-typing.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageEvents)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageEvents, but was %T", message)
	suite.Require().NotNil(actual, "Unmarshaled message should not be nil")
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("73cb7fb7-c2db-4996-88f3-0a83a4fea1da"), actual.Channel.ID)
	suite.Assert().Equal("Inbound", actual.Direction)
	suite.Require().Len(actual.Events, 1, "Unmarshaled message should have 1 event")

	messageEvent, ok := actual.Events[0].(*gcloudcx.OpenMessageTypingEvent)
	suite.Require().True(ok, "Unmarshaled message event should be of type *OpenMessageTypingEvent, but was %T", actual.Events[0])
	suite.Assert().True(messageEvent.IsTyping)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalReceipt() {
	payload := suite.LoadTestData("openmessaging-inbound-receipt-delivered.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageReceipt)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageReceipt, but was %T", message)
	suite.Require().NotNil(actual, "Unmarshaled message should not be nil")
	suite.Assert().Equal(uuid.MustParse("73cb7fb7-c2db-4996-88f3-0a83a4fea1da"), actual.Channel.ID)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Assert().False(actual.IsFailed(), "Receipt should not be failed")
	suite.Assert().Equal("Delivered", actual.Status)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalReceiptWithErrors() {
	payload := suite.LoadTestData("openmessaging-inbound-receipt-failure.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageReceipt)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageReceipt, but was %T", message)
	suite.Require().NotNil(actual, "Unmarshaled message should not be nil")
	suite.Assert().Equal(uuid.MustParse("73cb7fb7-c2db-4996-88f3-0a83a4fea1da"), actual.Channel.ID)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Assert().True(actual.IsFailed(), "Receipt should be failed")
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)

	err = actual.AsError()
	suite.Require().Error(err, "Receipt should be convert to an error")
	suite.ErrorIs(err, gcloudcx.GeneralError, "Receipt should convert to a GeneralError")
	suite.ErrorIs(err, gcloudcx.RateLimited, "Receipt should convert to a RateLimited")
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOutboundTypingEvent() {
	payload := suite.LoadTestData("openmessaging-outbound-event-typing.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageEvents)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageEvents, but was %T", message)
	suite.Require().NotNil(actual, "Unmarshaled message should not be nil")
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("73cb7fb7-c2db-4996-88f3-0a83a4fea1da"), actual.Channel.ID)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Require().Len(actual.Events, 1, "Unmarshaled message should have 1 event")

	messageEvent, ok := actual.Events[0].(*gcloudcx.OpenMessageTypingEvent)
	suite.Require().True(ok, "Unmarshaled message event should be of type *OpenMessageTypingEvent, but was %T", actual.Events[0])
	suite.Assert().True(messageEvent.IsTyping)
	suite.Assert().Equal(5*time.Second, messageEvent.Duration)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOutboundTextMessage() {
	payload := suite.LoadTestData("openmessaging-outbound-text.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageText)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageText, but was %T", message)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"), actual.ConversationID)
	suite.Assert().Equal("Hello World", actual.Text)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Assert().Equal("Human", actual.OriginatingEntity)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageStructuredWithNotification() {
	payload := suite.LoadTestData("openmessaging-outbound-notification.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageStructured)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageStructured, but was %T", message)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"), actual.ConversationID)
	suite.Assert().Equal("Hi Happy, How can I help you?", actual.Text)

	suite.Require().NotEmpty(actual.Content, "Content should not be empty")
	content := actual.Content[0]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("Notification", content.GetType())

	notification, ok := content.(*gcloudcx.NormalizedMessageNotificationContent)
	suite.Require().True(ok, "Content should be of type OpenMessageNotificationContent, but was %T", content)
	suite.Require().NotNil(notification, "Notification should not be nil")

	suite.Require().NotNil(notification.Header, "Notification Header should not be nil")
	suite.Assert().Equal("Hello", notification.Header.Text)

	suite.Assert().Equal("Hi Happy, How can I help you?", notification.Body.Text)
	suite.Require().NotEmpty(notification.Body.Parameters, "Notification Body Parameters should not be empty")
	value, found := notification.Body.Parameters["name"]
	suite.Require().True(found, "Notification Body Parameters should contain 'name'")
	suite.Assert().Equal("Happy", value)

	suite.Require().NotEmpty(notification.Footer, "Notification Footer should not be empty")
	suite.Assert().Equal("Goodbye", notification.Footer.Text)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageStructuredWithCarousel() {
	payload := suite.LoadTestData("openmessaging-inbound-carousel.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageStructured)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageStructured, but was %T", message)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"), actual.ConversationID)
	suite.Require().NotEmpty(actual.Content, "Content should not be empty")

	content := actual.Content[0]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("Carousel", content.GetType())

	carousel, ok := content.(*gcloudcx.NormalizedMessageCarouselContent)
	suite.Require().True(ok, "Content should be of type OpenMessageCarouselContent, but was %T", content)
	suite.Require().NotNil(carousel, "Carousel should not be nil")

	suite.Require().Len(carousel.Cards, 3, "Carousel should contain 3 cards")

	card1 := carousel.Cards[0]
	suite.Require().NotNil(card1, "Card 1 should not be nil")
	suite.Assert().Equal("Card 1", card1.Title)
	suite.Assert().Equal("", card1.Description)
	suite.Require().NotNil(card1.ImageURL, "Card 1 ImageURL should not be nil")
	suite.Assert().Equal("https://www.acme.com/image1.png", card1.ImageURL.String())
	suite.Assert().Nil(card1.VideoURL)
	suite.Require().Len(card1.Actions, 3, "Card 1 should have 3 actions")
	suite.Assert().Equal("Postback", card1.Actions[0].GetType())
	suite.Assert().Equal("Option1", card1.Actions[0].Text)
	suite.Assert().Equal("Option1", card1.Actions[0].Payload)
	suite.Assert().Equal("Postback", card1.Actions[1].GetType())
	suite.Assert().Equal("Option2", card1.Actions[1].Text)
	suite.Assert().Equal("Option2", card1.Actions[1].Payload)
	suite.Assert().Equal("Postback", card1.Actions[2].GetType())
	suite.Assert().Equal("Option3", card1.Actions[2].Text)
	suite.Assert().Equal("Option3", card1.Actions[2].Payload)

	card2 := carousel.Cards[1]
	suite.Assert().Equal("Card 2", card2.Title)
	suite.Assert().Nil(card2.ImageURL, "Card 2 ImageURL should be nil")
	suite.Assert().Nil(card2.VideoURL, "Card 2 VideoURL should be nil")
	suite.Require().Len(card2.Actions, 3, "Card 2 should have 3 actions")
	suite.Assert().Equal("Postback", card2.Actions[0].GetType())
	suite.Assert().Equal("Option4", card2.Actions[0].Text)
	suite.Assert().Equal("Option4", card2.Actions[0].Payload)
	suite.Assert().Equal("Postback", card2.Actions[1].GetType())
	suite.Assert().Equal("Option5", card2.Actions[1].Text)
	suite.Assert().Equal("Option5", card2.Actions[1].Payload)
	suite.Assert().Equal("Postback", card2.Actions[2].GetType())
	suite.Assert().Equal("Option6", card2.Actions[2].Text)
	suite.Assert().Equal("Option6", card2.Actions[2].Payload)

	card3 := carousel.Cards[2]
	suite.Assert().Equal("Card 3", card3.Title)
	suite.Assert().Nil(card3.ImageURL, "Card 3 ImageURL should be nil")
	suite.Assert().Nil(card3.VideoURL, "Card 3 VideoURL should be nil")
	suite.Require().Len(card3.Actions, 3, "Card 3 should have 3 actions")
	suite.Assert().Equal("Postback", card3.Actions[0].GetType())
	suite.Assert().Equal("Option7", card3.Actions[0].Text)
	suite.Assert().Equal("Option7", card3.Actions[0].Payload)
	suite.Assert().Equal("Postback", card3.Actions[1].GetType())
	suite.Assert().Equal("Option8", card3.Actions[1].Text)
	suite.Assert().Equal("Option8", card3.Actions[1].Payload)
	suite.Assert().Equal("Postback", card3.Actions[2].GetType())
	suite.Assert().Equal("Option9", card3.Actions[2].Text)
	suite.Assert().Equal("Option9", card3.Actions[2].Payload)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageStructuredWithQuickReply() {
	payload := suite.LoadTestData("openmessaging-inbound-quickreply.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageStructured)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageStructured, but was %T", message)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"), actual.ConversationID)
	suite.Assert().Equal("Do you want to proceed?", actual.Text)
	suite.Require().Len(actual.Content, 2, "Content should 2 items")

	content := actual.Content[0]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("QuickReply", content.GetType())

	quickreply, ok := content.(*gcloudcx.NormalizedMessageQuickReplyContent)
	suite.Require().True(ok, "Content should be of type OpenMessageQuickReplyContent, but was %T", content)
	suite.Require().NotNil(quickreply, "QuickReply should not be nil")
	suite.Assert().Equal("Yes", quickreply.Text)
	suite.Assert().Equal("Yes", quickreply.Payload)
	suite.Assert().Equal("Message", quickreply.Action)

	content = actual.Content[1]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("QuickReply", content.GetType())

	quickreply, ok = content.(*gcloudcx.NormalizedMessageQuickReplyContent)
	suite.Require().True(ok, "Content should be of type OpenMessageQuickReplyContent, but was %T", content)
	suite.Require().NotNil(quickreply, "QuickReply should not be nil")
	suite.Assert().Equal("No", quickreply.Text)
	suite.Assert().Equal("No", quickreply.Payload)
	suite.Assert().Equal("Message", quickreply.Action)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageStructuredWithDatePicker() {
	payload := suite.LoadTestData("openmessaging-inbound-datepicker.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageStructured)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageStructured, but was %T", message)
	suite.Assert().Equal("c327c2078ca056db130c55ce648d9fa2", actual.ID)
	suite.Assert().Equal(uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"), actual.ConversationID)
	suite.Require().Len(actual.Content, 1, "Content should 1 item")

	content := actual.Content[0]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("DatePicker", content.GetType())

	datepicker, ok := content.(*gcloudcx.NormalizedMessageDatePickerContent)
	suite.Require().True(ok, "Content should be of type OpenMessageDatePickerContent, but was %T", content)
	suite.Require().NotNil(datepicker, "DatePicker should not be nil")
	suite.Assert().Equal("When would you be available?", datepicker.Title)
	suite.Assert().Equal("Pick a date and time", datepicker.Subtitle)
	suite.Require().Len(datepicker.AvailableTimes, 2, "DatePicker should have 2 available times")
	suite.Assert().Equal("2025-05-30T12:00:00Z", datepicker.AvailableTimes[0].Time.Format(time.RFC3339))
	suite.Assert().Equal(1800*time.Second, datepicker.AvailableTimes[0].Duration)
	suite.Assert().Equal("2025-06-30T13:00:00Z", datepicker.AvailableTimes[1].Time.Format(time.RFC3339))
	suite.Assert().Equal(900*time.Second, datepicker.AvailableTimes[1].Duration)
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

func (suite *OpenMessagingSuite) TestCanRedactOpenMessage() {
	output := CaptureStdout(func() {
		log := logger.Create("test", &logger.StdoutStream{})
		defer log.Flush()
		log.SetFilterLevel(logger.TRACE)
		message := gcloudcx.OpenMessageText{
			ID: "12345678",
			Channel: gcloudcx.OpenMessageChannel{
				Platform:  "Open",
				Type:      "Private",
				MessageID: "gmAy9zNkhf4ermFvHH9mB5",
				To:        &gcloudcx.OpenMessageTo{ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
				From: &gcloudcx.OpenMessageFrom{
					ID:        "abcdef12345",
					Type:      "Email",
					Firstname: "Bob",
					Lastname:  "Minion",
					Nickname:  "Bobby",
				},
				Time: time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC),
			},
			Direction:      "inbound",
			Text:           "text message",
			ConversationID: uuid.MustParse("d06cb41e-f938-4dcf-b823-c8af1a39d7e5"),
		}

		suite.Logger.Record("message", message).Infof("message")
		log.Record("message", message).Infof("message")
	})
	suite.Require().NotEmpty(output, "There was no output")
	lines := strings.Split(output, "\n")
	lines = lines[0 : len(lines)-1] // remove the last empty line
	suite.Require().Len(lines, 1, "There should be 1 line in the log output, found %d", len(lines))
	suite.LogLineEqual(lines[0], map[string]string{
		"message":  `map\[channel:map\[from:map\[firstName:REDACTED-[0-9a-f]+ id:abcdef12345 idType:Email lastName:REDACTED-[0-9a-f]+ nickname:REDACTED-[0-9a-f]+\] messageId:gmAy9zNkhf4ermFvHH9mB5 platform:Open time:2021-04-09T04:43:33Z to:map\[id:edce4efa-4abf-468b-ada7-cd6d35e7bbaf\] type:Private\] conversationId:d06cb41e-f938-4dcf-b823-c8af1a39d7e5 direction:inbound id:12345678 text:REDACTED-[0-9a-f]+ type:Text\]`,
		"hostname": `[a-zA-Z_0-9\-\.]+`,
		"level":    "30",
		"msg":      "message",
		"name":     "test",
		"pid":      "[0-9]+",
		"scope":    "main",
		"tid":      "[0-9]+",
		"time":     `[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+Z`,
		"topic":    "main",
		"v":        "0",
	})
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessageChannel() {
	output := CaptureStdout(func() {
		log := logger.Create("test", &logger.StdoutStream{})
		defer log.Flush()
		log.SetFilterLevel(logger.TRACE)
		channel := gcloudcx.OpenMessageChannel{
			Platform:  "Open",
			Type:      "Private",
			MessageID: "gmAy9zNkhf4ermFvHH9mB5",
			To:        &gcloudcx.OpenMessageTo{ID: "edce4efa-4abf-468b-ada7-cd6d35e7bbaf"},
			From: &gcloudcx.OpenMessageFrom{
				ID:        "abcdef12345",
				Type:      "Email",
				Firstname: "Bob",
				Lastname:  "Minion",
				Nickname:  "Bobby",
			},
			Time: time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC),
		}
		log.Record("channel", channel).Infof("channel")
		suite.Logger.Record("channel", channel).Infof("channel")
	})
	suite.Require().NotEmpty(output, "There was no output")
	lines := strings.Split(output, "\n")
	lines = lines[0 : len(lines)-1] // remove the last empty line
	suite.Require().Len(lines, 1, "There should be 1 line in the log output, found %d", len(lines))
	suite.LogLineEqual(lines[0], map[string]string{
		"channel":  `map\[from:map\[firstName:REDACTED-[0-9a-f]+ id:abcdef12345 idType:Email lastName:REDACTED-[0-9a-f]+ nickname:REDACTED-[0-9a-f]+\] messageId:gmAy9zNkhf4ermFvHH9mB5 platform:Open time:2021-04-09T04:43:33Z to:map\[id:edce4efa-4abf-468b-ada7-cd6d35e7bbaf\] type:Private\]`,
		"hostname": `[a-zA-Z_0-9\-\.]+`,
		"level":    "30",
		"msg":      "channel",
		"name":     "test",
		"pid":      "[0-9]+",
		"scope":    "main",
		"tid":      "[0-9]+",
		"time":     `[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+Z`,
		"topic":    "main",
		"v":        "0",
	})
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessageFrom() {
	output := CaptureStdout(func() {
		log := logger.Create("test", &logger.StdoutStream{})
		defer log.Flush()
		log.SetFilterLevel(logger.TRACE)
		from := gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		}
		log.Record("from", from).Infof("from")
		suite.Logger.Record("from", from).Infof("from")
	})
	suite.Require().NotEmpty(output, "There was no output")
	lines := strings.Split(output, "\n")
	lines = lines[0 : len(lines)-1] // remove the last empty line
	suite.Require().Len(lines, 1, "There should be 1 line in the log output, found %d", len(lines))
	suite.LogLineEqual(lines[0], map[string]string{
		"from":     `map\[firstName:REDACTED-[0-9a-f]+ id:abcdef12345 idType:Email lastName:REDACTED-[0-9a-f]+ nickname:REDACTED-[0-9a-f]+\]`,
		"hostname": `[a-zA-Z_0-9\-\.]+`,
		"level":    "30",
		"msg":      "from",
		"name":     "test",
		"pid":      "[0-9]+",
		"scope":    "main",
		"tid":      "[0-9]+",
		"time":     `[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+Z`,
		"topic":    "main",
		"v":        "0",
	})
}

func (suite *OpenMessagingSuite) TestCanRedactOpenMessageMetadata() {
	output := CaptureStdout(func() {
		log := logger.Create("test", &logger.StdoutStream{})
		defer log.Flush()
		log.SetFilterLevel(logger.TRACE)
		// TODO: call SendInboundMessage. We should pass a list of keys to redact in the metadata, somehow. Context?
		from := gcloudcx.OpenMessageFrom{
			ID:        "abcdef12345",
			Type:      "Email",
			Firstname: "Bob",
			Lastname:  "Minion",
			Nickname:  "Bobby",
		}
		log.Record("from", from).Infof("from")
		suite.Logger.Record("from", from).Infof("from")
	})
	suite.Require().NotEmpty(output, "There was no output")
	lines := strings.Split(output, "\n")
	lines = lines[0 : len(lines)-1] // remove the last empty line
	suite.Require().Len(lines, 1, "There should be 1 line in the log output, found %d", len(lines))
	suite.LogLineEqual(lines[0], map[string]string{
		"from":     `map\[firstName:REDACTED-[0-9a-f]+ id:abcdef12345 idType:Email lastName:REDACTED-[0-9a-f]+ nickname:REDACTED-[0-9a-f]+\]`,
		"hostname": `[a-zA-Z_0-9\-\.]+`,
		"level":    "30",
		"msg":      "from",
		"name":     "test",
		"pid":      "[0-9]+",
		"scope":    "main",
		"tid":      "[0-9]+",
		"time":     `[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+Z`,
		"topic":    "main",
		"v":        "0",
	})
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalIntegrationWithInvalidJSON() {
	var err error

	integration := gcloudcx.OpenMessagingIntegration{}
	err = json.Unmarshal([]byte(`{"Name": 15}`), &integration)
	suite.Assert().Error(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalChannelWithInvalidJSON() {
	var err error

	channel := gcloudcx.OpenMessageChannel{}
	err = json.Unmarshal([]byte(`{"Platform": 2}`), &channel)
	suite.Assert().Error(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalFromWithInvalidJSON() {
	var err error

	from := gcloudcx.OpenMessageFrom{}
	err = json.Unmarshal([]byte(`{"idType": 3}`), &from)
	suite.Assert().Error(err, "Data should not have been unmarshaled successfully")
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalMessageWithInvalidJSON() {
	_, err := gcloudcx.UnmarshalOpenMessage([]byte(`{"Direction": 6}`))
	suite.Assert().Error(err, "Data should not have been unmarshaled successfully")
}

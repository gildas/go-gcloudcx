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
	suite.IntegrationName = "TEST-GO-PURECLOUD"
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
		suite.Require().NoError(err, "Failed to login")
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
	integration, err := suite.Client.CreateOpenMessagingIntegration(context.Background(), suite.IntegrationName, webhookURL, webhookToken, nil)
	suite.Require().NoError(err, "Failed to create integration")
	suite.Logger.Record("integration", integration).Infof("Created a integration")
	for {
		if integration.IsCreated() {
			break
		}
		suite.Logger.Warnf("Integration %s is still in status: %s, waiting a bit", integration.ID, integration.CreateStatus)
		time.Sleep(time.Second)
		err = integration.Refresh(context.Background())
		suite.Require().NoError(err, "Failed to refresh integration")
	}
	suite.IntegrationID = integration.ID
}

func (suite *OpenMessagingSuite) TestCanFetchByID() {
	integration, err := gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().NoErrorf(err, "Failed to fetch Open Messaging Integration %s. %s", suite.IntegrationID, err)
	suite.Assert().Equal(suite.IntegrationID, integration.ID)
	suite.Assert().Equal(suite.IntegrationName, integration.Name)
}

func (suite *OpenMessagingSuite) TestCanFetchByName() {
	match := func(integration gcloudcx.OpenMessagingIntegration) bool {
		return integration.Name == suite.IntegrationName
	}
	integration, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	suite.Require().NoErrorf(err, "Failed to fetch Open Messaging Integration %s. %s", suite.IntegrationName, err)
	suite.Assert().Equal(suite.IntegrationID, integration.ID)
	suite.Assert().Equal(suite.IntegrationName, integration.Name)
}

func (suite *OpenMessagingSuite) TestCanFetchIntegrations() {
	integrations, err := gcloudcx.FetchAll[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client)
	suite.Require().NoError(err, "Failed to fetch OpenMessaging Integrations")
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
	integration, err := gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().NoErrorf(err, "Failed to fetch integration %s, Error: %s", suite.IntegrationID, err)
	suite.Logger.Record("integration", integration).Infof("Got a integration")
	suite.Require().True(integration.IsCreated(), "Integration should be created")
	err = integration.Delete(context.Background())
	suite.Require().NoErrorf(err, "Failed to delete integration %s, Error: %s", suite.IntegrationID, err)
	_, err = gcloudcx.Fetch[gcloudcx.OpenMessagingIntegration](context.Background(), suite.Client, suite.IntegrationID)
	suite.Require().Error(err, "Integration should not exist anymore")
	suite.Assert().ErrorIsf(err, gcloudcx.NotFoundError, "Expected NotFoundError, got %s", err)
	suite.Assert().Truef(errors.Is(err, gcloudcx.NotFoundError), "Expected NotFoundError, got %s", err)
	details := gcloudcx.NotFoundError.Clone()
	suite.Require().ErrorAsf(err, &details, "Expected NotFoundError but got %s", err)
	suite.IntegrationID = uuid.Nil
}

func (suite *OpenMessagingSuite) TestCanUnmarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{}
	err := suite.UnmarshalData("openmessagingintegration.json", &integration)
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
}

func (suite *OpenMessagingSuite) TestCanMarshalIntegration() {
	integration := gcloudcx.OpenMessagingIntegration{
		ID:           uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"),
		Name:         "TEST-GO-PURECLOUD",
		WebhookURL:   core.Must(url.Parse("https://www.acme.com/gcloudcx")),
		WebhookToken: "DEADBEEF",
		SupportedContent: &gcloudcx.AddressableEntityRef{
			ID:      uuid.MustParse("832066dd-6030-46b1-baeb-b89b681c6636"),
			SelfURI: "/api/v2/conversations/messaging/supportedcontent/832066dd-6030-46b1-baeb-b89b681c6636",
		},
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
		CreateStatus: "Initiated",
	}

	data, err := json.Marshal(integration)
	suite.Require().NoErrorf(err, "Failed to marshal OpenMessagingIntegration. %s", err)
	expected := suite.LoadTestData("openmessagingintegration.json")
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *OpenMessagingSuite) TestShouldNotUnmarshalIntegrationWithInvalidJSON() {
	var err error

	integration := gcloudcx.OpenMessagingIntegration{}
	err = json.Unmarshal([]byte(`{"Name": 15}`), &integration)
	suite.Assert().Error(err, "Data should not have been unmarshaled successfully")
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

func (suite *OpenMessagingSuite) TestCanUnmarshalOpenMessageStructuredWithParameters() {
	payload := suite.LoadTestData("outbound-text-with-parameters.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageStructured)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageStructured, but was %T", message)
	suite.Assert().Equal("68d79558191e6f93ee7e3c1f996994cd", actual.ID)
	suite.Assert().Equal("Hi Happy, How can I help you?", actual.Text)
	suite.Require().NotEmpty(actual.Content, "Content should not be empty")

	content := actual.Content[0]
	suite.Require().NotNil(content, "Content should not be nil")
	suite.Require().Equal("Notification", content.Type)
	suite.Require().NotNil(content.Template, "Content Template should not be nil")
	suite.Assert().Equal("Hi Happy, How can I help you?", content.Template.Text)
	suite.Require().NotEmpty(content.Template.Parameters, "Content Template Parameters should not be empty")
	value, found := content.Template.Parameters["name"]
	suite.Require().True(found, "Content Template Parameters should contain 'name'")
	suite.Assert().Equal("Happy", value)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalReceipt() {
	var receipt gcloudcx.OpenMessageReceipt
	err := suite.UnmarshalData("response-inbound-receipt.json", &receipt)
	suite.Require().NoErrorf(err, "Failed to unmarshal OpenMessageReceipt. %s", err)
	suite.Assert().False(receipt.IsFailed(), "Receipt should not be failed")
	suite.Assert().Equal("Delivered", receipt.Status)
	suite.Assert().Equal("Outbound", receipt.Direction)
	suite.Assert().Equal("c8e53e498891dfc9400c79a278cc1863", receipt.ID)
}

func (suite *OpenMessagingSuite) TestCanUnmarshalReceiptWithErrors() {
	var receipt gcloudcx.OpenMessageReceipt
	err := suite.UnmarshalData("response-inbound-receipt-errors.json", &receipt)
	suite.Require().NoErrorf(err, "Failed to unmarshal OpenMessageReceipt. %s", err)
	suite.Assert().True(receipt.IsFailed(), "Receipt should be failed")
	suite.Assert().Equal("Outbound", receipt.Direction)
	suite.Assert().Equal("c8e53e498891dfc9400c79a278cc1863", receipt.ID)
	err = receipt.AsError()
	suite.Require().Error(err, "Receipt should be convert to an error")
	suite.ErrorIs(err, gcloudcx.GeneralError, "Receipt should convert to a GeneralError")
	suite.ErrorIs(err, gcloudcx.RateLimited, "Receipt should convert to a RateLimited")
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
			Direction: "inbound",
			Text:      "text message",
		}

		suite.Logger.Record("message", message).Infof("message")
		log.Record("message", message).Infof("message")
	})
	suite.Require().NotEmpty(output, "There was no output")
	lines := strings.Split(output, "\n")
	lines = lines[0 : len(lines)-1] // remove the last empty line
	suite.Require().Len(lines, 1, "There should be 1 line in the log output, found %d", len(lines))
	suite.LogLineEqual(lines[0], map[string]string{
		"message":  `map\[channel:map\[from:map\[firstName:REDACTED-[0-9a-f]+ id:abcdef12345 idType:Email lastName:REDACTED-[0-9a-f]+ nickname:REDACTED-[0-9a-f]+\] messageId:gmAy9zNkhf4ermFvHH9mB5 platform:Open time:2021-04-09T04:43:33Z to:map\[id:edce4efa-4abf-468b-ada7-cd6d35e7bbaf\] type:Private\] direction:inbound id:12345678 text:REDACTED-[0-9a-f]+ type:Text\]`,
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

func (suite *OpenMessagingSuite) TestCanMarshalTypingEvent() {
	event := gcloudcx.OpenMessageEvents{
		Channel: gcloudcx.OpenMessageChannel{
			Platform: "Open",
			Type:     "Private",
			From: &gcloudcx.OpenMessageFrom{
				ID:        "abcdef12345",
				Type:      "Email",
				Firstname: "Bob",
				Lastname:  "Minion",
				Nickname:  "Bobby",
			},
			Time: time.Date(2021, 4, 9, 4, 43, 33, 0, time.UTC),
		},
		Events: []gcloudcx.OpenMessageEvent{
			gcloudcx.OpenMessageTypingEvent{IsTyping: true},
		},
	}
	payload, err := json.Marshal(event)
	suite.Require().NoErrorf(err, "Failed to marshal OpenMessageEvents. %s", err)
	expected := suite.LoadTestData("openmessaging-event-typing.json")
	suite.Require().JSONEq(string(expected), string(payload))
}

func (suite *OpenMessagingSuite) TestCanUnmarshalTypingEvent() {
	payload := suite.LoadTestData("inbound-openmessaging-event-typing.json")
	message, err := gcloudcx.UnmarshalOpenMessage(payload)
	suite.Require().NoError(err, "Failed to unmarshal OpenMessage")
	suite.Require().NotNil(message, "Unmarshaled message should not be nil")

	actual, ok := message.(*gcloudcx.OpenMessageEvents)
	suite.Require().True(ok, "Unmarshaled message should be of type OpenMessageEvents, but was %T", message)
	suite.Require().NotNil(actual, "Unmarshaled message should not be nil")
	suite.Assert().Equal("6ffd815bca1570e46251fcc71c103837", actual.ID)
	suite.Assert().Equal(uuid.MustParse("1af69355-f1b0-477e-8ed9-66baff370209"), actual.Channel.ID)
	suite.Assert().Equal("Outbound", actual.Direction)
	suite.Require().Len(actual.Events, 1, "Unmarshaled message should have 1 event")

	messageEvent, ok := actual.Events[0].(*gcloudcx.OpenMessageTypingEvent)
	suite.Require().True(ok, "Unmarshaled message event should be of type *OpenMessageTypingEvent, but was %T", actual.Events[0])
	suite.Assert().True(messageEvent.IsTyping)
	suite.Assert().Equal(5*time.Second, messageEvent.Duration)
}

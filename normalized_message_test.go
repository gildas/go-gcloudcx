package gcloudcx_test

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type NormalizedMessageSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	IntegrationID   uuid.UUID
	IntegrationName string
	Client          *gcloudcx.Client
}

func TestNormalizedMessageSuite(t *testing.T) {
	suite.Run(t, new(NormalizedMessageSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *NormalizedMessageSuite) SetupSuite() {
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

func (suite *NormalizedMessageSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *NormalizedMessageSuite) BeforeTest(suiteName, testName string) {
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

func (suite *NormalizedMessageSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *NormalizedMessageSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *NormalizedMessageSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

func (suite *NormalizedMessageSuite) LogLineEqual(line string, records map[string]string) {
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

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *NormalizedMessageSuite) TestCanMarshalCarousel() {
	message := gcloudcx.NormalizedMessage{
		Type: gcloudcx.NormalizedMessageTypeStructured,
		Content: []gcloudcx.NormalizedMessageContent{
			gcloudcx.NormalizedMessageCarouselContent{
				Cards: []gcloudcx.NormalizedMessageCarouselCard{
					{
						Title:       "Card 1",
						Description: "Description 1",
						ImageURL:    core.Must(url.Parse("https://www.acme.com/image1.png")),
						Actions: []gcloudcx.NormalizedMessageCardAction{
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option1", Payload: "Option1"},
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option2", Payload: "Option2"},
							gcloudcx.NormalizedMessageCardLinkAction{Text: "Option3", URL: core.Must(url.Parse("https://www.acme.com/option3"))},
						},
					},
					{
						Title:       "Card 2",
						Description: "Description 2",
						ImageURL:    core.Must(url.Parse("https://www.acme.com/image2.png")),
						Actions: []gcloudcx.NormalizedMessageCardAction{
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option4", Payload: "Option4"},
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option5", Payload: "Option5"},
							gcloudcx.NormalizedMessageCardLinkAction{Text: "Option6", URL: core.Must(url.Parse("https://www.acme.com/option6"))},
						},
					},
					{
						Title:       "Card 3",
						Description: "Description 3",
						VideoURL:    core.Must(url.Parse("https://www.acme.com/video3.mp4")),
						Actions: []gcloudcx.NormalizedMessageCardAction{
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option7", Payload: "Option7"},
							gcloudcx.NormalizedMessageCardPostbackAction{Text: "Option8", Payload: "Option8"},
							gcloudcx.NormalizedMessageCardLinkAction{Text: "Option9", URL: core.Must(url.Parse("https://www.acme.com/option9"))},
						},
					},
				},
			},
		},
	}
	payload, err := json.Marshal(message)
	suite.Require().NoErrorf(err, "Failed to marshal NormalizedMessage. %s", err)
	expected := suite.LoadTestData("normalized-message-structured-carousel.json")
	suite.Assert().JSONEq(string(expected), string(payload))
}

func (suite *NormalizedMessageSuite) TestCanUnmarshalCarousel() {
	payload := suite.LoadTestData("normalized-message-structured-carousel.json")
	var message gcloudcx.NormalizedMessage
	err := json.Unmarshal(payload, &message)
	suite.Require().NoErrorf(err, "Failed to unmarshal NormalizedMessage. %s", err)
	suite.Assert().Equal(gcloudcx.NormalizedMessageTypeStructured, message.Type)
	suite.Assert().Len(message.Content, 1)

	carousel, ok := message.Content[0].(*gcloudcx.NormalizedMessageCarouselContent)
	suite.Assert().True(ok, "The content is not a *NormalizedMessageCarouselContent but a %T", message.Content[0])
	suite.Assert().Len(carousel.Cards, 3)
	suite.Assert().Equal("Carousel", carousel.GetType())

	suite.Assert().Equal("Card 1", carousel.Cards[0].Title)
	suite.Assert().Equal("Description 1", carousel.Cards[0].Description)
	suite.Require().NotNil(carousel.Cards[0].ImageURL)
	suite.Assert().Equal("https://www.acme.com/image1.png", carousel.Cards[0].ImageURL.String())
	suite.Assert().Len(carousel.Cards[0].Actions, 3)

	postbackAction, ok := carousel.Cards[0].Actions[0].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[0].Actions[0])
	suite.Assert().Equal("Option1", postbackAction.Text)
	suite.Assert().Equal("Option1", postbackAction.Payload)

	postbackAction, ok = carousel.Cards[0].Actions[1].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[0].Actions[1])
	suite.Assert().Equal("Option2", postbackAction.Text)
	suite.Assert().Equal("Option2", postbackAction.Payload)

	linkAction, ok := carousel.Cards[0].Actions[2].(*gcloudcx.NormalizedMessageCardLinkAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardLinkAction but a %T", carousel.Cards[0].Actions[2])
	suite.Assert().Equal("Option3", linkAction.Text)
	suite.Require().NotNil(linkAction.URL, "The link action URL should not be nil")
	suite.Assert().Equal("https://www.acme.com/option3", linkAction.URL.String())

	suite.Assert().Equal("Card 2", carousel.Cards[1].Title)
	suite.Assert().Equal("Description 2", carousel.Cards[1].Description)
	suite.Require().NotNil(carousel.Cards[1].ImageURL)
	suite.Assert().Equal("https://www.acme.com/image2.png", carousel.Cards[1].ImageURL.String())
	suite.Assert().Len(carousel.Cards[1].Actions, 3)

	postbackAction, ok = carousel.Cards[1].Actions[0].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[1].Actions[0])
	suite.Assert().Equal("Option4", postbackAction.Text)
	suite.Assert().Equal("Option4", postbackAction.Payload)

	postbackAction, ok = carousel.Cards[1].Actions[1].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[1].Actions[1])
	suite.Assert().Equal("Option5", postbackAction.Text)
	suite.Assert().Equal("Option5", postbackAction.Payload)

	linkAction, ok = carousel.Cards[1].Actions[2].(*gcloudcx.NormalizedMessageCardLinkAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardLinkAction but a %T", carousel.Cards[1].Actions[2])
	suite.Assert().Equal("Option6", linkAction.Text)
	suite.Require().NotNil(linkAction.URL, "The link action URL should not be nil")
	suite.Assert().Equal("https://www.acme.com/option6", linkAction.URL.String())

	suite.Assert().Equal("Card 3", carousel.Cards[2].Title)
	suite.Assert().Equal("Description 3", carousel.Cards[2].Description)
	suite.Require().NotNil(carousel.Cards[2].VideoURL)
	suite.Assert().Equal("https://www.acme.com/video3.mp4", carousel.Cards[2].VideoURL.String())
	suite.Assert().Len(carousel.Cards[2].Actions, 3)

	postbackAction, ok = carousel.Cards[2].Actions[0].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[2].Actions[0])
	suite.Assert().Equal("Option7", postbackAction.Text)
	suite.Assert().Equal("Option7", postbackAction.Payload)

	postbackAction, ok = carousel.Cards[2].Actions[1].(*gcloudcx.NormalizedMessageCardPostbackAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardPostbackAction but a %T", carousel.Cards[2].Actions[1])
	suite.Assert().Equal("Option8", postbackAction.Text)
	suite.Assert().Equal("Option8", postbackAction.Payload)

	linkAction, ok = carousel.Cards[2].Actions[2].(*gcloudcx.NormalizedMessageCardLinkAction)
	suite.Assert().True(ok, "The action is not a *NormalizedMessageCardLinkAction but a %T", carousel.Cards[2].Actions[2])
	suite.Assert().Equal("Option9", linkAction.Text)
	suite.Require().NotNil(linkAction.URL, "The link action URL should not be nil")
	suite.Assert().Equal("https://www.acme.com/option9", linkAction.URL.String())
}

func (suite *NormalizedMessageSuite) TestCanMarshalQuickReplies() {
	message := gcloudcx.NormalizedMessage{
		Type: gcloudcx.NormalizedMessageTypeStructured,
		Text: "Do you want to proceed?",
		Content: []gcloudcx.NormalizedMessageContent{
			gcloudcx.NormalizedMessageQuickReplyContent{Text: "Yes", Payload: "Yes", Action: "Message"},
			gcloudcx.NormalizedMessageQuickReplyContent{Text: "No", Payload: "No", Action: "Message"},
		},
	}
	payload, err := json.Marshal(message)
	suite.Require().NoErrorf(err, "Failed to marshal NormalizedMessage. %s", err)
	expected := suite.LoadTestData("normalized-message-structured-quick-replies.json")
	suite.Assert().JSONEq(string(expected), string(payload))
}

func (suite *NormalizedMessageSuite) TestCanUnmarshalQuickReplies() {
	payload := suite.LoadTestData("normalized-message-structured-quick-replies.json")
	var message gcloudcx.NormalizedMessage
	err := json.Unmarshal(payload, &message)
	suite.Require().NoErrorf(err, "Failed to unmarshal NormalizedMessage. %s", err)
	suite.Assert().Equal(gcloudcx.NormalizedMessageTypeStructured, message.Type)
	suite.Assert().Equal("Do you want to proceed?", message.Text)
	suite.Assert().Len(message.Content, 2)

	quickReply, ok := message.Content[0].(*gcloudcx.NormalizedMessageQuickReplyContent)
	suite.Assert().True(ok, "The content is not a *NormalizedMessageQuickReplyContent but a %T", message.Content[0])
	suite.Assert().Equal("Yes", quickReply.Text)
	suite.Assert().Equal("Yes", quickReply.Payload)
	suite.Assert().Equal("Message", quickReply.Action)

	quickReply, ok = message.Content[1].(*gcloudcx.NormalizedMessageQuickReplyContent)
	suite.Assert().True(ok, "The content is not a *NormalizedMessageQuickReplyContent but a %T", message.Content[1])
	suite.Assert().Equal("No", quickReply.Text)
	suite.Assert().Equal("No", quickReply.Payload)
	suite.Assert().Equal("Message", quickReply.Action)
}

func (suite *NormalizedMessageSuite) TestCanMarshalDatePicker() {
	message := gcloudcx.NormalizedMessage{
		Type: gcloudcx.NormalizedMessageTypeStructured,
		Content: []gcloudcx.NormalizedMessageContent{
			gcloudcx.NormalizedMessageDatePickerContent{
				Title:    "When would you be available?",
				Subtitle: "Pick a date and time",
				AvailableTimes: []gcloudcx.NormalizedAvailableTime{
					{Time: time.Date(2025, 5, 30, 12, 0, 0, 0, time.UTC), Duration: 30 * time.Minute},
					{Time: time.Date(2025, 6, 30, 13, 0, 0, 0, time.UTC), Duration: 15 * time.Minute},
				},
			},
		},
	}
	payload, err := json.Marshal(message)
	suite.Require().NoErrorf(err, "Failed to marshal NormalizedMessage. %s", err)
	expected := suite.LoadTestData("normalized-message-structured-datepicker.json")
	suite.Assert().JSONEq(string(expected), string(payload))
}

func (suite *NormalizedMessageSuite) TestCanUnmarshalDatePicker() {
	payload := suite.LoadTestData("normalized-message-structured-datepicker.json")
	var message gcloudcx.NormalizedMessage
	err := json.Unmarshal(payload, &message)
	suite.Require().NoErrorf(err, "Failed to unmarshal NormalizedMessage. %s", err)
	suite.Assert().Equal(gcloudcx.NormalizedMessageTypeStructured, message.Type)
	suite.Assert().Len(message.Content, 1)

	datePicker, ok := message.Content[0].(*gcloudcx.NormalizedMessageDatePickerContent)
	suite.Assert().True(ok, "The content is not a *NormalizedMessageDatePickerContent but a %T", message.Content[0])
	suite.Assert().Equal("When would you be available?", datePicker.Title)
	suite.Assert().Equal("Pick a date and time", datePicker.Subtitle)
	suite.Assert().Len(datePicker.AvailableTimes, 2)

	suite.Assert().Equal(time.Date(2025, 5, 30, 12, 0, 0, 0, time.UTC), datePicker.AvailableTimes[0].Time)
	suite.Assert().Equal(30*time.Minute, datePicker.AvailableTimes[0].Duration)
	suite.Assert().Equal(time.Date(2025, 6, 30, 13, 0, 0, 0, time.UTC), datePicker.AvailableTimes[1].Time)
	suite.Assert().Equal(15*time.Minute, datePicker.AvailableTimes[1].Duration)
}

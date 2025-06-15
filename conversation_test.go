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
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type ConversationSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	ConversationID uuid.UUID
	Client         *gcloudcx.Client
}

func TestConversationSuite(t *testing.T) {
	suite.Run(t, new(ConversationSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *ConversationSuite) SetupSuite() {
	var err error
	var value string

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
		region   = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret   = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	)

	value = core.GetEnvAsString("CONVERSATION_ID", "")
	suite.Require().NotEmpty(value, "CONVERSATION_ID is not set in your environment")

	suite.ConversationID, err = uuid.Parse(value)
	suite.Require().NoError(err, "CONVERSATION_ID is not a valid UUID")

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: region,
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *ConversationSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *ConversationSuite) BeforeTest(suiteName, testName string) {
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

func (suite *ConversationSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *ConversationSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *ConversationSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *ConversationSuite) TestCanUnmarshal() {
	conversation := gcloudcx.Conversation{}
	err := suite.UnmarshalData("conversation.json", &conversation)
	suite.Require().NoErrorf(err, "Failed to unmarshal conversation. %s", err)
	suite.Logger.Record("Conversation", conversation).Infof("Got a conversation")
	suite.Assert().NotEmpty(conversation.ID)
}

func (suite *ConversationSuite) TestCanFetchByID() {
	suite.T().Skip("We need a reliable way to fetch a conversation by ID forever")
	conversation, correlationID, err := gcloudcx.Fetch[gcloudcx.Conversation](context.Background(), suite.Client, suite.ConversationID)
	suite.Require().NoErrorf(err, "Failed to fetch Conversation %s. %s", suite.ConversationID, err)
	suite.Assert().Equal(suite.ConversationID, conversation.ID)
	suite.Logger.Infof("Correlation: %s", correlationID)
}

func (suite *ConversationSuite) TestCanUnmarshalRecordings() {
	recordings := []gcloudcx.Recording{}
	err := suite.UnmarshalData("recordings.json", &recordings)
	suite.Require().NoErrorf(err, "Failed to unmarshal recordings. %s", err)
	suite.Logger.Record("Recordings", recordings).Infof("Got recordings")
	suite.Assert().Len(recordings, 2)
}

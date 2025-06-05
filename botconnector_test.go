package gcloudcx_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type BotConnectorSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestBotConnectorSuite(t *testing.T) {
	suite.Run(t, new(BotConnectorSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *BotConnectorSuite) SetupSuite() {
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
}

func (suite *BotConnectorSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *BotConnectorSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *BotConnectorSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *BotConnectorSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *BotConnectorSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *BotConnectorSuite) TestCanUnmarshalIncomingMessageRequest() {
	var request gcloudcx.BotConnectorIncomingMessageRequest
	err := suite.UnmarshalData("botconnector-incoming-message-request.json", &request)
	suite.Require().NoError(err, "Failed to unmarshal BotConnectorIncomingMessageRequest")
	suite.Logger.Record("IncomingMessageRequest", request).Infof("Got IncomingMessageRequest")
}

func (suite *BotConnectorSuite) TestCanMarshalIncomingMessageResponse() {
	var response gcloudcx.BotConnectorIncomingMessageResponse
	err := suite.UnmarshalData("botconnector-incoming-message-response-complete.json", &response)
	suite.Require().NoError(err, "Failed to unmarshal BotConnectorIncomingMessageResponse")
	suite.Logger.Record("IncomingMessageResponse", response).Infof("Got IncomingMessageResponse")
}

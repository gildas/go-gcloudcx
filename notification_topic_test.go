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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type NotificationTopicSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestNotificationTopicSuite(t *testing.T) {
	suite.Run(t, new(NotificationTopicSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *NotificationTopicSuite) SetupSuite() {
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

func (suite *NotificationTopicSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *NotificationTopicSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))

	suite.Start = time.Now()
}

func (suite *NotificationTopicSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *NotificationTopicSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *NotificationTopicSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *NotificationTopicSuite) TestCanUnmarshalConverstionMessageTopic() {
	payload := suite.LoadTestData("notification_topic_chat_message.json")

	topic, err := gcloudcx.UnmarshalNotificationTopic(payload)
	suite.Require().NoErrorf(err, "Failed to Unmarshal Notification Topic. %s", err)
	suite.Require().NotNil(topic, "Unmarshal Notification Topic returned nil")

	actual, ok := topic.(gcloudcx.ConversationChatMessageTopic)
	suite.Require().Truef(ok, "Expected a ConversationChatMessageTopic, got %T", topic)
	suite.Require().NotNil(actual, "Cast Notification Topic returned nil")
}

func (suite *NotificationTopicSuite) TestCanUnmarshalConversationChatMemberTopic() {
	payload := suite.LoadTestData("notification_topic_chat_member.json")

	topic, err := gcloudcx.UnmarshalNotificationTopic(payload)
	suite.Require().NoErrorf(err, "Failed to Unmarshal Notification Topic. %s", err)
	suite.Require().NotNil(topic, "Unmarshal Notification Topic returned nil")

	actual, ok := topic.(gcloudcx.ConversationChatMemberTopic)
	suite.Require().Truef(ok, "Expected a ConversationChatMemberTopic, got %T", topic)
	suite.Require().NotNil(actual, "Cast Notification Topic returned nil")
}

func (suite *NotificationTopicSuite) TestCanInstantiateFromString() {
	expected := gcloudcx.ConversationChatMemberTopic{}
	topic, err := gcloudcx.NotificationTopicFrom("v2.conversations.chats.aa06a6fc-1fdf-4e59-b8a1-df3ca44f523e.members")
	suite.Require().NoErrorf(err, "Failed to instantiate Notification Topic. %s", err)
	suite.Require().NotNil(topic, "Instantiate Notification Topic returned nil")
	suite.Require().Equal(expected.GetType(), topic.GetType())
	suite.Require().Len(topic.GetTargets(), 1)
	suite.Require().Equal(uuid.MustParse("aa06a6fc-1fdf-4e59-b8a1-df3ca44f523e"), topic.GetTargets()[0].GetID())
}

func (suite *NotificationTopicSuite) TestCanGetTopicName() {
	expected := "v2.conversations.chats.aa06a6fc-1fdf-4e59-b8a1-df3ca44f523e.members"
	entityRef := gcloudcx.EntityRef{ID: uuid.MustParse("aa06a6fc-1fdf-4e59-b8a1-df3ca44f523e")}
	topic := gcloudcx.ConversationChatMemberTopic{}
	suite.Require().Equal(expected, topic.With(entityRef).String())
}

package gcloudcx_test

import (
	"context"
	"fmt"
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

type GroupSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	GroupID   uuid.UUID
	GroupName string
	Client    *gcloudcx.Client
}

func TestGroupSuite(t *testing.T) {
	suite.Run(t, new(GroupSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *GroupSuite) SetupSuite() {
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
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID     = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		deploymentID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""))
	)

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})

	value = core.GetEnvAsString("GROUP_ID", "")
	suite.Require().NotEmpty(value, "GROUP_ID is not set in your environment")

	suite.GroupID, err = uuid.Parse(value)
	suite.Require().NoError(err, "GROUP_ID is not a valid UUID")

	suite.GroupName = core.GetEnvAsString("GROUP_NAME", "")
	suite.Require().NotEmpty(suite.GroupName, "GROUP_NAME is not set in your environment")

	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *GroupSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *GroupSuite) BeforeTest(suiteName, testName string) {
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

func (suite *GroupSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *GroupSuite) TestCanFetchByID() {
	group, correlationID, err := gcloudcx.Fetch[gcloudcx.Group](context.Background(), suite.Client, suite.GroupID)
	suite.Require().NoErrorf(err, "Failed to fetch Group %s. %s", suite.GroupID, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Assert().Equal(suite.GroupID, group.ID)
	suite.Assert().Equal(suite.GroupName, group.Name)
	suite.Assert().Equal("public", group.Visibility)
}

func (suite *GroupSuite) TestCanFetchByName() {
	match := func(group gcloudcx.Group) bool {
		return group.Name == suite.GroupName
	}
	group, correlationID, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	suite.Require().NoErrorf(err, "Failed to fetch Group %s. %s", suite.GroupName, err)
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Assert().Equal(suite.GroupID, group.ID)
	suite.Assert().Equal(suite.GroupName, group.Name)
	suite.Assert().Equal("public", group.Visibility)
}

func (suite *GroupSuite) TestCanStringify() {
	id := uuid.New()
	group := gcloudcx.Group{
		ID:   id,
		Name: "Hello",
	}
	suite.Assert().Equal("Hello", group.String())
	group.Name = ""
	suite.Assert().Equal(id.String(), group.String())
}

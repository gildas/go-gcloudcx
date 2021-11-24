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

type UserSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	UserID   uuid.UUID
	UserName string
	Client   *gcloudcx.Client
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) TestCanUnmarshal() {
	user := gcloudcx.User{}
	err := LoadObject("user.json", &user)
	suite.Require().Nil(err, "Failed to unmarshal user. %s", err)
	suite.Logger.Record("User", user).Infof("Got a user")
	suite.Assert().NotEmpty(user.ID)
	suite.Assert().Equal("John Doe", user.Name)
}

func (suite *UserSuite) TestCanMarshal() {
	user := gcloudcx.User{
		ID:       uuid.MustParse("06ffcd2e-1ada-412e-a5f5-30d7853246dd"),
		Name:     "John Doe",
		UserName: "john.doe@acme.com",
		Mail:     "john.doe@acme.com",
		Title:    "Junior",
		Division: &gcloudcx.Division{
			ID:      uuid.MustParse("06ffcd2e-1ada-412e-a5f5-30d7853246dd"),
			Name:    "",
			SelfURI: "/api/v2/authorization/divisions/06ffcd2e-1ada-412e-a5f5-30d7853246dd",
		},
		Chat: &gcloudcx.Jabber{
			ID: "98765432d220541234567654@genesysapacanz.orgspan.com",
		},
		Addresses: []*gcloudcx.Contact{},
		PrimaryContact: []*gcloudcx.Contact{
			{
				Type:      "PRIMARY",
				MediaType: "EMAIL",
				Address:   "john.doe@acme.com",
			},
		},
		Images: []*gcloudcx.UserImage{
			{
				Resolution: "x96",
				ImageURL:   core.Must(url.Parse("https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/7fac0a12/4643/4d0e/86f3/2467894311b5.jpg")).(*url.URL),
			},
		},
		AcdAutoAnswer: false,
		State:         "active",
		SelfURI:       "/api/v2/users/06ffcd2e-1ada-412e-a5f5-30d7853246dd",
		Version:       29,
	}

	data, err := json.Marshal(user)
	suite.Require().Nil(err, "Failed to marshal User. %s", err)
	expected, err := LoadFile("user.json")
	suite.Require().Nil(err, "Failed to Load Data. %s", err)
	suite.Assert().JSONEq(string(expected), string(data))
}

func (suite *UserSuite) TestCanFetchByID() {
	user := gcloudcx.User{}
	err := suite.Client.Fetch(context.Background(), &user, suite.UserID)
	suite.Require().Nilf(err, "Failed to fetch User %s. %s", suite.UserID, err)
	suite.Assert().Equal(suite.UserName, user.Name)
}

// Suite Tools

func (suite *UserSuite) SetupSuite() {
	var err error
	var value string

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
		region   = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret   = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	)

	value = core.GetEnvAsString("USER_ID", "")
	suite.Require().NotEmpty(value, "USER_ID is not set in your environment")

	suite.UserID, err = uuid.Parse(value)
	suite.Require().Nil(err, "USER_ID is not a valid UUID")

	suite.UserName = core.GetEnvAsString("USER_NAME", "")
	suite.Require().NotEmpty(suite.UserName, "USER_NAME is not set in your environment")

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: region,
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *UserSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *UserSuite) BeforeTest(suiteName, testName string) {
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

func (suite *UserSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

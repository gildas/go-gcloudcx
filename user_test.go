package purecloud_test

import (
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
	"github.com/stretchr/testify/suite"

	purecloud "github.com/gildas/go-purecloud"
)

type UserSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *purecloud.Client
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(OrganizationSuite))
}

func (suite *UserSuite) TestCanUnmarshal() {
	source := `{"id":"06f3052e-1ddc-41bd-a7d5-30d44682dfdd","name":"Matt McPhee","division":{"id":"8676cfc3-c94c-49fc-85a0-8c70a5943d8e","name":"","selfUri":"/api/v2/authorization/divisions/8676cfc3-c94c-49fc-85a0-8c70a5943d8e"},"chat":{"jabberId":"59a49de3d220541cb2589365@genesysapacanz.orgspan.com"},"email":"mcphee11@gmail.com","primaryContactInfo":[{"address":"mcphee11@gmail.com","mediaType":"EMAIL","type":"PRIMARY"}],"addresses":[],"state":"active","title":"Junior","username":"mcphee11@gmail.com","images":[{"resolution":"x96","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/7fac0a12/4643/4d0e/86f3/24e3fd9311b5.jpg"},{"resolution":"x128","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/12e66044/828b/40e7/93d0/bd76358a1d5d.jpg"},{"resolution":"x200","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/2b169f90/4bc2/4712/a866/b3351739f4b2.jpg"},{"resolution":"x48","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/c9199e42/9161/4dc4/a6ef/2c367bb582c2.jpg"},{"resolution":"x300","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazonaws.com/eef83617/e7b9/4d51/b54a/e5ca9fdb2713.jpg"},{"resolution":"x400","imageUri":"https://prod-apse2-inin-directory-service-profile.s3-ap-southeast-2.amazon aws.com/dce99686/d9f2/4fed/8990/f2edf782bc63.jpg"}],"version":29,"acdAutoAnswer":false,"selfUri":"/api/v2/users/06f3052e-1ddc-41bd-a7d5-30d44682dfdd"}` 

	user := purecloud.User{}
	err := json.Unmarshal([]byte(source), &user)
	suite.Require().Nil(err, "Failed to unmarshal user. %s", err)
	log.Record("User", user).Infof("Got a user")
	suite.Assert().NotEmpty(user.ID)
	suite.Assert().Equal("Matt McPhee", user.Name)
}

// Suite Tools

func (suite *UserSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	logFolder := filepath.Join(".", "log")
	err := os.MkdirAll(logFolder, os.ModePerm)
	suite.Require().Nil(err, "Failed to create folder: %s", logFolder)
	suite.Logger = logger.CreateWithDestination("test", fmt.Sprintf("file://%s/test-%s.log", logFolder, strings.ToLower(suite.Name)))
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region   = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
		secret   = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	)

	suite.Client = purecloud.New(purecloud.ClientOptions{
		Region:       region,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "PureCloud Client is nil")
}

func (suite *UserSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *UserSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()

	// Reuse tokens as much as we can
	if !suite.Client.IsAuthorized() {
		suite.Logger.Infof("Client is not logged in...")
		err := suite.Client.Login()
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
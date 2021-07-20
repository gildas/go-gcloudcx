// +build integration

package gcloudcx_test

import (
	"fmt"
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

type OrganizationSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *gcloudcx.Client
}

func TestOrganizationSuite(t *testing.T) {
	suite.Run(t, new(OrganizationSuite))
}

func (suite *OrganizationSuite) TestCanFetchOrganization() {
	organization := &gcloudcx.Organization{}
	err := suite.Client.Fetch(organization)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().Nil(err, "Failed to fetch my Organization")
	suite.Require().NotNil(organization, "Client's Organization is not loaded")
	suite.Assert().NotEmpty(organization.ID, "Client's Orgization has no ID")
	suite.Assert().NotEmpty(organization.Name, "Client's Orgization has no Name")
	suite.Assert().NotEmpty(organization.Features, "Client's Orgization has no features")
	suite.T().Logf("Organization: %s", organization.Name)
	suite.Logger.Record("org", organization).Infof("Organization Details")
}

func (suite *OrganizationSuite) TestOrganizationHasName() {
	organization, err := suite.Client.GetMyOrganization()
	suite.Require().Nil(err, "Failed to fetch my Organization")
	suite.Require().NotNil(organization, "Client's Organization is not loaded")
	suite.Assert().NotEmpty(organization.ID, "Client's Orgization has no ID")
	suite.Assert().NotEmpty(organization.Name, "Client's Orgization has no Name")
	suite.Assert().NotEmpty(organization.Features, "Client's Orgization has no features")
	suite.T().Logf("Organization: %s", organization.Name)
	suite.Logger.Record("org", organization).Infof("Organization Details")
}

// Suite Tools

func (suite *OrganizationSuite) SetupSuite() {
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
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *OrganizationSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *OrganizationSuite) BeforeTest(suiteName, testName string) {
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

func (suite *OrganizationSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

//go:build integration
// +build integration

package gcloudcx_test

import (
	"context"
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

	OrganizationName string
	OrganizationID   uuid.UUID
	Client           *gcloudcx.Client
}

func TestOrganizationSuite(t *testing.T) {
	suite.Run(t, new(OrganizationSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *OrganizationSuite) SetupSuite() {
	var err error
	var value string

	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	region := core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com")

	value = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_CLIENTID is not set")

	clientID, err := uuid.Parse(value)
	suite.Require().NoError(err, "PURECLOUD_CLIENTID is not a valid UUID")

	secret := core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	suite.Require().NotEmpty(secret, "PURECLOUD_CLIENTSECRET is not set")

	value = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_DEPLOYMENTID is not set")

	deploymentID, err := uuid.Parse(value)
	suite.Require().NoError(err, "PURECLOUD_DEPLOYMENTID is not a valid UUID")

	value = core.GetEnvAsString("ORGANIZATION_ID", "")
	suite.Require().NotEmpty(value, "ORGANIZATION_ID is not set in your environment")

	suite.OrganizationID, err = uuid.Parse(value)
	suite.Require().NoError(err, "ORGANIZATION_ID is not a valid UUID")

	suite.OrganizationName = core.GetEnvAsString("ORGANIZATION_NAME", "")
	suite.Require().NotEmpty(suite.OrganizationName, "ORGANIZATION_NAME is not set in your environment")

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
		err := suite.Client.Login(context.Background())
		suite.Require().NoError(err, "Failed to login")
		suite.Logger.Infof("Client is now logged in...")
	} else {
		suite.Logger.Infof("Client is already logged in...")
	}
}

func (suite *OrganizationSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *OrganizationSuite) TestCanFetchMyOrganization() {
	organization := gcloudcx.Organization{}
	err := suite.Client.Fetch(context.Background(), &organization)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoError(err, "Failed to fetch my Organization")
	suite.Assert().Equal(suite.OrganizationID, organization.GetID(), "Client's Organization ID is not the same")
	suite.Assert().Equal(suite.OrganizationName, organization.String(), "Client's Organization Name is not the same")
	suite.Assert().NotEmpty(organization.Features, "Client's Organization has no features")
	suite.T().Logf("Organization: %s", organization.Name)
	suite.Logger.Record("org", organization).Infof("Organization Details")
}

func (suite *OrganizationSuite) TestCanFetchOrganizationByID() {
	organization := gcloudcx.Organization{}
	err := suite.Client.Fetch(context.Background(), &organization, suite.OrganizationID)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoError(err, "Failed to fetch my Organization")
	suite.Assert().Equal(suite.OrganizationID, organization.GetID(), "Client's Organization ID is not the same")
	suite.Assert().Equal(suite.OrganizationName, organization.String(), "Client's Organization Name is not the same")
	suite.Assert().NotEmpty(organization.Features, "Client's Organization has no features")
	suite.T().Logf("Organization: %s", organization.Name)
	suite.Logger.Record("org", organization).Infof("Organization Details")
}

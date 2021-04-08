// +build integration

package purecloud_test

import (
	"net/url"

	purecloud "github.com/gildas/go-purecloud"
	"github.com/google/uuid"
)


func (suite *OpenMessagingSuite) TestCanFetchIntegrations() {
	integrations, err := purecloud.FetchOpenMessagingIntegrations(suite.Client)
	suite.Require().Nil(err, "Failed to fetch my Organization")
	suite.T().Logf("Found %d integrations", len(integrations))
	if len(integrations) > 0 {
		for _, integration := range integrations {
			suite.Logger.Record("integration", integration).Infof("Got a integration")
			suite.Assert().NotEmpty(integration.ID)
			suite.Assert().NotEmpty(integration.Name)
			suite.Assert().NotEmpty(integration.WebhookToken)
			suite.Assert().NotNil(integration.WebhookURL)
			suite.T().Logf("  Integration: %s (%s)", integration.Name, integration.ID)
		}
	}
}

func (suite *OpenMessagingSuite) TestCanCreateIntegration() {
	integration := &purecloud.OpenMessagingIntegration{}
	_ = integration.Initialize(suite.Client)
	name := "TEST-GO-PURECLOUD"
	webhookURL, _ := url.Parse("https://www.breizh.org/sim/purecloud")
	webhookToken := "DEADBEEF"
	err := integration.Create(name, webhookURL, webhookToken)
	suite.Require().Nil(err, "Failed to create integration")
	suite.Logger.Record("integration", integration).Infof("Created a integration")
}

func (suite *OpenMessagingSuite) TestCanDeleteIntegration() {
	integration := &purecloud.OpenMessagingIntegration{}
	err := integration.Initialize(suite.Client, uuid.MustParse("34071108-1569-4cb0-9137-a326b8a9e815"))
	suite.Require().Nil(err, "Failed to fetch integration")
	suite.Logger.Record("integration", integration).Infof("Got a integration")
	err = integration.Delete()
	suite.Require().Nil(err, "Failed to delete integration")
	err = integration.Initialize(suite.Client)
	suite.Require().NotNil(err, "Integration should not exist anymore")
}

package purecloud_test

import (
	"testing"

	"github.com/gildas/go-purecloud"
	"github.com/stretchr/testify/suite"
)


type URISuite struct {
	suite.Suite
	Name string
}

func TestURISuite(t *testing.T) {
	suite.Run(t, new(URISuite))
}

func (suite *URISuite) TestCanInstantiate() {
	uri := purecloud.NewURI("/api/v2/users/me")
	suite.Assert().Equal("/api/v2/users/me", uri.String())
}

func (suite *URISuite) TestCanInstantiateWithParameters() {
	uri := purecloud.NewURI("/api/v2/users/%s/status/%s", "me", "away")
	suite.Assert().Equal("/api/v2/users/me/status/away", uri.String())
}

func (suite *URISuite) TestCanJoinURIs() {
	uri1 := purecloud.NewURI("/api/v2")
	uri2 := purecloud.NewURI("/users/me")
	uri := uri1.Join(uri2)
	suite.Assert().Equal("/api/v2/users/me", uri.String())
}

func (suite *URISuite) TestHasPrefix() {
	suite.Assert().True(purecloud.NewURI("/api/v2/users/me").HasPrefix("/api/v2"))
	suite.Assert().False(purecloud.NewURI("/users/me").HasPrefix("/api/v2"))
}

func (suite *URISuite) TestHasProtocol() {
	suite.Assert().True(purecloud.NewURI("https://www.acme.com/api/v2/users/me").HasProtocol())
	suite.Assert().False(purecloud.NewURI("/users/me").HasProtocol())
}

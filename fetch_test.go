package gcloudcx_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type FetchSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestFetchSuite(t *testing.T) {
	suite.Run(t, new(FetchSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *FetchSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))
}

func (suite *FetchSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *FetchSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *FetchSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *FetchSuite) TestCanFetchObjectByID() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := Stuff{}
	err := client.Fetch(context.Background(), &stuff, idFromGCloud)
	suite.Require().NoError(err, "Failed to fetch stuff")
}

func (suite *FetchSuite) TestCanFetchObjectByName() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := Stuff{}
	err := client.Fetch(context.Background(), &stuff, nameFromGCloud)
	suite.Require().NoError(err, "Failed to fetch stuff")
}

func (suite *FetchSuite) TestShouldFailFetchingObjectWithUnknownID() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := Stuff{}
	err := client.Fetch(context.Background(), &stuff, uuid.New())
	suite.Require().Error(err, "Failed to fetch stuff")
	suite.Assert().ErrorIs(err, errors.NotFound)
	// TODO: Check error has the unknown ID
}

func (suite *FetchSuite) TestShouldFailFetchingObjectWithUnknownName() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := Stuff{}
	err := client.Fetch(context.Background(), &stuff, "unknown")
	suite.Require().Error(err, "Failed to fetch stuff")
	suite.Assert().ErrorIs(err, errors.NotFound)
	// TODO: Check error has the unknown name
}

type Stuff struct {
	ID      uuid.UUID        `json:"id"`
	Name    string           `json:"name"`
	SelfURI gcloudcx.URI     `json:"selfUri"`
	client  *gcloudcx.Client `json:"-"`
	log     *logger.Logger   `json:"-"`
}

var ( // simulate the fetching of the object via Genesys Cloud
	idFromGCloud   = uuid.MustParse("d4f8f8f8-f8f8-f8f8-f8f8-f8f8f8f8f8f8")
	nameFromGCloud = "stuff"
	uriFromGCloud  = gcloudcx.URI(fmt.Sprintf("/api/v2/stuff/%s", idFromGCloud))
)

// GetID gets the identifier of this
//   implements Identifiable
func (stuff Stuff) GetID() uuid.UUID {
	return stuff.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (stuff Stuff) GetURI() gcloudcx.URI {
	return stuff.SelfURI
}

func (stuff *Stuff) Fetch(ctx context.Context, client *gcloudcx.Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, stuff, parameters...)

	if id != uuid.Nil { // Here, we should fetch the object from Genesys Cloud
		if id != idFromGCloud {
			return errors.NotFound.With("id", id.String())
		}
		stuff.ID = idFromGCloud
		stuff.Name = nameFromGCloud
		stuff.SelfURI = uriFromGCloud
		stuff.client = client
		stuff.log = log
		return nil
	}
	if len(name) > 0 { // Here, we should fetch the object from Genesys Cloud
		if name != nameFromGCloud {
			return errors.NotFound.With("name", name)
		}
		stuff.ID = idFromGCloud
		stuff.Name = nameFromGCloud
		stuff.SelfURI = uriFromGCloud
		stuff.client = client
		stuff.log = log
		return nil
	}
	if len(selfURI) > 0 { // Here, we should fetch the object from Genesys Cloud
		if selfURI != uriFromGCloud {
			return errors.NotFound.With("selfURI", selfURI.String())
		}
		stuff.ID = idFromGCloud
		stuff.Name = nameFromGCloud
		stuff.SelfURI = uriFromGCloud
		stuff.client = client
		stuff.log = log
		return nil
	}
	return errors.NotFound.WithStack()
}

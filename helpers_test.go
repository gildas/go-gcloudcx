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

type HelpersSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestHelpersSuite(t *testing.T) {
	suite.Run(t, new(HelpersSuite))
}
func (suite *HelpersSuite) TestCanRunInitializable() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := &Stuff{}
	err := stuff.Initialize(client)
	suite.Assert().Nil(err, "Failed to fetch stuff")
}

func (suite *HelpersSuite) TestCanInitializeWithFetch() {
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		DeploymentID: uuid.New(),
		Logger:       suite.Logger,
	})

	stuff := &Stuff{}
	err := client.Fetch(context.Background(), stuff)
	suite.Assert().Nil(err, "Failed to fetch stuff")
}

type Stuff struct {
	ID     string            `json:"id"`
	Client *gcloudcx.Client `json:"-"`
	Logger *logger.Logger    `json:"-"`
}

func (stuff *Stuff) Initialize(parameters ...interface{}) error {
	var (
		client *gcloudcx.Client
		log    *logger.Logger
	)

	for _, parameter := range parameters {
		switch object := parameter.(type) {
		case *gcloudcx.Client:
			client = object
		case *logger.Logger:
			log = object
		}
	}
	stuff.Client = client
	if log != nil {
		stuff.Logger = log
	} else {
		stuff.Logger = client.Logger.Child("stuff", "stuff")
	}
	if stuff.Client == nil {
		return errors.ArgumentMissing.WithStack()
	}
	return nil
}

// Suite Tools

func (suite *HelpersSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:        fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:  true,
			FilterLevel: logger.TRACE,
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))
}

func (suite *HelpersSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *HelpersSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *HelpersSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

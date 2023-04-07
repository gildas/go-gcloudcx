package gcloudcx_test

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type DataTableSuite struct {
	suite.Suite
	Name    string
	Start   time.Time
	Context context.Context
	Logger  *logger.Logger

	Client  *gcloudcx.Client
	TableID uuid.UUID
}

func TestDataTableSuite(t *testing.T) {
	suite.Run(t, new(DataTableSuite))
}

// *****************************************************************************
// Suite Tools

func (suite *DataTableSuite) SetupSuite() {
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
	suite.Context = suite.Logger.ToContext(context.Background())
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region   = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret   = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	)

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: region,
		Logger: suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")

	// Create a DataTable
	suite.TableID = uuid.MustParse("a5c9b6e1-9a36-405d-8728-67fc31973e5e")
}

func (suite *DataTableSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
		// Delete the DataTable
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *DataTableSuite) BeforeTest(suiteName, testName string) {
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

func (suite *DataTableSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *DataTableSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	if err != nil {
		panic(err)
	}
	return data
}

func (suite *DataTableSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************
// Suite Tests

func (suite *DataTableSuite) TestCanMarshal() {
	expected := suite.LoadTestData("datatable.json")
	table := gcloudcx.DataTable{
		ID:          suite.TableID,
		Name:        "Unit Test",
		Description: "unit testing",
	}
	payload, err := json.Marshal(table)
	suite.Require().NoError(err, "Failed to marshal DataTable")
	suite.Assert().JSONEq(string(expected), string(payload), "DataTable JSON is wrong")
}

func (suite *DataTableSuite) TestCanFetchByID() {
	table, err := gcloudcx.Fetch[gcloudcx.DataTable](suite.Context, suite.Client, suite.TableID)
	suite.Require().NoError(err, "Failed to fetch DataTable")
	suite.Require().NotNil(table, "DataTable is nil")
	suite.Assert().Equal(suite.TableID, table.GetID(), "DataTable ID is wrong")
	suite.Assert().Equal("Unit Test", table.String(), "DataTable stringer is wrong")
	suite.Assert().Equal("Unit Test", table.Name, "DataTable name is wrong")
	suite.Assert().Equal("unit testing", table.Description, "DataTable description is wrong")
	if table.Division != nil {
		suite.Assert().Equal("Home", table.Division.Name, "DataTable Division name is wrong")
	}
}

func (suite *DataTableSuite) TestCanFetchByName() {
	table, err := gcloudcx.FetchBy(suite.Context, suite.Client, func(dt gcloudcx.DataTable) bool { return dt.Name == "Unit Test" })
	suite.Require().NoError(err, "Failed to fetch DataTable")
	suite.Require().NotNil(table, "DataTable is nil")
	suite.Assert().Equal(suite.TableID, table.GetID(), "DataTable ID is wrong")
	suite.Assert().Equal("Unit Test", table.String(), "DataTable stringer is wrong")
	suite.Assert().Equal("Unit Test", table.Name, "DataTable name is wrong")
	suite.Assert().Equal("unit testing", table.Description, "DataTable description is wrong")
	if table.Division != nil {
		suite.Assert().Equal("Home", table.Division.Name, "DataTable Division name is wrong")
	}
}

func (suite *DataTableSuite) TestFetchShoudFailWithUnknownID() {
	table, err := gcloudcx.Fetch[gcloudcx.DataTable](suite.Context, suite.Client, uuid.New())
	suite.Require().Error(err, "DataTable should not be found")
	suite.Require().Nil(table, "DataTable should nil")
	suite.Logger.Errorf("Expected error:", err)
	// suite.Assert().ErrorIs(err, gcloudcx.NotFoundError, "Error should be NotFound")
}

func (suite *DataTableSuite) TestFetchShoudFailWithUnknownName() {
	table, err := gcloudcx.FetchBy(suite.Context, suite.Client, func(dt gcloudcx.DataTable) bool { return dt.Name == "ZZZZZZZZ" })
	suite.Require().Error(err, "DataTable should not be found")
	suite.Require().Nil(table, "DataTable should nil")
	suite.Logger.Errorf("Expected error:", err)
	// suite.Assert().ErrorIs(err, errors.NotFound, "Error should be NotFound")
}

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
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type DataTableRowSuite struct {
	suite.Suite
	Name    string
	Start   time.Time
	Context context.Context
	Logger  *logger.Logger

	Client *gcloudcx.Client
	Table  *gcloudcx.DataTable
	Key    string
	Value  string
}

func TestDataTableRowSuite(t *testing.T) {
	suite.Run(t, new(DataTableRowSuite))
}

// *****************************************************************************
// Suite Tools

func (suite *DataTableRowSuite) SetupSuite() {
	var err error
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
		region        = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID      = uuid.MustParse(core.GetEnvAsString("PURECLOUD_CLIENTID", ""))
		secret        = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		correlationID string
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
	tableID := uuid.MustParse("a5c9b6e1-9a36-405d-8728-67fc31973e5e")
	suite.Table, correlationID, err = gcloudcx.Fetch[gcloudcx.DataTable](suite.Context, suite.Client, tableID)
	suite.Require().NoError(err, "Failed to fetch DataTable")
	suite.Require().NotNil(suite.Table, "DataTable is nil")
	suite.Logger.Infof("Correlation: %s", correlationID)

	suite.Key = "test-key"
	suite.Value = "test-value"
}

func (suite *DataTableRowSuite) TearDownSuite() {
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

func (suite *DataTableRowSuite) BeforeTest(suiteName, testName string) {
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

	// Create a Row
	correlationID, err := suite.Table.AddRow(suite.Context, suite.Key, gcloudcx.DataTableRow{"text": suite.Value})
	suite.Require().NoError(err, "Failed to add DataTable row")
	suite.Logger.Infof("Correlation: %s", correlationID)
}

func (suite *DataTableRowSuite) AfterTest(suiteName, testName string) {
	correlationID, err := suite.Table.DeleteRow(suite.Context, suite.Key)
	suite.Require().NoError(err, "Failed to delete DataTable row")
	suite.Logger.Infof("Correlation: %s", correlationID)

	duration := time.Since(suite.Start)
	if suite.T().Failed() {
		suite.Logger.Errorf("Test %s failed", testName)
	}
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *DataTableRowSuite) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	if err != nil {
		panic(err)
	}
	return data
}

func (suite *DataTableRowSuite) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// *****************************************************************************
// Suite Tests

func (suite *DataTableRowSuite) TestCanUpdateRow() {
	correlationID, err := suite.Table.UpdateRow(suite.Context, suite.Key, gcloudcx.DataTableRow{"text": "updated"})
	suite.Require().NoError(err, "Failed to update DataTable rows")
	suite.Logger.Infof("Correlation: %s", correlationID)

	row, correlationID, err := suite.Table.GetRow(suite.Context, suite.Key)
	suite.Require().NoError(err, "Failed to fetch DataTable rows")
	suite.Require().NotNil(row, "DataTable rows is nil")
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Logger.Record("row", row).Infof("DataTable row")
	suite.Assert().Equal(suite.Key, row["key"], "DataTable row key is wrong")
	suite.Assert().Equal("updated", row["text"], "DataTable row value is wrong")
}

func (suite *DataTableRowSuite) TestCanGetAllRows() {
	rows, correlationID, err := suite.Table.GetRows(suite.Context)
	suite.Require().NoError(err, "Failed to fetch DataTable rows")
	suite.Require().NotNil(rows, "DataTable rows is nil")
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Require().Len(rows, 1, "DataTable rows is wrong length")
	suite.Logger.Record("row", rows[0]).Infof("DataTable row 0")
	suite.Assert().Equal(suite.Key, rows[0]["key"], "DataTable row 0 key is wrong")
	suite.Assert().Equal(suite.Value, rows[0]["text"], "DataTable row 0 value is wrong")
}

func (suite *DataTableRowSuite) TestCanGetRowByKey() {
	row, correlationID, err := suite.Table.GetRow(suite.Context, suite.Key)
	suite.Require().NoError(err, "Failed to fetch DataTable rows")
	suite.Require().NotNil(row, "DataTable rows is nil")
	suite.Logger.Infof("Correlation: %s", correlationID)
	suite.Logger.Record("row", row).Infof("DataTable row")
	suite.Assert().Equal(suite.Key, row["key"], "DataTable row key is wrong")
	suite.Assert().Equal(suite.Value, row["text"], "DataTable row value is wrong")
}

func (suite *DataTableRowSuite) TestShouldFailGettingUnknowRow() {
	_, correlationID, err := suite.Table.GetRow(suite.Context, "unknown")
	suite.Require().Error(err, "Should have failed to fetch DataTable rows")
	suite.Logger.Infof("Correlation: %s", correlationID)
}

func (suite *DataTableRowSuite) TestShouldFailWithNilRow() {
	correlationID, err := suite.Table.AddRow(suite.Context, "nil", nil)
	suite.Require().Error(err, "Should have failed to add DataTable row")
	suite.Assert().ErrorIs(err, errors.ArgumentMissing, "Should have failed with ArgumentMissing")
	suite.Logger.Infof("Correlation: %s", correlationID)

	correlationID, err = suite.Table.UpdateRow(suite.Context, "nil", nil)
	suite.Require().Error(err, "Should have failed to update DataTable row")
	suite.Assert().ErrorIs(err, errors.ArgumentMissing, "Should have failed with ArgumentMissing")
	suite.Logger.Infof("Correlation: %s", correlationID)
}

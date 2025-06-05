package gcloudcx_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type SlotEntityTest struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time
}

func TestSlotEntityTest(t *testing.T) {
	suite.Run(t, new(SlotEntityTest))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *SlotEntityTest) SetupSuite() {
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
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))
}

func (suite *SlotEntityTest) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *SlotEntityTest) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()
}

func (suite *SlotEntityTest) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

func (suite *SlotEntityTest) LoadTestData(filename string) []byte {
	data, err := os.ReadFile(filepath.Join(".", "testdata", filename))
	suite.Require().NoErrorf(err, "Failed to Load Data. %s", err)
	return data
}

func (suite *SlotEntityTest) UnmarshalData(filename string, v interface{}) error {
	data := suite.LoadTestData(filename)
	suite.Logger.Infof("Loaded %s: %s", filename, string(data))
	return json.Unmarshal(data, v)
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *SlotEntityTest) TestCanUnmarshalEntities() {
	var data []json.RawMessage
	err := suite.UnmarshalData("slotentities.json", &data)
	suite.Require().NoError(err, "Failed to unmarshal SlotEntity")

	var slotEntities []gcloudcx.SlotEntity
	for _, raw := range data {
		entity, err := gcloudcx.UnmarshalSlotEntity(raw)
		suite.Require().NoError(err, "Failed to unmarshal SlotEntity from raw data: %s", string(raw))
		slotEntities = append(slotEntities, entity)
	}
	suite.Logger.Record("SlotEntities", slotEntities).Infof("Got %d SlotEntities", len(slotEntities))
}

func (suite *SlotEntityTest) TestCanMarshalEntities() {
	expected := suite.LoadTestData("slotentities.json")

	slotEntities := []gcloudcx.SlotEntity{
		gcloudcx.StringSlotEntity{Name: "ProductName", Value: "Chocolate Chip Cookie"},
		gcloudcx.IntegerSlotEntity{Name: "Size", Value: 12},
		gcloudcx.DecimalSlotEntity{Name: "Weight", Value: 85.6},
		gcloudcx.DurationSlotEntity{Name: "ConsumeBefore", Value: 16 * 24 * time.Hour},
		gcloudcx.BooleanSlotEntity{Name: "Diet", Value: false},
		gcloudcx.CurrencySlotEntity{Name: "CurrentPrice", Amount: 3.49, Currency: "USD"},
		gcloudcx.DatetimeSlotEntity{Name: "ExpiryDate", Value: time.Date(2024, 3, 15, 23, 59, 59, 0, time.UTC)},
		gcloudcx.StringCollectionSlotEntity{Name: "Ingredients", Values: []string{
			"flour",
			"sugar",
			"butter",
			"chocolate chips",
			"eggs",
		}},
		gcloudcx.IntegerCollectionSlotEntity{Name: "Presentations", Values: []int64{6, 12, 24}},
		gcloudcx.DecimalCollectionSlotEntity{Name: "AvailableWeights", Values: []float64{50, 85.5, 100.0}},
		gcloudcx.DurationCollectionSlotEntity{Name: "ShelLifeOptions", Values: []time.Duration{
			7 * 24 * time.Hour,
			15 * 24 * time.Hour,
			27 * 24 * time.Hour,
		}},
		gcloudcx.BooleanCollectionSlotEntity{Name: "ProductAttributes", Values: []bool{true, false, true}},
		gcloudcx.CurrencyCollectionSlotEntity{Name: "PreviousPrices", Values: []gcloudcx.CurrencyValue{
			{Amount: 3.49, Currency: "USD"},
			{Amount: 3.29, Currency: "USD"},
			{Amount: 2.99, Currency: "USD"},
		}},
		gcloudcx.DatetimeCollectionSlotEntity{Name: "BatchProductionDates", Values: []time.Time{
			time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 2, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 3, 10, 0, 0, 0, time.UTC),
		}},
	}

	payload, err := json.Marshal(slotEntities)
	suite.Require().NoError(err, "Failed to marshal SlotEntity")
	suite.Logger.Record("MarshaledPayload", string(payload)).Infof("Marshaled %d SlotEntities", len(slotEntities))
	suite.Assert().JSONEq(string(expected), string(payload), "The marshaled SlotEntity does not match the expected JSON")
}

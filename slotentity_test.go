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
	for _index, entity := range slotEntities {
		suite.Logger.Record("SlotEntity", entity).Infof("SlotEntity[%d]: %s (%s)", _index, entity.GetName(), entity)
		suite.Assert().NotEmpty(entity.GetName(), "SlotEntity name should not be empty")
		suite.Assert().NotEmpty(entity.GetType(), "SlotEntity type should not be empty")
	}
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

func (suite *SlotEntityTest) TestCanParseValue() {
	var err error

	_, err = gcloudcx.BooleanSlotEntity{Name: "IsAvailable"}.ParseValue("true")
	suite.Assert().NoError(err, "Failed to parse BooleanSlotEntity value")
	_, err = gcloudcx.CurrencySlotEntity{Name: "Price"}.ParseValue("3.49 USD")
	suite.Assert().NoError(err, "Failed to parse CurrencySlotEntity value")
	_, err = gcloudcx.DatetimeSlotEntity{Name: "CreatedAt"}.ParseValue("2024-03-15T12:00:00Z")
	suite.Assert().NoError(err, "Failed to parse DatetimeSlotEntity value")
	_, err = gcloudcx.DecimalSlotEntity{Name: "Weight"}.ParseValue("85.6")
	suite.Assert().NoError(err, "Failed to parse DecimalSlotEntity value")
	_, err = gcloudcx.DurationSlotEntity{Name: "Duration"}.ParseValue("PT16H")
	suite.Assert().NoError(err, "Failed to parse DurationSlotEntity value")
	_, err = gcloudcx.IntegerSlotEntity{Name: "Count"}.ParseValue("42")
	suite.Assert().NoError(err, "Failed to parse IntegerSlotEntity value")
	_, err = gcloudcx.StringSlotEntity{Name: "Description"}.ParseValue("A delicious cookie")
	suite.Assert().NoError(err, "Failed to parse StringSlotEntity value")
	_, err = gcloudcx.BooleanCollectionSlotEntity{Name: "Attributes"}.ParseValue("true, false, true")
	suite.Assert().NoError(err, "Failed to parse BooleanCollectionSlotEntity value")
	_, err = gcloudcx.CurrencyCollectionSlotEntity{Name: "Prices"}.ParseValue("3.49 USD, 2.99 USD")
	suite.Assert().NoError(err, "Failed to parse CurrencyCollectionSlotEntity value")
	_, err = gcloudcx.DatetimeCollectionSlotEntity{Name: "ProductionDates"}.ParseValue("2024-03-01T10:00:00Z, 2024-03-02T10:00:00Z")
	suite.Assert().NoError(err, "Failed to parse DatetimeCollectionSlotEntity value")
	_, err = gcloudcx.DecimalCollectionSlotEntity{Name: "AvailableWeights"}.ParseValue("50.0, 85.5, 100.0")
	suite.Assert().NoError(err, "Failed to parse DecimalCollectionSlotEntity value")
	_, err = gcloudcx.DurationCollectionSlotEntity{Name: "ShelfLife"}.ParseValue("P7D, P15D, P30D")
	suite.Assert().NoError(err, "Failed to parse DurationCollectionSlotEntity value")
	_, err = gcloudcx.IntegerCollectionSlotEntity{Name: "Sizes"}.ParseValue("6, 12, 24")
	suite.Assert().NoError(err, "Failed to parse IntegerCollectionSlotEntity value")
	_, err = gcloudcx.StringCollectionSlotEntity{Name: "Ingredients"}.ParseValue("flour, sugar, butter, chocolate chips, eggs")
	suite.Assert().NoError(err, "Failed to parse StringCollectionSlotEntity value")
}

func (suite *SlotEntityTest) TestShouldFailValidateWithEmptyName() {
	var err error

	booleanCollectionEntity := gcloudcx.BooleanCollectionSlotEntity{Name: ""}
	err = booleanCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	currencyCollectionEntity := gcloudcx.CurrencyCollectionSlotEntity{Name: ""}
	err = currencyCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	datetimeCollectionEntity := gcloudcx.DatetimeCollectionSlotEntity{Name: ""}
	err = datetimeCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	decimalCollectionEntity := gcloudcx.DecimalCollectionSlotEntity{Name: ""}
	err = decimalCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	durationCollectionEntity := gcloudcx.DurationCollectionSlotEntity{Name: ""}
	err = durationCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	integerCollectionEntity := gcloudcx.IntegerCollectionSlotEntity{Name: ""}
	err = integerCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	stringCollectionEntity := gcloudcx.StringCollectionSlotEntity{Name: ""}
	err = stringCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")

	booleanEntity := gcloudcx.BooleanSlotEntity{Name: ""}
	err = booleanEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	currencyEntity := gcloudcx.CurrencySlotEntity{Name: ""}
	err = currencyEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	datetimeEntity := gcloudcx.DatetimeSlotEntity{Name: ""}
	err = datetimeEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	decimalEntity := gcloudcx.DecimalSlotEntity{Name: ""}
	err = decimalEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	durationEntity := gcloudcx.DurationSlotEntity{Name: ""}
	err = durationEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	integerEntity := gcloudcx.IntegerSlotEntity{Name: ""}
	err = integerEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
	stringEntity := gcloudcx.StringSlotEntity{Name: ""}
	err = stringEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to empty name")
}

func (suite *SlotEntityTest) TestShouldFailValidateWithTooLongName() {
	var err error

	booleanCollectionEntity := gcloudcx.BooleanCollectionSlotEntity{Name: strings.Repeat("b", 101)} // 101 characters long
	err = booleanCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	currencyCollectionEntity := gcloudcx.CurrencyCollectionSlotEntity{Name: strings.Repeat("f", 101)} // 101 characters long
	err = currencyCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	datetimeCollectionEntity := gcloudcx.DatetimeCollectionSlotEntity{Name: strings.Repeat("e", 101)} // 101 characters long
	err = datetimeCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	decimalCollectionEntity := gcloudcx.DecimalCollectionSlotEntity{Name: strings.Repeat("d", 101)} // 101 characters long
	err = decimalCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	durationCollectionEntity := gcloudcx.DurationCollectionSlotEntity{Name: strings.Repeat("g", 101)} // 101 characters long
	err = durationCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	integerCollectionEntity := gcloudcx.IntegerCollectionSlotEntity{Name: strings.Repeat("c", 101)} // 101 characters long
	err = integerCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	stringCollectionEntity := gcloudcx.StringCollectionSlotEntity{Name: strings.Repeat("l", 101)} // 101 characters long
	err = stringCollectionEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")

	booleanEntity := gcloudcx.BooleanSlotEntity{Name: strings.Repeat("g", 101)} // 101 characters long
	err = booleanEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	currencyEntity := gcloudcx.CurrencySlotEntity{Name: strings.Repeat("i", 101)} // 101 characters long
	err = currencyEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	datetimeEntity := gcloudcx.DatetimeSlotEntity{Name: strings.Repeat("h", 101)} // 101 characters long
	err = datetimeEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	decimalEntity := gcloudcx.DecimalSlotEntity{Name: strings.Repeat("f", 101)} // 101 characters long
	err = decimalEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	durationEntity := gcloudcx.DurationSlotEntity{Name: strings.Repeat("j", 101)} // 101 characters long
	err = durationEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	integerEntity := gcloudcx.IntegerSlotEntity{Name: strings.Repeat("k", 101)} // 101 characters long
	err = integerEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
	stringEntity := gcloudcx.StringSlotEntity{Name: strings.Repeat("a", 101)} // 101 characters long
	err = stringEntity.Validate()
	suite.Assert().Error(err, "Expected validation to fail due to name length")
}

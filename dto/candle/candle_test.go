package candle_test

import (
	"testing"
	"time"

	"github.com/tekramkcots/sdk/dto/candle"
)

func TestNew(t *testing.T) {
	candle := candle.New(10)
	if candle.Close != 10 {
		t.Errorf("Expected candle.Close to be 10, got %f", candle.Close)
	}
	if candle.High != 10 {
		t.Errorf("Expected candle.High to be 10, got %f", candle.High)
	}
	if candle.Low != 10 {
		t.Errorf("Expected candle.Low to be 10, got %f", candle.Low)
	}
	if candle.Open != 10 {
		t.Errorf("Expected candle.Open to be 10, got %f", candle.Open)
	}
}

var candleUpdateTestcases = []struct {
	name          string
	startPrice    float32
	updatePrice   float32
	expectedOpen  float32
	expectedHigh  float32
	expectedLow   float32
	expectedClose float32
}{
	{"update with higher price", 10, 20, 10, 20, 10, 20},
	{"update with lower price", 10, 5, 10, 10, 5, 5},
}

func TestCandleUpdate(t *testing.T) {
	for _, tc := range candleUpdateTestcases {
		t.Run(tc.name, func(t *testing.T) {
			candle := candle.New(tc.startPrice)
			candle.Update(tc.updatePrice)
			if candle.Close != tc.expectedClose {
				t.Errorf("Expected candle.Close to be %f, got %f", tc.expectedClose, candle.Close)
			}
			if candle.High != tc.expectedHigh {
				t.Errorf("Expected candle.High to be %f, got %f", tc.expectedHigh, candle.High)
			}
			if candle.Low != tc.expectedLow {
				t.Errorf("Expected candle.Low to be %f, got %f", tc.expectedLow, candle.Low)
			}
			if candle.Open != tc.expectedOpen {
				t.Errorf("Expected candle.Open to be %f, got %f", tc.expectedOpen, candle.Open)
			}
		})
	}
}

var candleTypeIsMinuteCandleTestcases = []struct {
	name                   string
	candleType             candle.CandleType
	expectedIsMinuteCandle bool
}{
	{"minute candle", candle.Minute, true},
	{"five minute candle", candle.FiveMinute, true},
	{"ten minute candle", candle.TenMinute, true},
	{"fifteen minute candle", candle.FifteenMinute, true},
}

func TestCandleTypeIsMinuteCandle(t *testing.T) {
	for _, tc := range candleTypeIsMinuteCandleTestcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.candleType.IsMinuteCandle() != tc.expectedIsMinuteCandle {
				t.Errorf("Expected IsMinuteCandle to be %t, got %t", tc.expectedIsMinuteCandle, tc.candleType.IsMinuteCandle())
			}
		})
	}
}

var candleTypeGetStartTimeTestcases = []struct {
	name              string
	candleType        candle.CandleType
	startTime         time.Time
	expectedStartTime time.Time
}{
	{
		name:              "Minute - Normal",
		candleType:        candle.Minute,
		startTime:         time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
	},
	{
		name:              "Minute - offset",
		candleType:        candle.Minute,
		startTime:         time.Date(2020, time.January, 1, 1, 1, 30, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
	},
	{
		name:              "Five Minute - Normal",
		candleType:        candle.FiveMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
	},
	{
		name:              "Five Minute - offset",
		candleType:        candle.FiveMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
	},
	{
		name:              "Ten Minute - Normal",
		candleType:        candle.TenMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
	},
	{
		name:              "Ten Minute - offset",
		candleType:        candle.TenMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 11, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
	},
	{
		name:              "Fifteen Minute - Normal",
		candleType:        candle.FifteenMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 15, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 15, 0, 0, time.UTC),
	},
	{
		name:              "Fifteen Minute - offset",
		candleType:        candle.FifteenMinute,
		startTime:         time.Date(2020, time.January, 1, 1, 16, 0, 0, time.UTC),
		expectedStartTime: time.Date(2020, time.January, 1, 1, 15, 0, 0, time.UTC),
	},
}

func TestCandleTypeGetStartTime(t *testing.T) {
	for _, tc := range candleTypeGetStartTimeTestcases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := tc.candleType.GetStartTime(tc.startTime)
			if startTime != tc.expectedStartTime {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedStartTime, startTime)
			}
		})
	}
}

var candleDataNewtestcases = []struct {
	name               string
	candleType         candle.CandleType
	expectedCandleType candle.CandleType
	price              float32
	expectedClose      float32
	startTime          time.Time
	expectedStartTime  time.Time
}{
	{
		name:               "Minute - Normal",
		candleType:         candle.Minute,
		expectedCandleType: candle.Minute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
	},
	{
		name:               "Minute - offset",
		candleType:         candle.Minute,
		expectedCandleType: candle.Minute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 1, 30, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
	},
}

func TestCandleDataNewData(t *testing.T) {
	for _, tc := range candleDataNewtestcases {
		t.Run(tc.name, func(t *testing.T) {
			cd := candle.NewData(tc.candleType, tc.price, tc.startTime)
			if !cd.From.Equal(tc.expectedStartTime) {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedStartTime, cd.From)
			}
			if cd.Candle.Close != tc.expectedClose {
				t.Errorf("Expected close to be %f, got %f", tc.expectedClose, cd.Candle.Close)
			}
			if cd.Type != tc.expectedCandleType {
				t.Errorf("Expected candleType to be %d, got %d", tc.expectedCandleType, cd.Type)
			}
		})
	}
}

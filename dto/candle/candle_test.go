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
	startPrice    float64
	updatePrice   float64
	expectedOpen  float64
	expectedHigh  float64
	expectedLow   float64
	expectedClose float64
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

var candleNextTestcases = []struct {
	name          string
	startPrice    float64
	updatePrice   float64
	expectedOpen  float64
	expectedHigh  float64
	expectedLow   float64
	expectedClose float64
}{
	{"update with higher price", 10, 20, 20, 20, 20, 20},
	{"update with lower price", 10, 5, 5, 5, 5, 5},
}

func TestCandleNext(t *testing.T) {
	for _, tc := range candleNextTestcases {
		t.Run(tc.name, func(t *testing.T) {
			candle := candle.New(tc.startPrice)
			candle.Update(tc.updatePrice)
			candle = candle.Next()
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
	candleType             candle.Type
	expectedIsMinuteCandle bool
}{
	{"minute candle", candle.Minute, true},
	{"five minute candle", candle.FiveMinute, true},
	{"ten minute candle", candle.TenMinute, true},
	{"fifteen minute candle", candle.FifteenMinute, true},
}

func TestTypeIsMinuteCandle(t *testing.T) {
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
	candleType        candle.Type
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

func TestTypeGetStartTime(t *testing.T) {
	for _, tc := range candleTypeGetStartTimeTestcases {
		t.Run(tc.name, func(t *testing.T) {
			startTime := tc.candleType.GetStartTime(tc.startTime)
			if startTime != tc.expectedStartTime {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedStartTime, startTime)
			}
		})
	}
}

var candleTypeDurationTestcases = []struct {
	name             string
	candleType       candle.Type
	expectedDuration time.Duration
}{
	{"Minute", candle.Minute, time.Minute},
	{"Five Minute", candle.FiveMinute, time.Minute * 5},
	{"Ten Minute", candle.TenMinute, time.Minute * 10},
	{"Fifteen Minute", candle.FifteenMinute, time.Minute * 15},
}

func TestTypeDuration(t *testing.T) {
	for _, tc := range candleTypeDurationTestcases {
		t.Run(tc.name, func(t *testing.T) {
			duration := tc.candleType.Duration()
			if duration != tc.expectedDuration {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedDuration, duration)
			}
		})
	}
}

var candleDataNewtestcases = []struct {
	name               string
	candleType         candle.Type
	expectedCandleType candle.Type
	price              float64
	expectedClose      float64
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

var candleDataNextTestcases = []struct {
	name                   string
	candleType             candle.Type
	expectedNextCandleType candle.Type
	price                  float64
	updatePrice            float64
	expectedOpen           float64
	startTime              time.Time
	expectedStartTime      time.Time
}{
	{
		name:                   "Minute - Normal",
		candleType:             candle.Minute,
		expectedNextCandleType: candle.Minute,
		price:                  1.0,
		updatePrice:            2.0,
		expectedOpen:           2.0,
		startTime:              time.Date(2020, time.January, 1, 1, 1, 0, 0, time.UTC),
		expectedStartTime:      time.Date(2020, time.January, 1, 1, 2, 0, 0, time.UTC),
	},
	{
		name:                   "Minute - offset",
		candleType:             candle.Minute,
		expectedNextCandleType: candle.Minute,
		price:                  1.0,
		updatePrice:            2.0,
		expectedOpen:           2.0,
		startTime:              time.Date(2020, time.January, 1, 1, 1, 30, 0, time.UTC),
		expectedStartTime:      time.Date(2020, time.January, 1, 1, 2, 0, 0, time.UTC),
	},
}

func TestCandleDataNext(t *testing.T) {
	for _, tc := range candleDataNextTestcases {
		t.Run(tc.name, func(t *testing.T) {
			cd := candle.NewData(tc.candleType, tc.price, tc.startTime)
			cd.Candle.Update(tc.updatePrice)
			newCd := cd.Next()
			if !newCd.From.Equal(tc.expectedStartTime) {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedStartTime, newCd.From)
			}
			if newCd.Candle.Open != tc.expectedOpen {
				t.Errorf("Expected close to be %f, got %f", tc.expectedOpen, cd.Candle.Open)
			}
			if newCd.Type != tc.expectedNextCandleType {
				t.Errorf("Expected candleType to be %d, got %d", tc.expectedNextCandleType, newCd.Type)
			}
		})
	}
}

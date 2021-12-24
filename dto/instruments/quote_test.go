package instruments_test

import (
	"testing"
	"time"

	"github.com/tekramkcots/sdk/dto/candle"
	"github.com/tekramkcots/sdk/dto/instruments"
)

var quoteTickTestcases = []struct {
	name       string
	ltp        float64
	volume     uint32
	tickLtp    float64
	tickVolume uint32
}{
	{"normal", 10, 20, 10, 20},
}

func TestQuoteTick(t *testing.T) {
	for _, tc := range quoteTickTestcases {
		t.Run(tc.name, func(t *testing.T) {
			q := instruments.Quote{Token: 1, LTP: tc.ltp, Volume: tc.volume}
			tick := q.Tick()
			if tick.LastPrice != tc.tickLtp {
				t.Errorf("Expected tick.LastPrice to be %f, got %f", tc.tickLtp, tick.LastPrice)
			}
			if tick.VolumeTraded != tc.tickVolume {
				t.Errorf("Expected tick.VolumeTraded to be %d, got %d", tc.tickVolume, tick.VolumeTraded)
			}
		})
	}
}

var candleNextTestcases = []struct {
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
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 2, 0, 0, time.UTC),
	},
	{
		name:               "Minute - offset",
		candleType:         candle.Minute,
		expectedCandleType: candle.Minute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 1, 30, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 2, 0, 0, time.UTC),
	},
	{
		name:               "5 Minute - Normal",
		candleType:         candle.FiveMinute,
		expectedCandleType: candle.FiveMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 5, 0, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
	},
	{
		name:               "5 Minute - offset",
		candleType:         candle.FiveMinute,
		expectedCandleType: candle.FiveMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 7, 30, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
	},
	{
		name:               "10 Minute - Normal",
		candleType:         candle.TenMinute,
		expectedCandleType: candle.TenMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 10, 0, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 20, 0, 0, time.UTC),
	},
	{
		name:               "10 Minute - offset",
		candleType:         candle.TenMinute,
		expectedCandleType: candle.TenMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 17, 30, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 20, 0, 0, time.UTC),
	},
	{
		name:               "15 Minute - Normal",
		candleType:         candle.FifteenMinute,
		expectedCandleType: candle.FifteenMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 15, 0, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 30, 0, 0, time.UTC),
	},
	{
		name:               "15 Minute - offset",
		candleType:         candle.FifteenMinute,
		expectedCandleType: candle.FifteenMinute,
		price:              1.0,
		expectedClose:      1.0,
		startTime:          time.Date(2020, time.January, 1, 1, 17, 30, 0, time.UTC),
		expectedStartTime:  time.Date(2020, time.January, 1, 1, 30, 0, 0, time.UTC),
	},
}

func TestCandleNext(t *testing.T) {
	for _, tc := range candleNextTestcases {
		t.Run(tc.name, func(t *testing.T) {
			cd := instruments.Candle{Token: 1, Candle: candle.NewData(tc.candleType, tc.price, tc.startTime)}
			cd = cd.Next()
			if !cd.Candle.From.Equal(tc.expectedStartTime) {
				t.Errorf("Expected startTime to be %s, got %s", tc.expectedStartTime, cd.Candle.From)
			}
			if cd.Candle.Candle.Close != tc.expectedClose {
				t.Errorf("Expected close to be %f, got %f", tc.expectedClose, cd.Candle.Candle.Close)
			}
			if cd.Candle.Type != tc.expectedCandleType {
				t.Errorf("Expected candleType to be %d, got %d", tc.expectedCandleType, cd.Candle.Type)
			}
		})
	}
}

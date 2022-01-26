package listener_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tekramkcots/sdk/dto/candle"
	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/markets/indian"
	"github.com/tekramkcots/sdk/quotes/listener"
	"github.com/zerodha/gokiteconnect/v4/models"
)

var typeGetValueTestcases = []struct {
	name  string
	t     listener.Type
	tick  models.Tick
	value float64
}{
	{"ltp", listener.Price, models.Tick{LastPrice: 10}, 10},
	{"volume", listener.Volume, models.Tick{VolumeTraded: 10}, 10},
	{"zero ltp", listener.Price, models.Tick{VolumeTraded: 10}, 0},
	{"zero volume", listener.Volume, models.Tick{LastPrice: 10}, 0},
}

func TestTypeGetValue(t *testing.T) {
	for _, tc := range typeGetValueTestcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.t.GetValue(tc.tick) != tc.value {
				t.Errorf("Expected %f, got %f", tc.value, tc.t.GetValue(tc.tick))
			}
		})
	}
}

var newCandleTestcases = []struct {
	name string
	t    listener.Type
	c    candle.Type
}{
	{"price", listener.Price, candle.FiveMinute},
	{"volume", listener.Volume, candle.TenMinute},
}

func TestNewCandle(t *testing.T) {
	for _, tc := range newCandleTestcases {
		t.Run(tc.name, func(t *testing.T) {
			stopSignal := make(chan bool)
			candleChans := []chan []instruments.Candle{make(chan []instruments.Candle)}
			c := listener.NewCandle(tc.t, tc.c, candleChans, stopSignal)
			go c.Listen([]instruments.Quote{}, indian.MarketStartTime())
			if c.Type != tc.t {
				t.Errorf("Expected %d, got %d", tc.t, c.Type)
			}
			if c.C == nil {
				t.Error("Expected a channel, got nil")
			}
			if c.CandleType != tc.c {
				t.Errorf("Expected %d, got %d", tc.c, c.CandleType)
			}
			if len(c.Chans) != len(candleChans) {
				t.Errorf("Expected %d, got %d", len(candleChans), len(c.Chans))
			}
			stopSignal <- true
		})
	}
}

var now = time.Now()

var candleListenTestcases = []struct {
	name            string
	t               listener.Type
	c               candle.Type
	ins             []instruments.Quote
	ticks           [][]models.Tick
	tickesEvery     time.Duration
	expectedCandles map[uint32][]instruments.Candle
}{
	{
		name: "with multiple result candles",
		t:    listener.Price,
		c:    candle.Minute,
		ins: []instruments.Quote{
			{Token: 1, LTP: 100, Volume: 10000},
			{Token: 2, LTP: 500, Volume: 20000},
		},
		ticks: [][]models.Tick{
			{{InstrumentToken: 1, LastPrice: 100, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 500, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 200, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 600, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 300, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 700, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 400, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 800, VolumeTraded: 20000}},
		},
		tickesEvery: time.Second * 30,
		expectedCandles: map[uint32][]instruments.Candle{
			1: {
				{Token: 1, Candle: &candle.Data{
					From:   indian.MarketStartTime(),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 100, High: 200, Low: 100, Close: 200},
				}},
				{Token: 1, Candle: &candle.Data{
					From:   indian.MarketStartTime().Add(time.Minute),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 300, High: 400, Low: 300, Close: 400},
				}},
			},
			2: {
				{Token: 2, Candle: &candle.Data{
					From:   indian.MarketStartTime(),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 500, High: 600, Low: 500, Close: 600},
				}},
				{Token: 2, Candle: &candle.Data{
					From:   indian.MarketStartTime().Add(time.Minute),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 700, High: 800, Low: 700, Close: 800},
				}},
			},
		},
	},
	{
		name: "with different OHLC",
		t:    listener.Price,
		c:    candle.Minute,
		ins: []instruments.Quote{
			{Token: 1, LTP: 100, Volume: 10000},
			{Token: 2, LTP: 500, Volume: 20000},
		},
		ticks: [][]models.Tick{
			{{InstrumentToken: 1, LastPrice: 250, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 500, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 200, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 450, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 400, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 700, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 430, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 600, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 500, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 500, VolumeTraded: 20000}},
		},
		tickesEvery: time.Second * 19,
		expectedCandles: map[uint32][]instruments.Candle{
			1: {
				{Token: 1, Candle: &candle.Data{
					From:   indian.MarketStartTime(),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 250, High: 430, Low: 200, Close: 430},
				}},
			},
			2: {
				{Token: 2, Candle: &candle.Data{
					From:   indian.MarketStartTime(),
					Type:   candle.Minute,
					Candle: &candle.Stick{Open: 500, High: 700, Low: 450, Close: 600},
				}},
			},
		},
	},
}

func TestCandleListen(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(len(candleListenTestcases))
	for _, tc := range candleListenTestcases {
		go t.Run(tc.name, func(t *testing.T) {
			stopSignal1 := make(chan bool)
			stopSignal2 := make(chan bool)
			out := make(chan map[uint32][]instruments.Candle)
			candleChans := []chan []instruments.Candle{make(chan []instruments.Candle)}
			c := listener.NewCandle(tc.t, tc.c, candleChans, stopSignal1)
			go c.Listen(tc.ins, indian.MarketStartTime())
			go collectCandles(candleChans, stopSignal2, out)
			sendTicks(tc.ticks, c.C, tc.tickesEvery)
			stopSignal1 <- true
			stopSignal2 <- true
			candles := <-out
			err := verifyCandles(candles, tc.expectedCandles)
			if err != nil {
				t.Errorf("error verifying candles received via listener %s", err)
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func sendTicks(ticks [][]models.Tick, tickChan chan models.Tick, ticksEvery time.Duration) {
	for _, tickArray := range ticks {
		for _, tick := range tickArray {
			tickChan <- tick
		}
		time.Sleep(ticksEvery)
	}
}

func collectCandles(candleChans []chan []instruments.Candle, stopSignal chan bool, out chan map[uint32][]instruments.Candle) {
	candles := map[uint32][]instruments.Candle{}
	for {
		for _, candleChan := range candleChans {
			select {
			case <-stopSignal:
				out <- candles
				return
			case cns := <-candleChan:
				for _, cn := range cns {
					candles[cn.Token] = append(candles[cn.Token], cn)
				}
			}
		}
	}
}

func verifyCandles(resp map[uint32][]instruments.Candle, expectedCandles map[uint32][]instruments.Candle) error {
	if len(resp) != len(expectedCandles) {
		return fmt.Errorf("length of received candles doesn't match. expected %d. got %d", len(expectedCandles), len(resp))
	}
	for token, expCandles := range expectedCandles {
		actualCandles, ok := resp[token]
		if !ok {
			return fmt.Errorf("token %d not found in received candles", token)
		}
		for i, ec := range expCandles {
			ac := actualCandles[i]
			if !ac.Candle.From.Equal(ec.Candle.From) {
				return fmt.Errorf("candle start time doesn't match for index %d. expected %v. got %v", i, ec.Candle.From, ac.Candle.From)
			}
			if ac.Candle.Type != ec.Candle.Type {
				return fmt.Errorf("candle type doesn't match for index %d. expected %v. got %v", i, ec.Candle.Type, ac.Candle.Type)
			}
			if ac.Candle.Candle.Open != ec.Candle.Candle.Open {
				return fmt.Errorf("candle open value doesn't match for index %d. expected %v. got %v", i, ec.Candle.Candle.Open, ac.Candle.Candle.Open)
			}
			if ac.Candle.Candle.Close != ec.Candle.Candle.Close {
				return fmt.Errorf("candle close value doesn't match for index %d. expected %v. got %v", i, ec.Candle.Candle.Close, ac.Candle.Candle.Close)
			}
			if ac.Candle.Candle.High != ec.Candle.Candle.High {
				return fmt.Errorf("candle high doesn't match for index %d. expected %v. got %v", i, ec.Candle.Candle.High, ac.Candle.Candle.High)
			}
			if ac.Candle.Candle.Low != ec.Candle.Candle.Low {
				return fmt.Errorf("candle low doesn't match for index %d. expected %v. got %v", i, ec.Candle.Candle.Low, ac.Candle.Candle.Low)
			}
		}
	}
	return nil
}

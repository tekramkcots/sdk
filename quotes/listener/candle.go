package listener

import (
	"time"

	"github.com/tekramkcots/sdk/dto/candle"
	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type Type uint

const (
	Price  Type = 1
	Volume Type = 2
)

func (t Type) GetValue(tick models.Tick) float64 {
	switch t {
	case Price:
		return tick.LastPrice
	case Volume:
		return float64(tick.VolumeTraded)
	}
	return 0
}

type Candle struct {
	Type       Type
	CandleType candle.Type
	C          chan models.Tick
	Chans      []chan []instruments.Candle
	stopChan   chan bool
}

func NewCandle(listenerType Type, candleType candle.Type, chans []chan []instruments.Candle, stopSignal chan bool) *Candle {
	return &Candle{Type: listenerType, CandleType: candleType, C: make(chan models.Tick), Chans: chans, stopChan: stopSignal}
}

func (c *Candle) Listen(ins []instruments.Quote, startTime time.Time) {
	candles := map[uint32]instruments.Candle{}
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startTime.Location())
	for _, inst := range ins {
		candles[inst.Token] = instruments.Candle{Token: inst.Token, Candle: candle.NewData(c.CandleType, startTime)}
	}
	ticker := time.NewTicker(c.CandleType.Duration())
	for {
		select {
		case <-ticker.C:
			newCandles := []instruments.Candle{}
			for k, cand := range candles {
				candles[k] = cand.Next()
				newCandles = append(newCandles, cand)
			}
			for _, ch := range c.Chans {
				go func(ch chan []instruments.Candle) {
					ch <- newCandles
				}(ch)
			}
		case tick := <-c.C:
			ins, ok := candles[tick.InstrumentToken]
			if !ok {
				continue
			}
			ins.Candle.Candle.Update(float64(c.Type.GetValue(tick)))
			candles[tick.InstrumentToken] = ins
		case <-c.stopChan:
			ticker.Stop()
			return
		}
	}
}

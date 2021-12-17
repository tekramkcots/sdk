package candle

import "time"

type Candle struct {
	Open  float32
	High  float32
	Low   float32
	Close float32
}

func New(price float32) *Candle {
	return &Candle{
		Open:  price,
		High:  price,
		Low:   price,
		Close: price,
	}
}

func (c *Candle) Update(price float32) {
	c.Close = price
	if c.High < price {
		c.High = price
	}
	if c.Low > price {
		c.Low = price
	}
}

type CandleType uint

func (c CandleType) GetStartTime(t time.Time) time.Time {
	timeParameter := t.Minute()
	adjust := time.Duration(timeParameter % int(c))
	startTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	startTime = startTime.Add(-adjust * time.Minute)
	return startTime
}

func (c CandleType) IsMinuteCandle() bool {
	_, ok := minuteCandles[c]
	return ok
}

const (
	Minute        CandleType = 1
	FiveMinute    CandleType = 5
	TenMinute     CandleType = 10
	FifteenMinute CandleType = 15
)

var minuteCandles = map[CandleType]struct{}{
	Minute:        {},
	FiveMinute:    {},
	TenMinute:     {},
	FifteenMinute: {},
}

type CandleData struct {
	From   time.Time
	Type   CandleType
	Candle *Candle
}

func NewData(candleType CandleType, price float32, t time.Time) *CandleData {
	return &CandleData{
		From:   candleType.GetStartTime(t),
		Type:   candleType,
		Candle: New(price),
	}
}

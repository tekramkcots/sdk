package candle

import "time"

type Stick struct {
	Open    float64
	High    float64
	Low     float64
	Close   float64
	updated bool
}

func New() *Stick {
	return &Stick{}
}

func (c *Stick) Update(value float64) {
	c.Close = value
	if c.High < value {
		c.High = value
	}
	if c.Low > value {
		c.Low = value
	}
	if !c.updated {
		c.Open = value
		c.High = value
		c.Low = value
		c.Low = value
	}
	c.updated = true
}

func (c Stick) Next() *Stick {
	return New()
}

type Type uint

func (c Type) GetStartTime(t time.Time) time.Time {
	timeParameter := t.Minute()
	adjust := time.Duration(timeParameter % int(c))
	startTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	startTime = startTime.Add(-adjust * time.Minute)
	return startTime
}

func (c Type) Duration() time.Duration {
	isMinuteCandle := c.IsMinuteCandle()
	multiplier := time.Second
	if isMinuteCandle {
		multiplier = time.Minute
	}
	return multiplier * time.Duration(c)
}

func (c Type) IsMinuteCandle() bool {
	_, ok := minuteCandles[c]
	return ok
}

const (
	Minute        Type = 1
	FiveMinute    Type = 5
	TenMinute     Type = 10
	FifteenMinute Type = 15
)

var minuteCandles = map[Type]struct{}{
	Minute:        {},
	FiveMinute:    {},
	TenMinute:     {},
	FifteenMinute: {},
}

type Data struct {
	From   time.Time
	Type   Type
	Candle *Stick
}

func NewData(Type Type, t time.Time) *Data {
	return &Data{
		From:   Type.GetStartTime(t),
		Type:   Type,
		Candle: New(),
	}
}

func (c Data) Next() *Data {
	return &Data{
		From:   c.From.Add(c.Type.Duration()),
		Type:   c.Type,
		Candle: c.Candle.Next(),
	}
}

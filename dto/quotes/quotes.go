package quotes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zerodha/gokiteconnect/v4/models"
)

// DepthItem represents a single market depth entry.
type DepthItem struct {
	Price    float64
	Quantity uint32
	Orders   uint32
}

func (d DepthItem) Encode() []byte {
	return []byte(fmt.Sprintf("%f:%d:%d", d.Price, d.Quantity, d.Orders))
}

func (d DepthItem) Equal(dp DepthItem) bool {
	return d.Price == dp.Price && d.Quantity == dp.Quantity && d.Orders == dp.Orders
}

func DecodeDepthItem(data []byte) (*DepthItem, error) {
	parts := strings.Split(string(data), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid depth item. expected 3 parts. got %s", data)
	}
	p, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid depth item. expected float64 price. got %s", parts[0])
	}
	q, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid depth item. expected uint32 quantity. got %s", parts[1])
	}
	o, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid depth item. expected uint32 orders. got %s", parts[2])
	}
	return &DepthItem{p, uint32(q), uint32(o)}, nil
}

// Depth represents a group of buy/sell market depths.
type Depth struct {
	Buy  [5]DepthItem
	Sell [5]DepthItem
}

func (d Depth) Encode() []byte {
	buyEncode := fmt.Sprintf("%s|%s|%s|%s|%s", string(d.Buy[0].Encode()), string(d.Buy[1].Encode()), string(d.Buy[2].Encode()), string(d.Buy[3].Encode()), string(d.Buy[4].Encode()))
	sellEncode := fmt.Sprintf("%s|%s|%s|%s|%s", string(d.Sell[0].Encode()), string(d.Sell[1].Encode()), string(d.Sell[2].Encode()), string(d.Sell[3].Encode()), string(d.Sell[4].Encode()))
	return []byte(fmt.Sprintf("%s|%s", buyEncode, sellEncode))
}

func (d Depth) Equal(e Depth) bool {
	for i := 0; i < 5; i++ {
		if !d.Buy[i].Equal(e.Buy[i]) {
			return false
		}
		if !d.Sell[i].Equal(e.Sell[i]) {
			return false
		}
	}
	return true
}

func DecodeDepth(data []byte) (*Depth, error) {
	result := &Depth{}
	parts := strings.Split(string(data), "|")
	if len(parts) != 10 {
		return nil, fmt.Errorf("invalid depth. expected 10 parts. got %s", data)
	}
	for i := 0; i < 5; i++ {
		itemB, err := DecodeDepthItem([]byte(parts[i]))
		if err != nil {
			return nil, fmt.Errorf("invalid depth buy item at %d. expected DepthItem. got %w", i, err)
		}
		itemS, err := DecodeDepthItem([]byte(parts[i+5]))
		if err != nil {
			return nil, fmt.Errorf("invalid depth sell item at %d. expected DepthItem. got %w", i+5, err)
		}
		result.Buy[i] = *itemB
		result.Sell[i] = *itemS
	}
	return result, nil
}

type OHLC struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func (o OHLC) Encode() []byte {
	return []byte(fmt.Sprintf("%f:%f:%f:%f", o.Open, o.High, o.Low, o.Close))
}

func (o OHLC) Equal(e OHLC) bool {
	return o.Open == e.Open && o.High == e.High && o.Low == e.Low && o.Close == e.Close
}

func DecodeOHLC(data []byte) (*OHLC, error) {
	parts := strings.Split(string(data), ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid OHLC. expected 4 parts. got %s", data)
	}
	o, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid OHLC. expected float64 open. got %s", parts[0])
	}
	h, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid OHLC. expected float64 high. got %s", parts[1])
	}
	l, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid OHLC. expected float64 low. got %s", parts[2])
	}
	c, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid OHLC. expected float64 close. got %s", parts[3])
	}
	return &OHLC{o, h, l, c}, nil
}

type LTP struct {
	InstrumentToken uint32

	LastTradeTime time.Time
	LastPrice     float64
	VolumeTraded  uint32
	TotalBuy      uint32
	TotalSell     uint32
}

type Full struct {
	InstrumentToken uint32

	LastTradeTime time.Time
	LastPrice     float64
	VolumeTraded  uint32
	TotalBuy      uint32
	TotalSell     uint32

	OHLC  OHLC
	Depth Depth
}

func FromTick(tick models.Tick) Full {
	return Full{
		InstrumentToken: tick.InstrumentToken,
		LastTradeTime:   tick.LastTradeTime.Time,
		LastPrice:       tick.LastPrice,
		VolumeTraded:    tick.VolumeTraded,
		TotalBuy:        tick.TotalBuy,
		TotalSell:       tick.TotalSell,
		OHLC: OHLC{
			Open:  tick.OHLC.Open,
			High:  tick.OHLC.High,
			Low:   tick.OHLC.Low,
			Close: tick.OHLC.Close,
		},
		Depth: Depth{
			Buy: [5]DepthItem{
				{tick.Depth.Buy[0].Price, tick.Depth.Buy[0].Quantity, tick.Depth.Buy[0].Orders},
				{tick.Depth.Buy[1].Price, tick.Depth.Buy[1].Quantity, tick.Depth.Buy[1].Orders},
				{tick.Depth.Buy[2].Price, tick.Depth.Buy[2].Quantity, tick.Depth.Buy[2].Orders},
				{tick.Depth.Buy[3].Price, tick.Depth.Buy[3].Quantity, tick.Depth.Buy[3].Orders},
				{tick.Depth.Buy[4].Price, tick.Depth.Buy[4].Quantity, tick.Depth.Buy[4].Orders},
			},
			Sell: [5]DepthItem{
				{tick.Depth.Sell[0].Price, tick.Depth.Sell[0].Quantity, tick.Depth.Sell[0].Orders},
				{tick.Depth.Sell[1].Price, tick.Depth.Sell[1].Quantity, tick.Depth.Sell[1].Orders},
				{tick.Depth.Sell[2].Price, tick.Depth.Sell[2].Quantity, tick.Depth.Sell[2].Orders},
				{tick.Depth.Sell[3].Price, tick.Depth.Sell[3].Quantity, tick.Depth.Sell[3].Orders},
				{tick.Depth.Sell[4].Price, tick.Depth.Sell[4].Quantity, tick.Depth.Sell[4].Orders},
			},
		},
	}
}

func (f Full) LTPQuote() LTP {
	return LTP{
		InstrumentToken: f.InstrumentToken,
		LastTradeTime:   f.LastTradeTime,
		LastPrice:       f.LastPrice,
		VolumeTraded:    f.VolumeTraded,
		TotalBuy:        f.TotalBuy,
		TotalSell:       f.TotalSell,
	}
}

func (f Full) NatsChannelName() string {
	return fmt.Sprintf("f%d", f.InstrumentToken)
}

func (f Full) Encode() []byte {
	return []byte(fmt.Sprintf(
		"%d,%d,%f,%d,%d,%d,%s,%s",
		f.InstrumentToken,
		f.LastTradeTime.Unix(),
		f.LastPrice,
		f.VolumeTraded,
		f.TotalBuy,
		f.TotalSell,
		string(f.OHLC.Encode()),
		string(f.Depth.Encode()),
	))
}

func (f Full) Equal(e Full) bool {
	return f.InstrumentToken == e.InstrumentToken &&
		f.LastTradeTime.Equal(e.LastTradeTime) &&
		f.LastPrice == e.LastPrice &&
		f.VolumeTraded == e.VolumeTraded &&
		f.TotalBuy == e.TotalBuy &&
		f.TotalSell == e.TotalSell &&
		f.OHLC.Equal(e.OHLC) &&
		f.Depth.Equal(e.Depth)
}

func DecodeFull(data []byte) (*Full, error) {
	parts := strings.Split(string(data), ",")
	if len(parts) != 8 {
		return nil, fmt.Errorf("invalid Full quote. expected 8 parts. got %s", data)
	}
	token, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected uint32 instrument token. got %s", parts[0])
	}
	lt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected int64 last traded time. got %s", parts[1])
	}
	ltt := time.Unix(lt, 0)
	ltp, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected float64 last price. got %s", parts[2])
	}
	vt, err := strconv.ParseUint(parts[3], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected uint32 volume traded. got %s", parts[3])
	}
	tb, err := strconv.ParseUint(parts[4], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected uint32 total buy. got %s", parts[4])
	}
	ts, err := strconv.ParseUint(parts[5], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. expected uint32 total sell. got %s", parts[5])
	}
	o, err := DecodeOHLC([]byte(parts[6]))
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. %w", err)
	}
	d, err := DecodeDepth([]byte(parts[7]))
	if err != nil {
		return nil, fmt.Errorf("invalid Full quote. %w", err)
	}
	return &Full{uint32(token), ltt, ltp, uint32(vt), uint32(tb), uint32(ts), *o, *d}, nil
}

package instruments

import (
	"github.com/tekramkcots/sdk/markets/indian/index"
	"github.com/tekramkcots/sdk/models"
)

type InstrumentFetcher interface {
	Fetch() ([]Instrument, error)
}

type HistoricalData struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
	OI     int64
}

type Instrument struct {
	ModelID         uint                        `json:"-"`
	InstrumentToken int                         `json:"instrument_token"`
	ExchangeToken   int                         `json:"exchange_token"`
	Tradingsymbol   string                      `json:"tradingsymbol"`
	Name            string                      `json:"name"`
	LastPrice       float64                     `json:"last_price"`
	Expiry          int64                       `json:"expiry"`
	StrikePrice     float64                     `json:"strike"`
	TickSize        float64                     `json:"tick_size"`
	LotSize         float64                     `json:"lot_size"`
	InstrumentType  string                      `json:"instrument_type"`
	Segment         string                      `json:"segment"`
	Exchange        string                      `json:"exchange"`
	Candles         map[string][]HistoricalData `json:"-"`
}

func FromModel(ins models.Instrument) Instrument {
	return Instrument{
		ModelID:         ins.ID,
		InstrumentToken: ins.InstrumentToken,
		ExchangeToken:   ins.ExchangeToken,
		Tradingsymbol:   ins.Tradingsymbol,
		Name:            ins.Name,
		LastPrice:       ins.LastPrice,
		Expiry:          ins.Expiry,
		StrikePrice:     ins.StrikePrice,
		TickSize:        ins.TickSize,
		LotSize:         ins.LotSize,
		InstrumentType:  ins.InstrumentType,
		Segment:         ins.Segment,
		Exchange:        ins.Exchange,
	}
}

func FromModels(ins []models.Instrument) Instruments {
	var result []Instrument
	for _, i := range ins {
		result = append(result, FromModel(i))
	}
	return *NewInstruments(result)
}

func (i Instrument) ToModel() models.Instrument {
	return models.Instrument{
		InstrumentToken: i.InstrumentToken,
		ExchangeToken:   i.ExchangeToken,
		Tradingsymbol:   i.Tradingsymbol,
		Name:            i.Name,
		LastPrice:       i.LastPrice,
		Expiry:          i.Expiry,
		StrikePrice:     i.StrikePrice,
		TickSize:        i.TickSize,
		LotSize:         i.LotSize,
		InstrumentType:  i.InstrumentType,
		Segment:         i.Segment,
		Exchange:        i.Exchange,
	}
}

func (i Instruments) ToModels() []models.Instrument {
	var result []models.Instrument
	for _, i := range i.ins {
		result = append(result, i.ToModel())
	}
	return result
}

func (h HistoricalData) ToModel(instrumentID uint, candleType models.CandleType) models.HistoricalData {
	return models.HistoricalData{
		InstrumentID: instrumentID,
		CandleType:   candleType,
		Time:         h.Time,
		Open:         h.Open,
		High:         h.High,
		Low:          h.Low,
		Close:        h.Close,
		Volume:       h.Volume,
		OI:           h.OI,
	}
}

func FromHistoricalDataModel(h models.HistoricalData) HistoricalData {
	return HistoricalData{
		Time:   h.Time,
		Open:   h.Open,
		High:   h.High,
		Low:    h.Low,
		Close:  h.Close,
		Volume: h.Volume,
		OI:     h.OI,
	}
}

func FromHistoricalDataModels(hData []models.HistoricalData) []HistoricalData {
	result := []HistoricalData{}
	for _, h := range hData {
		result = append(result, FromHistoricalDataModel(h))
	}
	return result
}

func (i Instruments) HistoricalData(candleType models.CandleType) []models.HistoricalData {
	var result = []models.HistoricalData{}
	for _, ins := range i.ins {
		candles := ins.Candles[candleType.String()]
		for _, candle := range candles {
			result = append(result, candle.ToModel(ins.ModelID, candleType))
		}
	}
	return result
}

type Instruments struct {
	ins      []Instrument
	insTrMap map[string]Instrument
	insToMap map[int]Instrument
}

func NewInstruments(ins []Instrument) *Instruments {
	insTrMap := make(map[string]Instrument)
	insToMap := make(map[int]Instrument)
	for _, i := range ins {
		insTrMap[i.Tradingsymbol] = i
		insToMap[i.InstrumentToken] = i
	}
	return &Instruments{
		ins:      ins,
		insTrMap: insTrMap,
		insToMap: insToMap,
	}
}

func (i Instruments) Instruments() []Instrument {
	return i.ins
}

func (ins Instruments) WithInstrumentType(t index.Series) Instruments {
	var result []Instrument
	for _, i := range ins.ins {
		if i.InstrumentType == t.String() {
			result = append(result, i)
		}
	}
	return *NewInstruments(result)
}

func (ins Instruments) WithExchange(exch index.Exchange) Instruments {
	var result []Instrument
	for _, i := range ins.ins {
		if i.Exchange == exch.String() {
			result = append(result, i)
		}
	}
	return *NewInstruments(result)
}

func (i Instruments) WithIndexAndEquity() Instruments {
	var result []Instrument
	bnf := index.BankNifty()
	filtered := i.WithExchange(index.NSE).WithInstrumentType(index.EQ).ins
	for _, ins := range filtered {
		if len(ins.Name) == 0 {
			continue
		}
		if ins.LotSize != 1 && ins.LotSize != 0 {
			continue
		}
		if ins.Segment == index.Indices.String() && (ins.Tradingsymbol != bnf.GetSymbol()) {
			continue
		}
		result = append(result, ins)

	}
	return *NewInstruments(result)
}

func (ins Instruments) WithBankNiftyIndexAndEquity() Instruments {
	bnf := index.BankNifty()
	bnfIns := append(bnf.GetIndexStockSymbols(), bnf.GetSymbol())
	return ins.WithIndexAndEquity().WithTradingSymbols(bnfIns)
}

func (ins Instruments) WithAdjustedVolumeForBankNifty() Instruments {
	bnfIndex := -1
	volumeMap := map[string]map[int64]float32{}
	bnfSec := index.BankNifty()
	weightage := bnfSec.GetStockWeightage()
	for k, i := range ins.ins {
		if i.Tradingsymbol == bnfSec.GetSymbol() {
			bnfIndex = k
		}
		weight, ok := weightage[i.Tradingsymbol]
		if !ok {
			continue
		}
		for interval, candles := range i.Candles {
			candleCollection, ok := volumeMap[interval]
			if !ok {
				candleCollection = map[int64]float32{}
			}
			for _, candle := range candles {
				candleVolume, ok := candleCollection[candle.Time]
				if !ok {
					candleVolume = 0
				}
				candleVolume += (weight * float32(candle.Volume)) / 100
				candleCollection[candle.Time] = candleVolume
			}
			volumeMap[interval] = candleCollection
		}
	}
	if bnfIndex == -1 {
		return ins
	}
	bnf := ins.ins[bnfIndex]
	for interval, candleCollection := range bnf.Candles {
		for i := 0; i < len(candleCollection); i++ {
			candle := candleCollection[i]
			intervalMap, ok := volumeMap[interval]
			if !ok {
				continue
			}
			candleVolume, ok := intervalMap[candle.Time]
			if !ok {
				continue
			}
			candle.Volume = int64(candleVolume)
			candleCollection[i] = candle
		}
		bnf.Candles[interval] = candleCollection
	}
	ins.ins[bnfIndex] = bnf

	return *NewInstruments(ins.ins)
}

func (i Instruments) WithTradingSymbols(ts []string) Instruments {
	var result []Instrument
	for _, ins := range ts {
		val, ok := i.insTrMap[ins]
		if ok {
			result = append(result, val)
		}
	}
	return *NewInstruments(result)
}

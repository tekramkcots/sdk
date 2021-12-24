package instruments

import (
	"github.com/tekramkcots/sdk/dto/candle"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type Quote struct {
	Token  uint32
	LTP    float64
	Volume uint32
}

func (q Quote) Tick() models.Tick {
	return models.Tick{InstrumentToken: q.Token, LastPrice: q.LTP, VolumeTraded: q.Volume}
}

type Candle struct {
	Token  uint32
	Candle *candle.Data
}

func (c Candle) Next() Candle {
	return Candle{Token: c.Token, Candle: c.Candle.Next()}
}

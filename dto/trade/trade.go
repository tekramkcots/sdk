package trade

import "github.com/tekramkcots/sdk/dto/instruments"

type Trade struct {
	Buy        bool
	ExitTrade  bool
	ExitFor    string
	ID         string
	Instrument instruments.Instrument
	Price      float64
	Qty        uint64
	Time       int64
}

type PnL struct {
	StrategyID              string
	BuyPrice                float64
	SellPrice               float64
	BuyTrade                Trade
	SellTrade               Trade
	Instrument              instruments.Instrument
	Profit                  float64
	ProfitWithDelta         float64
	ProfitWithThetaAndDelta float64
}

package strategies

import (
	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/technicals/backtesting"
)

func GetAllStrategies(instrument instruments.Instrument) []backtesting.Strategy {
	allStratgies := []backtesting.Strategy{NewDailyTrade(instrument, 15), NewMACDStrategy(instrument, 15)}
	return allStratgies
}

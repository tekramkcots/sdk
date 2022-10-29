package strategies

import (
	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/technicals/backtesting"
)

func GetAllStrategies(instrument instruments.Instrument) []backtesting.Strategy {
	allStratgies := []backtesting.Strategy{NewDailyTrade(instrument, 25)}
	return allStratgies
}

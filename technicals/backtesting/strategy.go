package backtesting

import (
	"sync"
	"time"

	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/dto/trade"
)

type Strategy interface {
	ID() string
	Name() string
	Description() string
	Run(ch <-chan instruments.HistoricalData) <-chan trade.Trade
}

type Backtest struct {
	strategies    []Strategy
	quoteChannels []chan instruments.HistoricalData
	openTrades    map[string]map[string]trade.Trade
	wg            sync.WaitGroup
	l             sync.Mutex
}

func NewBacktest(strategies ...Strategy) *Backtest {
	quoteChannels := []chan instruments.HistoricalData{}
	for range strategies {
		quoteChannels = append(quoteChannels, make(chan instruments.HistoricalData))
	}
	return &Backtest{
		strategies:    strategies,
		openTrades:    make(map[string]map[string]trade.Trade),
		quoteChannels: quoteChannels,
	}
}

func (b *Backtest) Run(quoteChan <-chan instruments.HistoricalData) <-chan trade.PnL {
	pnlChannel := make(chan trade.PnL)
	go b.trackQuotes(quoteChan)
	for i, s := range b.strategies {
		b.wg.Add(1)
		trCh := s.Run(b.quoteChannels[i])
		go b.trackTrades(trCh, s.ID(), pnlChannel)
	}
	go b.finish(pnlChannel)
	return pnlChannel
}

func (b *Backtest) finish(ch chan trade.PnL) {
	b.wg.Wait()
	close(ch)
}

func (b *Backtest) trackQuotes(quoteChan <-chan instruments.HistoricalData) {
	for q := range quoteChan {
		for _, ch := range b.quoteChannels {
			sendQuote(q, ch)
		}
	}
	for _, ch := range b.quoteChannels {
		close(ch)
	}
}

func sendQuote(quote instruments.HistoricalData, ch chan<- instruments.HistoricalData) {
	ch <- quote
}

func (b *Backtest) trackTrades(trCh <-chan trade.Trade, id string, resultChannel chan<- trade.PnL) {
	for tr := range trCh {
		b.processTrade(tr, id, resultChannel)
	}
	b.wg.Done()
}

func (b *Backtest) processTrade(tr trade.Trade, id string, resultChannel chan<- trade.PnL) {
	b.l.Lock()
	defer b.l.Unlock()
	openTr, ok := b.openTrades[id]
	if !ok {
		openTr = make(map[string]trade.Trade)
	}
	openTrd, ok := openTr[tr.ExitFor]
	if tr.ExitTrade && !ok {
		return
	}
	if !tr.ExitTrade {
		openTr[tr.ID] = tr
		b.openTrades[id] = openTr
	} else {
		delete(openTr, tr.ExitFor)
		profit := (tr.Price - openTrd.Price) * float64(tr.Qty)
		if tr.IsLong {
			profit *= -1
		}
		withDelta := WithDelta(profit, tr.Time)
		withThetaAndDelta := WithTheta(withDelta, openTrd.Time, tr.Time)
		pnl := trade.PnL{
			StrategyID:              id,
			BuyPrice:                openTrd.Price,
			SellPrice:               tr.Price,
			BuyTrade:                openTrd,
			SellTrade:               tr,
			Instrument:              tr.Instrument,
			Profit:                  profit,
			ProfitWithDelta:         withDelta,
			ProfitWithThetaAndDelta: withThetaAndDelta,
		}
		b.wg.Add(1)
		go b.sendPnl(pnl, resultChannel)
	}
}

func (b *Backtest) sendPnl(pnl trade.PnL, resultChannel chan<- trade.PnL) {
	resultChannel <- pnl
	b.wg.Done()
}

func WithTheta(price float64, startTime, endTime int64) float64 {
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)
	hours := end.Sub(start).Hours()
	multilier := -1.0
	if price < 0 {
		multilier = 1
	}
	return price * (1 + (multilier * (hours * 0.04)))
}

func WithDelta(price float64, ts int64) float64 {
	t := time.Unix(ts, 0)
	switch t.Weekday() {
	case time.Thursday:
		return price * 0.7
	case time.Friday:
		return price * 0.6
	case time.Monday:
		return price * 0.7
	case time.Tuesday:
		return price * 0.8
	case time.Wednesday:
		return price * 0.9
	}
	return 0.0
}

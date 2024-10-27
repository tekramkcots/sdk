package strategies

import (
	"fmt"
	"sync"
	"time"

	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/dto/trade"
)

type openPosition struct {
	id     string
	qty    uint64
	time   int64
	isLong bool
}

// DailyTrade is an example strategy where we take the trade on the first quote of the day
// then exit it at 3:15 pm
type DailyTrade struct {
	id              string
	name            string
	description     string
	instrument      instruments.Instrument
	defaultQty      uint64
	currentPosition *openPosition
	currentDay      time.Time
	trades          int
	wg              sync.WaitGroup
	l               sync.Mutex
}

func NewDailyTrade(instrument instruments.Instrument, qty uint64) *DailyTrade {
	return &DailyTrade{
		id:          fmt.Sprintf("daily-trade-%s", instrument.Tradingsymbol),
		name:        fmt.Sprintf("Example Daily Trade on %s", instrument.Name),
		description: fmt.Sprintf("Take trade on the first quote and exit at the day's end with %d quantity", qty),
		instrument:  instrument,
		defaultQty:  qty,
	}
}

func (d *DailyTrade) ID() string {
	return d.id
}

func (d *DailyTrade) Name() string {
	return d.name
}

func (d *DailyTrade) Description() string {
	return d.description
}

func (d *DailyTrade) Run(ch <-chan instruments.HistoricalData) <-chan trade.Trade {
	trCh := make(chan trade.Trade)
	go d.startListeningToQuotes(ch, trCh)
	return trCh
}

func (d *DailyTrade) startListeningToQuotes(ch <-chan instruments.HistoricalData, trCh chan trade.Trade) {
	for q := range ch {
		d.processQuote(q, trCh)
	}
	d.wg.Wait()
	close(trCh)
}

func (d *DailyTrade) processQuote(q instruments.HistoricalData, trCh chan trade.Trade) {
	d.l.Lock()
	defer d.l.Unlock()
	t := time.Unix(q.Time, 0)
	dayT := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	if d.currentPosition == nil && dayT.After(d.currentDay) {
		d.currentPosition = &openPosition{
			qty:  d.defaultQty,
			time: q.Time,
			id:   fmt.Sprintf("%s-%d", d.ID(), d.trades+1),
		}
		d.currentDay = t
		d.wg.Add(1)
		go d.sendTrade(trade.Trade{
			IsLong:     true,
			ExitTrade:  false,
			ExitFor:    "",
			ID:         d.currentPosition.id,
			Instrument: d.instrument,
			Price:      q.Close,
			Qty:        d.currentPosition.qty,
			Time:       q.Time,
		}, trCh)
		d.trades += 1
		return
	}
	if d.currentPosition != nil && t.Hour() >= 15 && t.Minute() >= 15 {
		d.wg.Add(1)
		go d.sendTrade(trade.Trade{
			IsLong:     false,
			ExitTrade:  true,
			ExitFor:    d.currentPosition.id,
			ID:         fmt.Sprintf("%s-%d", d.ID(), d.trades+1),
			Instrument: d.instrument,
			Price:      q.Close,
			Qty:        d.currentPosition.qty,
			Time:       q.Time,
		}, trCh)
		d.trades += 1
		d.currentDay = t
		d.currentPosition = nil
	}
}

func (d *DailyTrade) sendTrade(trd trade.Trade, resultChannel chan<- trade.Trade) {
	resultChannel <- trd
	d.wg.Done()
}

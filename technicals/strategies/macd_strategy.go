package strategies

import (
	"fmt"
	"sync"

	"github.com/tekramkcots/sdk/dto/candle"
	"github.com/tekramkcots/sdk/dto/instruments"
	"github.com/tekramkcots/sdk/dto/trade"
	"github.com/tekramkcots/sdk/technicals/indicators"
)

// DailyTrade is an example strategy where we take the trade on the first quote of the day
// then exit it at 3:15 pm
type MACDStrategy struct {
	id              string
	name            string
	description     string
	instrument      instruments.Instrument
	defaultQty      uint64
	currentPosition *openPosition
	trades          int
	wg              sync.WaitGroup
	l               sync.Mutex
	closingValues   []float64
	candleType      candle.Type
	lastCandle      int64
}

func NewMACDStrategy(instrument instruments.Instrument, qty uint64) *MACDStrategy {
	return &MACDStrategy{
		id:            fmt.Sprintf("macd-%s", instrument.Tradingsymbol),
		name:          fmt.Sprintf("Take trade with macd and ema on %s", instrument.Name),
		description:   fmt.Sprintf("Take trade on the  %d quantity", qty),
		instrument:    instrument,
		defaultQty:    qty,
		closingValues: []float64{},
	}
}

func (m *MACDStrategy) ID() string {
	return m.id
}

func (m *MACDStrategy) Name() string {
	return m.name
}

func (m *MACDStrategy) Description() string {
	return m.description
}

func (m *MACDStrategy) Run(ch <-chan instruments.HistoricalData) <-chan trade.Trade {
	trCh := make(chan trade.Trade)
	go m.startListeningToQuotes(ch, trCh)
	return trCh
}

func (m *MACDStrategy) startListeningToQuotes(ch <-chan instruments.HistoricalData, trCh chan trade.Trade) {
	for q := range ch {
		m.processQuote(q, trCh)
	}
	m.wg.Wait()
	close(trCh)
}

func (m *MACDStrategy) processQuote(q instruments.HistoricalData, trCh chan trade.Trade) {
	/*
	 *
	 * Use MACD to identify the trades and use ema to verify the trend
	 *
	 */
	m.l.Lock()
	defer m.l.Unlock()
	// find out the incoming candles from last candle time and current candle time
	// check for 60 seconds is enough. Just in case added 5 sec for all
	if m.lastCandle != 0 && m.candleType.NotSet() && q.Time-m.lastCandle <= 65 {
		m.candleType = candle.Minute
	} else if m.lastCandle != 0 && m.candleType.NotSet() && q.Time-m.lastCandle <= 305 {
		m.candleType = candle.FiveMinute
	} else if m.lastCandle != 0 && m.candleType.NotSet() && q.Time-m.lastCandle <= 605 {
		m.candleType = candle.TenMinute
	} else if m.lastCandle != 0 && m.candleType.NotSet() && q.Time-m.lastCandle <= 905 {
		m.candleType = candle.FifteenMinute
	}
	m.lastCandle = q.Time
	m.closingValues = append(m.closingValues, q.Close)

	// if there is a current position and it is the 3rd candle post that
	// then close the position
	if m.currentPosition != nil && m.currentPosition.time+(4*60*int64(m.candleType)) <= q.Time {
		// exit the position
		m.wg.Add(1)
		go m.sendTrade(trade.Trade{
			IsLong:     !m.currentPosition.isLong,
			ExitTrade:  true,
			ExitFor:    m.currentPosition.id,
			ID:         fmt.Sprintf("%s-%d", m.ID(), m.trades+1),
			Instrument: m.instrument,
			Price:      q.Close,
			Qty:        m.currentPosition.qty,
			Time:       q.Time,
		}, trCh)
		m.trades += 1
		m.currentPosition = nil
	}

	if m.currentPosition != nil {
		return
	}
	macd := indicators.MovingAverageConvergenceDivergence(12, 26, m.closingValues)
	macdSignalLine := indicators.MovingAverageConvergenceDivergenceSignalLine(9, macd)
	zeroLine := make([]float64, len(macd))
	for i := range zeroLine {
		zeroLine[i] = macd[i] - macdSignalLine[i]
	}

	// macd strategy where if it crosses the signal line and is below 0, then buy
	// if it crosses the signal line and is above 0, then sell
	direction, crossed := checkIfSignalsCrosed(macd, macdSignalLine)
	if !crossed {
		return
	}

	currentMacd := macd[len(macd)-1]
	isBullishMomentum := direction == MACD_GOING_ABOVE_SIGNAL
	isBearishMomentum := direction == MACD_GOING_BELOW_SIGNAL
	if !isBullishMomentum && !isBearishMomentum {
		return
	}

	// also add a validation with 200 day ema trend line to see we are not betting against trend
	ema200 := indicators.ExponentialMovingAverage(200, m.closingValues)
	currentEma := ema200[len(ema200)-1]
	isUptrend := q.Close > currentEma
	isDowntrend := q.Close < currentEma
	if ema200[len(ema200)-1] == 0 {
		// since values are not generated, let us not take the trade
		return
	}

	bullishMomentumConfirmed := (isBullishMomentum && isUptrend && currentMacd < -8) || (isBullishMomentum && currentMacd < -11)
	bearishMomentumConfirmed := (isBearishMomentum && isDowntrend && currentMacd > 8) || (isBearishMomentum && currentMacd > 11)
	if !bullishMomentumConfirmed && !bearishMomentumConfirmed {
		return
	}

	// everything aligns, take the trade
	m.currentPosition = &openPosition{
		qty:    m.defaultQty,
		time:   q.Time,
		id:     fmt.Sprintf("%s-%d", m.ID(), m.trades+1),
		isLong: bullishMomentumConfirmed,
	}
	m.wg.Add(1)
	go m.sendTrade(trade.Trade{
		IsLong:     bullishMomentumConfirmed, // either bullish or bearish
		ExitTrade:  false,
		ExitFor:    "",
		ID:         m.currentPosition.id,
		Instrument: m.instrument,
		Price:      q.Close,
		Qty:        m.currentPosition.qty,
		Time:       q.Time,
	}, trCh)
	m.trades += 1
}

func (m *MACDStrategy) sendTrade(trd trade.Trade, resultChannel chan<- trade.Trade) {
	resultChannel <- trd
	m.wg.Done()
}

const (
	MACD_GOING_BELOW_SIGNAL = 1
	MACD_GOING_ABOVE_SIGNAL = 2
)

func checkIfSignalsCrosed(macd []float64, macdSignalLine []float64) (int, bool) {
	// checking if the macd is going above or below the signal line with reference of last 3 values
	if len(macd) < 3 {
		return 0, false
	}
	if len(macdSignalLine) < 3 {
		return 0, false
	}
	// 3rd to and 2nd to last values above and last two below
	mLen := len(macd)
	msLen := len(macdSignalLine)
	if macd[mLen-3] < macdSignalLine[msLen-3] && macd[mLen-2] < macdSignalLine[msLen-2] && macd[mLen-1] >= macdSignalLine[msLen-1] {
		return MACD_GOING_ABOVE_SIGNAL, true
	}

	// 3rd to and 2nd to last values below and last two above
	if macd[mLen-3] > macdSignalLine[msLen-3] && macd[mLen-2] > macdSignalLine[msLen-2] && macd[mLen-1] <= macdSignalLine[msLen-1] {
		return MACD_GOING_BELOW_SIGNAL, true
	}
	return 0, false
}

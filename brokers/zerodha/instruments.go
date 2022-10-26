package zerodha

import (
	"fmt"
	"time"

	"github.com/tekramkcots/sdk/dto/instruments"
)

func (c Client) DownloadInstruments() (*instruments.Instruments, error) {
	// Get margins
	ins, err := c.GetInstruments()
	if err != nil {
		return nil, fmt.Errorf("error fetching the instruments list %w", err)
	}

	result := []instruments.Instrument{}
	for _, i := range ins {
		result = append(result, instruments.Instrument{
			InstrumentToken: i.InstrumentToken,
			ExchangeToken:   i.ExchangeToken,
			Tradingsymbol:   i.Tradingsymbol,
			Name:            i.Name,
			LastPrice:       i.LastPrice,
			Expiry:          i.Expiry.Unix(),
			StrikePrice:     i.StrikePrice,
			TickSize:        i.TickSize,
			LotSize:         i.LotSize,
			InstrumentType:  i.InstrumentType,
			Segment:         i.Segment,
			Exchange:        i.Exchange,
		})
	}

	return instruments.NewInstruments(result), nil
}

func (c Client) DownloadHistoricalData(ins instruments.Instruments, from, to time.Time, interval string) (*instruments.Instruments, error) {
	insI := ins.Instruments()
	for j := 0; j < len(insI); j++ {
		i := insI[j]
		// Get historical data
		currentFrom := from
		currentTo := from.AddDate(0, 0, 100)
		if currentTo.After(to) {
			currentTo = to
		}
		allHistorical := []instruments.HistoricalData{}
		for currentFrom.Before(currentTo) {
			historicalData, err := c.GetHistoricalData(i.InstrumentToken, interval, currentFrom, currentTo, false, false)
			if err != nil {
				return nil, fmt.Errorf("error fetching the historical data for %s %w", i.Tradingsymbol, err)
			}
			for _, h := range historicalData {
				allHistorical = append(allHistorical, instruments.HistoricalData{
					Time:   h.Date.Unix(),
					Open:   h.Open,
					High:   h.High,
					Low:    h.Low,
					Close:  h.Close,
					Volume: int64(h.Volume),
					OI:     int64(h.OI),
				})
			}
			currentFrom = currentTo
			currentTo = currentFrom.AddDate(0, 0, 100)
			if currentTo.After(to) {
				currentTo = to
			}
		}
		if i.Candles == nil {
			i.Candles = map[string][]instruments.HistoricalData{}
		}
		i.Candles[interval] = allHistorical
		insI[j] = i
		c.logger.Println("downloaded historical data for", i.Tradingsymbol)
	}
	return instruments.NewInstruments(insI), nil
}

package models

import "gorm.io/gorm"

type Instrument struct {
	ID              uint `gorm:"primaryKey"`
	InstrumentToken int
	ExchangeToken   int
	Tradingsymbol   string
	Name            string
	LastPrice       float64
	Expiry          int64
	StrikePrice     float64
	TickSize        float64
	LotSize         float64
	InstrumentType  string
	Segment         string
	Exchange        string
}

func SaveInstruments(db *gorm.DB, instruments []Instrument) error {
	batchSize := 500
	nearestBatchSizeTotalLen := (len(instruments)/batchSize)*batchSize + batchSize
	for i := 0; i < nearestBatchSizeTotalLen; i += batchSize {
		end := i + batchSize
		if end > len(instruments) {
			end = len(instruments)
		}
		records := instruments[i:end]
		if err := db.Create(&records).Error; err != nil {
			return err
		}
	}
	return nil
}

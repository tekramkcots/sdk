package models

import (
	"time"

	"gorm.io/gorm"
)

type HistoricalData struct {
	ID           uint       `gorm:"primaryKey"`
	InstrumentID uint       `gorm:"index"`
	CandleType   CandleType `gorm:"index"`
	Time         int64      `gorm:"index"`
	Open         float64
	High         float64
	Low          float64
	Close        float64
	Volume       int64
	OI           int64
}

type CandleType uint

const (
	OneMinute     CandleType = 1
	FiveMinute    CandleType = 5
	FifteenMinute CandleType = 15
)

var candleStrMap = map[CandleType]string{
	OneMinute:     "1minute",
	FiveMinute:    "5minute",
	FifteenMinute: "15minute",
}

func (c CandleType) String() string {
	return candleStrMap[c]
}

func GetHistoricalFor(db gorm.DB, instrumentID uint, candleType CandleType, from, to time.Time) ([]HistoricalData, error) {
	var historicalData []HistoricalData
	if err := db.Where("instrument_id = ? AND candle_type = ? AND time >= ? AND time <= ?", instrumentID, candleType, from, to).Find(&historicalData).Error; err != nil {
		return nil, err
	}
	return historicalData, nil
}

func SaveHistoricalData(db *gorm.DB, historicalData []HistoricalData) error {
	batchSize := 500
	nearestBatchSizeTotalLen := (len(historicalData)/batchSize)*batchSize + batchSize
	for i := 0; i < nearestBatchSizeTotalLen; i += batchSize {
		end := i + batchSize
		if end > len(historicalData) {
			end = len(historicalData)
		}
		records := historicalData[i:end]
		if err := db.Create(&records).Error; err != nil {
			return err
		}
	}
	return nil
}

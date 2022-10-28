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

func GetAllInstruments(db *gorm.DB) ([]Instrument, error) {
	var instruments []Instrument
	err := db.Find(&instruments).Error
	if err != nil {
		return nil, err
	}
	return instruments, nil
}

func GetInstrumentFor(db *gorm.DB, symbol string) (*Instrument, error) {
	var instrument Instrument
	if err := db.Where("tradingsymbol = ?", symbol).First(&instrument).Error; err != nil {
		return nil, err
	}
	return &instrument, nil
}

func DeleteAllInstruments(db *gorm.DB) error {
	err := db.Where("1 = 1").Delete(&Instrument{}).Error
	if err != nil {
		return err
	}
	return nil
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

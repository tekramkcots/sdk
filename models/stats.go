package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type StatItem uint

const (
	StatHistoricalDataDownloadedTill StatItem = 1
	StatInstrumentsUpdatedOn         StatItem = 2
)

var statStrMap = map[StatItem]string{
	StatHistoricalDataDownloadedTill: "Historical Data Downloaded Till",
	StatInstrumentsUpdatedOn:         "Instruments Updated On",
}

func (s StatItem) String() string {
	return statStrMap[s]
}

type Stat struct {
	ID    uint     `gorm:"primaryKey"`
	Stat  StatItem `gorm:"index"`
	Value string
}

func GetAllStats(db *gorm.DB) ([]Stat, error) {
	stats := []Stat{}
	if err := db.Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func GetStat(db *gorm.DB, st StatItem) (*Stat, error) {
	stat := &Stat{}
	err := db.Where("stat = ?", st).First(stat).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return stat, nil
}

func SaveOrUpdateStat(db *gorm.DB, stat Stat) error {
	err := db.Save(&stat).Error
	if err != nil {
		return fmt.Errorf("error saving stat: %w", err)
	}
	return nil
}

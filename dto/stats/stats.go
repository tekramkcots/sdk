package stats

import "github.com/tekramkcots/sdk/models"

type Stat struct {
	Stat  uint
	Key   string
	Value string
}

func FromModel(s models.Stat) Stat {
	return Stat{
		Stat:  uint(s.Stat),
		Key:   s.Stat.String(),
		Value: s.Value,
	}
}

func FromModels(stats []models.Stat) []Stat {
	var s []Stat
	for _, stat := range stats {
		s = append(s, FromModel(stat))
	}
	return s
}

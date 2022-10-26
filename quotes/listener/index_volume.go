package listener

import (
	"github.com/tekramkcots/sdk/markets/indian/index"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type IndexVolume struct {
	index    index.Index
	stopChan chan bool
	C        chan models.Tick
	outChan  []chan models.Tick
	// compo
}

func NewIndexVolume(index index.Index, stopChan chan bool) *IndexVolume {
	return &IndexVolume{index: index, stopChan: stopChan, C: make(chan models.Tick)}
}

func (i *IndexVolume) Listen() {
	for {
		select {
		case <-i.stopChan:
			return
			// case tick := <-i.C:
			// if _, ok := i.index.Composition[tick.InstrumentToken]; ok {
			// 	i.outChan <- tick
			// }
		}
	}
}

package listener_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tekramkcots/sdk/dto/quotes"
	"github.com/tekramkcots/sdk/providers"
	"github.com/tekramkcots/sdk/quotes/listener"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type NatsPublisher struct {
	m             sync.Mutex
	encodedQuotes map[string][][]byte
}

func NewNatsPublisher() *NatsPublisher {
	return &NatsPublisher{sync.Mutex{}, map[string][][]byte{}}
}

func (n *NatsPublisher) Drain() error {
	return nil
}

func (n *NatsPublisher) Publish(channelName string, data []byte) error {
	c, ok := n.encodedQuotes[channelName]
	if !ok {
		c = [][]byte{}
	}
	c = append(c, data)
	n.m.Lock()
	n.encodedQuotes[channelName] = c
	n.m.Unlock()
	return nil
}

func (n *NatsPublisher) GetQuotes() (map[uint32][]quotes.Full, error) {
	result := map[uint32][]quotes.Full{}
	for _, encodedQuotes := range n.encodedQuotes {
		for _, encodedQuote := range encodedQuotes {
			quote, err := quotes.DecodeFull(encodedQuote)
			if err != nil {
				return nil, err
			}
			qts, ok := result[quote.InstrumentToken]
			if !ok {
				qts = []quotes.Full{}
			}
			qts = append(qts, *quote)
			result[quote.InstrumentToken] = qts
		}
	}
	return result, nil
}

var natsListenTestcases = []struct {
	name           string
	tickesEvery    time.Duration
	ticks          [][]models.Tick
	expectedQuotes map[uint32][]quotes.Full
}{
	{
		name:        "normal",
		tickesEvery: time.Second,
		ticks: [][]models.Tick{
			{{InstrumentToken: 1, LastPrice: 100, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 500, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 200, VolumeTraded: 20000},
				{InstrumentToken: 2, LastPrice: 600, VolumeTraded: 10000}},
			{{InstrumentToken: 1, LastPrice: 300, VolumeTraded: 30000},
				{InstrumentToken: 2, LastPrice: 700, VolumeTraded: 20000}},
			{{InstrumentToken: 1, LastPrice: 400, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 800, VolumeTraded: 50000}},
		},
		expectedQuotes: map[uint32][]quotes.Full{
			1: {
				{InstrumentToken: 1, LastPrice: 100, VolumeTraded: 10000},
				{InstrumentToken: 1, LastPrice: 200, VolumeTraded: 20000},
				{InstrumentToken: 1, LastPrice: 300, VolumeTraded: 30000},
				{InstrumentToken: 1, LastPrice: 400, VolumeTraded: 10000},
			},
			2: {
				{InstrumentToken: 2, LastPrice: 500, VolumeTraded: 20000},
				{InstrumentToken: 2, LastPrice: 600, VolumeTraded: 10000},
				{InstrumentToken: 2, LastPrice: 700, VolumeTraded: 20000},
				{InstrumentToken: 2, LastPrice: 800, VolumeTraded: 50000},
			},
		},
	},
}

func TestNatsListen(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(len(natsListenTestcases))
	for _, tc := range natsListenTestcases {
		go t.Run(tc.name, func(t *testing.T) {
			defer wg.Done()
			stopSignal := make(chan bool)
			np := NewNatsPublisher()
			appCtx := providers.GetAppContext()
			nts := listener.NewNatsFromProvider(np, stopSignal)
			go nts.Listen(*appCtx)
			sendTicksToNats(tc.ticks, nts.C, tc.tickesEvery)
			stopSignal <- true
			quotesArray, err := np.GetQuotes()
			if err != nil {
				t.Errorf("error fetching quotes from nats publisher %s", err)
				return
			}
			err = verifyQuotes(quotesArray, tc.expectedQuotes)
			if err != nil {
				t.Errorf("error verifying quotes received via nats publisher %s", err)
				return
			}
		})
	}
	wg.Wait()
}

func sendTicksToNats(ticks [][]models.Tick, tickChan chan models.Tick, ticksEvery time.Duration) {
	for _, ticksArray := range ticks {
		for _, tick := range ticksArray {
			tickChan <- tick
		}
		time.Sleep(ticksEvery)
	}
}

func verifyQuotes(expected map[uint32][]quotes.Full, result map[uint32][]quotes.Full) error {
	for instrumentToken, expectedQuotes := range expected {
		resultQuotes, ok := result[instrumentToken]
		if !ok {
			return fmt.Errorf("instrument token %d not found in result", instrumentToken)
		}
		if len(expectedQuotes) != len(resultQuotes) {
			return fmt.Errorf("instrument token %d expected %d quotes but received %d", instrumentToken, len(expectedQuotes), len(resultQuotes))
		}
		for i, expectedQuote := range expectedQuotes {
			if !expectedQuote.Equal(resultQuotes[i]) {
				return fmt.Errorf("instrument token %d expected quote at %d to be %v but received %v", instrumentToken, i, expectedQuote, resultQuotes[i])
			}
		}
	}
	return nil
}

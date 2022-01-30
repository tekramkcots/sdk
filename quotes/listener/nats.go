package listener

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/tekramkcots/sdk/app"
	"github.com/tekramkcots/sdk/dto/quotes"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type NatsPublisher interface {
	Drain() error
	Publish(channelName string, data []byte) error
}

type Nats struct {
	Conn     NatsPublisher
	C        chan models.Tick
	stopChan chan bool
}

func NewNats(url string, stopChan chan bool) (*Nats, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %w", err)
	}
	return &Nats{nc, make(chan models.Tick), stopChan}, nil
}

func NewNatsFromProvider(provider NatsPublisher, stopChan chan bool) *Nats {
	return &Nats{provider, make(chan models.Tick), stopChan}
}

func (n *Nats) Listen(appCtx app.Context) {
	defer n.Conn.Drain()
	for {
		select {
		case tick := <-n.C:
			go n.Publish(appCtx, tick)
		case <-n.stopChan:
			return
		}
	}
}

func (n *Nats) Publish(appCtx app.Context, tick models.Tick) {
	fullQuote := quotes.FromTick(tick)
	channelName := fullQuote.NatsChannelName()
	convertedQuote := fullQuote.Encode()
	err := n.Conn.Publish(channelName, convertedQuote)
	if err != nil {
		appCtx.Logger.Errorf("error publishing quote of %d to nats: %s", fullQuote.InstrumentToken, err)
		return
	}
}

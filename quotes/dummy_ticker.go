package quotes

import (
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type DummyTicker struct {
	onClose       func(code int, reason string)
	onConnect     func()
	onError       func(err error)
	onMessage     func(msgType string, msg []byte)
	onOrderUpdate func(order kiteconnect.Order)
	onReconnect   func(attempt int, delay time.Duration)
	onTick        func(tick models.Tick)
}

func NewDummyTicker() *DummyTicker {
	return &DummyTicker{}
}

func (d *DummyTicker) OnClose(f func(code int, reason string)) {
	d.onClose = f
}
func (d *DummyTicker) Close(code int, reason string) {
	d.onClose(code, reason)
}

func (d *DummyTicker) OnConnect(f func()) {
	d.onConnect = f
}
func (d *DummyTicker) Connect() {
	d.onConnect()
}

func (d *DummyTicker) OnError(f func(err error)) {
	d.onError = f
}
func (d *DummyTicker) Error(err error) {
	d.onError(err)
}

func (d *DummyTicker) OnMessage(f func(msgType string, msg []byte)) {
	d.onMessage = f
}
func (d *DummyTicker) Message(msgType string, msg []byte) {
	d.onMessage(msgType, msg)
}

func (d *DummyTicker) OnOrderUpdate(f func(order kiteconnect.Order)) {
	d.onOrderUpdate = f
}
func (d *DummyTicker) OrderUpdate(order kiteconnect.Order) {
	d.onOrderUpdate(order)
}

func (d *DummyTicker) OnReconnect(f func(attempt int, delay time.Duration)) {
	d.onReconnect = f
}
func (d *DummyTicker) Reconnect(attempt int, delay time.Duration) {
	d.onReconnect(attempt, delay)
}

func (d *DummyTicker) OnTick(f func(tick models.Tick)) {
	d.onTick = f
}
func (d *DummyTicker) Tick(tick models.Tick) {
	d.onTick(tick)
}

package quotes

import (
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	"github.com/zerodha/gokiteconnect/v4/models"
)

type Ticker interface {
	OnClose(f func(code int, reason string))
	OnConnect(f func())
	OnError(f func(err error))
	OnMessage(f func(msgType string, msg []byte))
	OnOrderUpdate(f func(order kiteconnect.Order))
	OnReconnect(f func(attempt int, delay time.Duration))
	OnTick(f func(tick models.Tick))
}

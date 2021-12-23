package quotes

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tekramkcots/sdk/app"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	"github.com/zerodha/gokiteconnect/v4/models"
)

func ConnectToQuotes(appCtx app.Context, t Ticker, chans []chan models.Tick, orderChans []chan kiteconnect.Order) {
	/*
	 * We will initialize the following fucntions
	 * 1. OnClose
	 * 2. OnConnect
	 * 3. OnError
	 * 4. OnMessage
	 * 5. OnOrderUpdate
	 * 6. OnReconnect
	 * 7. OnTick
	 */
	//onclose
	t.OnClose(func(code int, reason string) {
		appCtx.Logger.WithFields(logrus.Fields{
			"code":   code,
			"reason": reason,
		}).Error("closed ticker connection")
	})

	//onconnect
	t.OnConnect(func() {
		appCtx.Logger.Info("connected to ticker")
	})

	//onerror
	t.OnError(func(err error) {
		appCtx.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("error in ticker")
	})

	//onmessage
	t.OnMessage(func(msgType string, msg []byte) {
		appCtx.Logger.WithFields(logrus.Fields{
			"msgType":  msgType,
			"messsage": string(msg),
		}).Info("message received")
	})

	//onorderupdate
	t.OnOrderUpdate(onOrderUpdate(appCtx, orderChans))

	//onreconnect
	t.OnReconnect(func(attempt int, delay time.Duration) {
		appCtx.Logger.WithFields(logrus.Fields{
			"attempt": attempt,
			"delay":   delay,
		}).Info("reconnected to ticker")
	})

	//ontick
	t.OnTick(onTick(appCtx, chans))
}

func onOrderUpdate(appCtx app.Context, chans []chan kiteconnect.Order) func(tick kiteconnect.Order) {
	t := time.NewTicker(time.Minute * 5)
	tickChan := make(chan kiteconnect.Order)
	go forwardOrderUpdates(appCtx, tickChan, chans, t)

	return handleOrderUpdate(tickChan)
}

func handleOrderUpdate(ch chan kiteconnect.Order) func(tick kiteconnect.Order) {
	return func(tick kiteconnect.Order) {
		ch <- tick
	}
}

func forwardOrderUpdates(appCtx app.Context, in chan kiteconnect.Order, out []chan kiteconnect.Order, t *time.Ticker) {
	count := 0
	for {
		select {
		case <-t.C:
			appCtx.Logger.WithFields(logrus.Fields{
				"count": count,
			}).Info("order update received in 5 mins")
			count = 0
		case order := <-in:
			count++

			appCtx.Logger.WithFields(logrus.Fields{
				"order": order,
			}).Info("order updated")
			for _, ch := range out {
				go func(ch chan kiteconnect.Order) {
					ch <- order
				}(ch)
			}
		}
	}
}

func onTick(appCtx app.Context, chans []chan models.Tick) func(tick models.Tick) {
	t := time.NewTicker(time.Minute * 5)
	tickChan := make(chan models.Tick)
	go forwardTicks(appCtx, tickChan, chans, t)

	return handleTick(tickChan)
}

func handleTick(ch chan models.Tick) func(tick models.Tick) {
	return func(tick models.Tick) {
		ch <- tick
	}
}

func forwardTicks(appCtx app.Context, in chan models.Tick, out []chan models.Tick, t *time.Ticker) {
	count := 0
	for {
		select {
		case <-t.C:
			appCtx.Logger.WithFields(logrus.Fields{
				"count": count,
			}).Info("tick received in 5 mins")
			count = 0
		case tick := <-in:
			count++
			for _, ch := range out {
				go func(ch chan models.Tick) {
					ch <- tick
				}(ch)
			}
		}
	}
}

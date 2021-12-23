package quotes_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tekramkcots/sdk/providers"
	"github.com/tekramkcots/sdk/quotes"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	"github.com/zerodha/gokiteconnect/v4/models"
)

func TestConnectToQuotes(t *testing.T) {
	ticks := []models.Tick{
		{InstrumentToken: 1, LastPrice: 10},
	}
	orders := []kiteconnect.Order{
		{OrderID: "xxx", Status: "success"},
	}
	appCtx := providers.GetAppContext()
	dummyTicker := quotes.NewDummyTicker()
	tickListener := make(chan models.Tick)
	orderUpdateListener := make(chan kiteconnect.Order)
	success := false
	defer func() {
		if !success {
			t.Error("failed connect to quotes function as ticker didn't receive the quotes as expected")
		}
	}()
	quotes.ConnectToQuotes(*appCtx, dummyTicker, []chan models.Tick{tickListener}, []chan kiteconnect.Order{orderUpdateListener})
	dummyTicker.Connect()
	dummyTicker.Error(fmt.Errorf("error receiving quotes"))
	dummyTicker.Reconnect(1, time.Second)
	dummyTicker.OrderUpdate(orders[0])
	dummyTicker.Message("TestMessage", []byte("nothing"))
	dummyTicker.Tick(ticks[0])
	receivedTick := <-tickListener
	if receivedTick.InstrumentToken != ticks[0].InstrumentToken {
		t.Errorf("Expected %d, received %d", ticks[0].InstrumentToken, receivedTick.InstrumentToken)
	}
	if receivedTick.LastPrice != ticks[0].LastPrice {
		t.Errorf("Expected %f, received %f", ticks[0].LastPrice, receivedTick.LastPrice)
	}
	receivedOrder := <-orderUpdateListener
	if receivedOrder.OrderID != orders[0].OrderID {
		t.Errorf("Expected %s, received %s", orders[0].OrderID, receivedOrder.OrderID)
	}
	if receivedOrder.Status != orders[0].Status {
		t.Errorf("Expected %s, received %s", orders[0].Status, receivedOrder.Status)
	}
	dummyTicker.Close(0, "done")
	success = true
}

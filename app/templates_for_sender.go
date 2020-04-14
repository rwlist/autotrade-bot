package app

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/binance"
)

func errorMessage(err error, str string) string {
	return fmt.Sprintf("Error while %v:\n\n%s", str, err)
}

func startMessage(order binance.Order) string {
	return fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for %s", order.Side(), order.Price(), sleepDur)
}

func orderStatusMessage(order binance.Order) string {
	return fmt.Sprintf("Side: %v\nDone %v / %v\nStatus: %v", order.Side(), order.ExecutedQuantity(), order.OrigQuantity(), order.Status())
}

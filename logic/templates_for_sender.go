package logic

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/trade/trigger"

	"github.com/rwlist/autotrade-bot/trade/binance"
)

func errorMessage(err error, str string) string {
	return fmt.Sprintf("Error while %v:\n\n%s", str, err)
}

func startMessage(order *binance.Order) string {
	return fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for %s", order.Side, order.Price, sleepDur)
}

func orderStatusMessage(order *binance.Order) string {
	return fmt.Sprintf("Side: %v\nDone %v / %v\nStatus: %v", order.Side, order.ExecutedQuantity, order.OrigQuantity, order.Status)
}

func triggerResponseMessage(resp *trigger.Response) string {
	return fmt.Sprintf("Current rate: %v\nFormula rate: %.2f\n\n"+
		"Absolute difference: %.2f\nRelative difference: %.2f%%\n\n"+
		"Start rate: %v\nRelative profit: %.2f%%\nAbsolute profit: %.2f",
		resp.CurRate, resp.FormulaRate, resp.AbsDif, resp.RelDif, resp.StartRate, resp.RelProfit, resp.AbsProfit)
}

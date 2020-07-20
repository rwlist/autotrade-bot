package logic

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/trigger"
)

func startMessage(order *trade.Order) string {
	return fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for %s", order.Side, order.Price, sleepDur)
}

func orderStatusMessage(order *trade.Order) string {
	return fmt.Sprintf("Side: %v\nDone %v / %v\nStatus: %v", order.Side, order.ExecutedQuantity, order.OrigQuantity, order.Status)
}

type infoToSend struct {
	resp   *trigger.Response
	isTest bool
}

func triggerResponseMessage(inf infoToSend) string {
	txt := ""
	txt += fmt.Sprintf("Current rate: %v\nFormula rate: %.2f\n\n",
		inf.resp.CurRate, inf.resp.FormulaRate)
	txt += fmt.Sprintf("Absolute difference: %.2f\nRelative difference: %.2f%%\n\n",
		inf.resp.AbsDif, inf.resp.RelDif)
	txt += fmt.Sprintf("Start rate: %v\nRelative profit: %.2f%%\nAbsolute profit: %.2f\n\n",
		inf.resp.StartRate, inf.resp.RelProfit, inf.resp.AbsProfit)
	txt += fmt.Sprintf("Error: %v\nUpdate time: %v\n\n",
		inf.resp.Err, inf.resp.T.Format("02.01.2006 15.04.05"))
	txt += fmt.Sprintf("Formula: %v\n\n", inf.resp.Formula)

	testTxt := "РЕЖИМ ТОРГОВЛИ ВКЛЮЧЕН"
	if inf.isTest {
		testTxt = "ТЕСТОВЫЙ РЕЖИМ ВКЛЮЧЕН"
	}
	txt += testTxt

	return txt
}

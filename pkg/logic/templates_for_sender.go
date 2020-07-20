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
	txt := fmt.Sprintf("Current rate: %v\n", inf.resp.CurRate)
	txt += fmt.Sprintf("Formula rate: %.2f\n\n", inf.resp.FormulaRate)
	txt += fmt.Sprintf("Absolute difference: %.2f\n", inf.resp.AbsDif)
	txt += fmt.Sprintf("Relative difference: %.2f%%\n\n", inf.resp.RelDif)
	txt += fmt.Sprintf("Start rate: %v\n", inf.resp.StartRate)
	txt += fmt.Sprintf("Relative profit: %.2f%%\n", inf.resp.RelProfit)
	txt += fmt.Sprintf("Absolute profit: %.2f\n\n", inf.resp.AbsProfit)
	txt += fmt.Sprintf("Error: %v\n", inf.resp.Err)
	txt += fmt.Sprintf("Update time: %v\n\n", inf.resp.T.Format("02.01.2006 15.04.05"))
	txt += fmt.Sprintf("Formula: %v\n\n", inf.resp.Formula)

	testTxt := "РЕЖИМ ТОРГОВЛИ ВКЛЮЧЕН"
	if inf.isTest {
		testTxt = "ТЕСТОВЫЙ РЕЖИМ ВКЛЮЧЕН"
	}
	txt += testTxt

	return txt
}

package logic

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/trigger"
)

func startMessage(order *trade.Order) string {
	txt := fmt.Sprintf("A %v BTC/USDT order was placed with price = %s.\n", order.Side, order.Price)
	txt += fmt.Sprintf("Waiting for %s", sleepDur)
	return txt
}

func orderStatusMessage(order *trade.Order) string {
	txt := fmt.Sprintf("Side: %v\n", order.Side)
	txt += fmt.Sprintf("Done %s / %s\n", order.ExecutedQuantity, order.OrigQuantity)
	txt += fmt.Sprintf("Status: %v", order.Status)
	return txt
}

type infoToSend struct {
	resp   *trigger.Response
	isTest bool
}

func triggerResponseMessage(inf infoToSend) string {
	txt := fmt.Sprintf("Current rate: %s\n", inf.resp.CurRate)
	txt += fmt.Sprintf("Formula rate: %s\n\n", inf.resp.FormulaRate.Truncate(convert.UsefulShift))
	txt += fmt.Sprintf("Absolute difference: %s\n", inf.resp.AbsDif.Truncate(convert.UsefulShift))
	txt += fmt.Sprintf("Relative difference: %s%%\n\n", inf.resp.RelDif.Truncate(convert.UsefulShift))
	txt += fmt.Sprintf("Start rate: %s\n", inf.resp.StartRate)
	txt += fmt.Sprintf("Relative profit: %s%%\n", inf.resp.RelProfit.Truncate(convert.UsefulShift))
	txt += fmt.Sprintf("Absolute profit: %s\n\n", inf.resp.AbsProfit.Truncate(convert.UsefulShift))
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

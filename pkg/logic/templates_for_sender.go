package logic

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/trigger"
)

func startMessage(order *trade.Order) string {
	txt := fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\n", order.Side, order.Price.String())
	txt += fmt.Sprintf("Waiting for %s", sleepDur)
	return txt
}

func orderStatusMessage(order *trade.Order) string {
	txt := fmt.Sprintf("Side: %v\n", order.Side)
	txt += fmt.Sprintf("Done %v / %v\n", order.ExecutedQuantity.String(), order.OrigQuantity.String())
	txt += fmt.Sprintf("Status: %v", order.Status)
	return txt
}

type infoToSend struct {
	resp   *trigger.Response
	isTest bool
}

func triggerResponseMessage(inf infoToSend) string {
	txt := fmt.Sprintf("Current rate: %v\n", inf.resp.CurRate.String())
	txt += fmt.Sprintf("Formula rate: %v\n\n", inf.resp.FormulaRate.Truncate(convert.UsefulShift).String())
	txt += fmt.Sprintf("Absolute difference: %v\n", inf.resp.AbsDif.Truncate(convert.UsefulShift).String())
	txt += fmt.Sprintf("Relative difference: %v%%\n\n", inf.resp.RelDif.Truncate(convert.UsefulShift).String())
	txt += fmt.Sprintf("Start rate: %v\n", inf.resp.StartRate.String())
	txt += fmt.Sprintf("Relative profit: %v%%\n", inf.resp.RelProfit.Truncate(convert.UsefulShift).String())
	txt += fmt.Sprintf("Absolute profit: %v\n\n", inf.resp.AbsProfit.Truncate(convert.UsefulShift).String())
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

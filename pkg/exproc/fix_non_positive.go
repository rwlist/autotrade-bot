package exproc

import (
	"fmt"
	"strings"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/rwlist/autotrade-bot/pkg/money"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

// tryFixNonPositive accepts two active orders, which are
func (f *Finder) tryFixNonPositive(order1, order2 chatexsdk.Order) bool { //nolint:funlen,unused
	opts, err := f.tradeOpts.GetAll()
	if err != nil {
		log.WithError(err).Error("failed to get opts")
		return false
	}

	var info []string
	info = append(
		info,
		"Trying to fix non positives",
	)

	pair1 := pairOf(order1)
	limitOptionName := "limit." + pair1.Sell
	myLimit := opts.Decimal(limitOptionName)

	minSell := opts.Decimal("satoshi." + pair1.Sell)
	minBuy := opts.Decimal("satoshi." + pair1.Buy)

	fullRoundRate := money.One.DivRound(order1.Rate, money.Precision).DivRound(order2.Rate, money.Precision)

	if !fullRoundRate.GreaterThanOrEqual(money.One) {
		info = append(info, "Error: non-positive loop, nothing to do")
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}

	if !minSell.IsPositive() || !minBuy.IsPositive() {
		info = append(info, "Error: satoshi is not specified, check it")
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}

	if myLimit.LessThan(minSell) {
		info = append(info, "Error: current limit is insufficient")
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}

	isDust1 := order1.Amount.LessThanOrEqual(minBuy)
	isDust2 := order2.Amount.LessThanOrEqual(minSell)

	if !isDust1 && !isDust2 {
		info = append(info, "Error: both orders are not small enough")
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}

	amount1 := decimal.Min(minBuy, order1.Amount)
	amount2 := decimal.Min(minSell, order2.Amount)

	trade1, err := f.makeTrade(order1.ID, chatexsdk.TradeRequest{
		Amount: amount1,
		Rate:   order1.Rate,
	})
	if err != nil {
		log.WithError(err).Error("failed to make trade1")
		info = append(info, "Error(trade1): "+err.Error())
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}
	log.WithField("trade1", trade1).Info("made trade")

	// TODO: update limits

	info = append(
		info,
		"Limits are not updated, WATCH OUT!",
		fmt.Sprintf(
			"trade1 = %v, received = %v, amount = %v",
			trade1.ID,
			trade1.ReceivedAmount,
			trade1.Amount,
		),
	)

	if !isDust2 {
		info = append(info, "ALl ok! Make first trade, that's enough")
		f.sender.Send(strings.Join(info, "\n"))
		return true
	}

	trade2, err := f.makeTrade(order2.ID, chatexsdk.TradeRequest{
		Amount: amount2,
		Rate:   order2.Rate,
	})
	if err != nil {
		log.WithError(err).Error("failed to make trade2")
		info = append(info, "Error(trade2): "+err.Error())
		f.sender.Send(strings.Join(info, "\n"))
		return false
	}
	log.WithField("trade2", trade2).Info("made trade")

	// TODO: update limits?

	info = append(
		info,
		fmt.Sprintf(
			"trade2 = %v, received = %v, amount = %v",
			trade2.ID,
			trade2.ReceivedAmount,
			trade2.Amount,
		),
		"All ok! Made both trades",
	)

	f.sender.Send(strings.Join(info, "\n"))

	return true
}

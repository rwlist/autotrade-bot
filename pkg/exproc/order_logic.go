package exproc

import (
	"github.com/shopspring/decimal"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/rwlist/autotrade-bot/pkg/money"
)

type OrderCalc struct {
	MidAmount1  decimal.Decimal
	MidAmount2  decimal.Decimal
	MidAmount   decimal.Decimal
	StartAmount decimal.Decimal
	NextAmount  decimal.Decimal
	LastAmount  decimal.Decimal

	MinStartAmount *decimal.Decimal
}

// calcOrders calculates amounts to do trades with order1,
// and then with order2.
//
// order1 allows to exchange start -> next.
// order2 allows to exchange next -> start.
//
// order1 => next/start
// order2 => start/next
//
// Order X/Y allows to buy X, and sell Y.
// order.Amount is the maximum X amount can buy.
// order.Rate is how much Y you need to sell to get one X.
func (o OrderCalc) CalcTrades(order1, order2 chatexsdk.Order) OrderCalc {
	// next -> start
	// how much next can be bought in order1
	midAmount1 := order1.Amount

	// next <- last_start
	// how much next can be sold in order2
	midAmount2 := order2.Amount.Mul(order2.Rate)

	// take min
	midAmount := decimal.Min(midAmount1, midAmount2)

	// start
	startAmount := midAmount.Mul(order1.Rate)

	// apply additional min, if exists
	if o.MinStartAmount != nil {
		startAmount = decimal.Min(startAmount, *o.MinStartAmount)
	}

	// apply order1
	nextAmount := startAmount.DivRound(order1.Rate, money.Precision)

	// apply order2
	lastAmount := nextAmount.DivRound(order2.Rate, money.Precision)

	return OrderCalc{
		MidAmount1:  midAmount1,
		MidAmount2:  midAmount2,
		MidAmount:   midAmount,
		StartAmount: startAmount,
		NextAmount:  nextAmount,
		LastAmount:  lastAmount,
	}
}

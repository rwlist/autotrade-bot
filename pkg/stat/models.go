package stat

import "github.com/shopspring/decimal"

type Status struct {
	Total    decimal.Decimal
	Rate     decimal.Decimal
	Balances []Balance
}

type Balance struct {
	Asset  string
	USD    decimal.Decimal
	Free   decimal.Decimal
	Locked decimal.Decimal
}

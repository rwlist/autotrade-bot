package trade

import "github.com/shopspring/decimal"

type Order struct {
	Symbol           string
	OrderID          int64
	Price            decimal.Decimal
	OrigQuantity     decimal.Decimal
	ExecutedQuantity decimal.Decimal
	Status           string
	Side             string
}

type Status struct {
	Order *Order
	Done  bool
	Err   error
}

type Balance struct {
	Asset  string          `json:"asset"`
	Free   decimal.Decimal `json:"free"`
	Locked decimal.Decimal `json:"locked"`
}

type Rate struct {
	Rate     decimal.Decimal // how much Base in 1 Currency
	Currency string
	Base     string
}

type Opts struct {
	Symbol string
	Scale  string
}

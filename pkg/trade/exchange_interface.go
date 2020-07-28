package trade

import (
	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/shopspring/decimal"
)

type IExchange interface {
	AccountBalance() ([]Balance, error)
	AccountSymbolBalance(symbol string) (decimal.Decimal, error)
	BalanceToUSD(bal *Balance) (decimal.Decimal, error)
	GetRate(symbol ...string) (decimal.Decimal, error)
	BuyAll(symbol ...string) *Status
	SellAll(symbol ...string) *Status
	GetOrder(id int64, symbol ...string) (*Order, error)
	CancelOrder(id int64, symbol ...string) error
	GetKlines(opts ...draw.KlinesOpts) (*draw.Klines, error)
	SetScale(scale string)
	SetSymbol(symbol string)
}

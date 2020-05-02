package binance

type Order interface {
	Symbol() string
	OrderID() int64
	Price() string
	OrigQuantity() string
	ExecutedQuantity() string
	Status() string
	Side() string
}

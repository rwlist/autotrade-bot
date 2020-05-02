package binance

import "github.com/adshao/go-binance"

type OrderNew struct {
	Ord *binance.CreateOrderResponse
}

func (ord *OrderNew) Symbol() string {
	return ord.Ord.Symbol
}

func (ord *OrderNew) OrderID() int64 {
	return ord.Ord.OrderID
}

func (ord *OrderNew) Price() string {
	return ord.Ord.Price
}

func (ord *OrderNew) OrigQuantity() string {
	return ord.Ord.OrigQuantity
}

func (ord *OrderNew) ExecutedQuantity() string {
	return ord.Ord.ExecutedQuantity
}

func (ord *OrderNew) Status() string {
	return string(ord.Ord.Status)
}

func (ord *OrderNew) Side() string {
	return string(ord.Ord.Side)
}

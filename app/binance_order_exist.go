package app

import "github.com/adshao/go-binance"

type OrderExist struct {
	Ord *binance.Order
}

func (ord *OrderExist) Symbol() string {
	return ord.Ord.Symbol
}

func (ord *OrderExist) OrderID() int64 {
	return ord.Ord.OrderID
}

func (ord *OrderExist) Price() string {
	return ord.Ord.Price
}

func (ord *OrderExist) OrigQuantity() string {
	return ord.Ord.OrigQuantity
}

func (ord *OrderExist) ExecutedQuantity() string {
	return ord.Ord.ExecutedQuantity
}

func (ord *OrderExist) Status() string {
	return string(ord.Ord.Status)
}

func (ord *OrderExist) Side() string {
	return string(ord.Ord.Side)
}

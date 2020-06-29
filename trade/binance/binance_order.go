package binance

import "github.com/adshao/go-binance"

type Order struct {
	Symbol           string
	OrderID          int64
	Price            string
	OrigQuantity     string
	ExecutedQuantity string
	Status           string
	Side             string
}

func convertOrderToOrder(ord *binance.Order) *Order {
	return &Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            ord.Price,
		OrigQuantity:     ord.OrigQuantity,
		ExecutedQuantity: ord.ExecutedQuantity,
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

func convertCreateOrderResponseToOrder(ord *binance.CreateOrderResponse) *Order {
	return &Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            ord.Price,
		OrigQuantity:     ord.OrigQuantity,
		ExecutedQuantity: ord.ExecutedQuantity,
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

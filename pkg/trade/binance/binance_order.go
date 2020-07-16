package binance

import (
	"github.com/adshao/go-binance"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

func convertOrderToOrder(ord *binance.Order) *trade.Order {
	return &trade.Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            ord.Price,
		OrigQuantity:     ord.OrigQuantity,
		ExecutedQuantity: ord.ExecutedQuantity,
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

func convertCreateOrderResponseToOrder(ord *binance.CreateOrderResponse) *trade.Order {
	return &trade.Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            ord.Price,
		OrigQuantity:     ord.OrigQuantity,
		ExecutedQuantity: ord.ExecutedQuantity,
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

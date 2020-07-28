package binance

import (
	"github.com/adshao/go-binance"
	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

func convertOrderToOrder(ord *binance.Order) *trade.Order {
	return &trade.Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            convert.UnsafeDecimal(ord.Price),
		OrigQuantity:     convert.UnsafeDecimal(ord.OrigQuantity),
		ExecutedQuantity: convert.UnsafeDecimal(ord.ExecutedQuantity),
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

func convertCreateOrderResponseToOrder(ord *binance.CreateOrderResponse) *trade.Order {
	return &trade.Order{
		Symbol:           ord.Symbol,
		OrderID:          ord.OrderID,
		Price:            convert.UnsafeDecimal(ord.Price),
		OrigQuantity:     convert.UnsafeDecimal(ord.OrigQuantity),
		ExecutedQuantity: convert.UnsafeDecimal(ord.ExecutedQuantity),
		Status:           string(ord.Status),
		Side:             string(ord.Side),
	}
}

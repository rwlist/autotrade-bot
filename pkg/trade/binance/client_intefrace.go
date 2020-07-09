package binance

import goBinance "github.com/adshao/go-binance"

type Client interface {
	AccountBalance() ([]goBinance.Balance, error)
	ListPrices(symbol string) ([]*goBinance.SymbolPrice, error)
	CreateOrder(req *orderReq) (*goBinance.CreateOrderResponse, error)
	GetOrder(req *orderID) (*goBinance.Order, error)
	CancelOrder(req *orderID) (*goBinance.CancelOrderResponse, error)
	GetKlines(req *klinesReq) ([]*goBinance.Kline, error)
}

type orderReq struct {
	Symbol   string
	Side     goBinance.SideType
	Type     goBinance.OrderType
	Tif      goBinance.TimeInForceType
	Price    string
	Quantity string
}

type orderID struct {
	Symbol string
	ID     int64
}

type klinesReq struct {
	Symbol    string
	T         string
	StartTime int64
}

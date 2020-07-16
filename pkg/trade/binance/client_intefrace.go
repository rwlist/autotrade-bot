package binance

import gobinance "github.com/adshao/go-binance"

type Client interface {
	AccountBalance() ([]gobinance.Balance, error)
	ListPrices(symbol string) ([]*gobinance.SymbolPrice, error)
	CreateOrder(req *orderReq) (*gobinance.CreateOrderResponse, error)
	GetOrder(req *orderID) (*gobinance.Order, error)
	CancelOrder(req *orderID) (*gobinance.CancelOrderResponse, error)
	GetKlines(req *klinesReq) ([]*gobinance.Kline, error)
}

type orderReq struct {
	Symbol   string
	Side     gobinance.SideType
	Type     gobinance.OrderType
	Tif      gobinance.TimeInForceType
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

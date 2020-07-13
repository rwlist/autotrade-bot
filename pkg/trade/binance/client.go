package binance

import (
	"context"

	gobinance "github.com/adshao/go-binance"
)

type CliDef struct {
	client *gobinance.Client
}

func NewClientDefault(cli *gobinance.Client) *CliDef {
	return &CliDef{
		client: cli,
	}
}

func (b *CliDef) AccountBalance() ([]gobinance.Balance, error) {
	info, err := b.client.NewGetAccountService().Do(context.Background())

	if err != nil {
		return nil, err
	}

	return info.Balances, err
}

func (b *CliDef) ListPrices(symbol string) ([]*gobinance.SymbolPrice, error) {
	list, err := b.client.NewListPricesService().Symbol(symbol).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *CliDef) CreateOrder(req *orderReq) (*gobinance.CreateOrderResponse, error) {
	order, err := b.client.NewCreateOrderService().
		Symbol(req.Symbol).
		Side(req.Side).
		Type(req.Type).
		TimeInForce(gobinance.TimeInForceTypeGTC).
		Price(req.Price).
		Quantity(req.Quantity).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return order, err
}

func (b *CliDef) GetOrder(req *orderID) (*gobinance.Order, error) {
	order, err := b.client.NewGetOrderService().
		Symbol(req.Symbol).
		OrderID(req.ID).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return order, err
}

func (b *CliDef) CancelOrder(req *orderID) (*gobinance.CancelOrderResponse, error) {
	res, err := b.client.NewCancelOrderService().
		Symbol(req.Symbol).
		OrderID(req.ID).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return res, err
}

func (b *CliDef) GetKlines(req *klinesReq) ([]*gobinance.Kline, error) {
	klines, err := b.client.NewKlinesService().
		Symbol(req.Symbol).
		Interval(req.T).
		StartTime(req.StartTime).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return klines, err
}

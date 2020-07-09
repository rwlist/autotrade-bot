package binance

import (
	"context"
	"time"

	goBinance "github.com/adshao/go-binance"
	log "github.com/sirupsen/logrus"
)

type CliLog struct {
	client *goBinance.Client
}

func (b *CliLog) log(method string, req, resp interface{}, err error, start time.Time) {
	logger := log.WithFields(log.Fields{
		"method":   method,
		"req":      req,
		"resp":     resp,
		"duration": time.Since(start).String(),
	})

	if err != nil {
		logger = logger.WithError(err)
	}

	logger.Debug("binance request finished")
}

func NewClientLog(apiKey, secretKey string) *CliLog {
	return &CliLog{
		client: goBinance.NewClient(apiKey, secretKey),
	}
}

func (b *CliLog) AccountBalance() ([]goBinance.Balance, error) {
	start := time.Now()

	info, err := b.client.NewGetAccountService().Do(context.Background())
	b.log("AccountBalance", nil, info, err, start)

	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

func (b *CliLog) ListPrices(symbol string) ([]*goBinance.SymbolPrice, error) {
	start := time.Now()

	list, err := b.client.NewListPricesService().Symbol(symbol).Do(context.Background())
	b.log("ListPrices", symbol, list, err, start)

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *CliLog) CreateOrder(req *orderReq) (*goBinance.CreateOrderResponse, error) {
	start := time.Now()

	order, err := b.client.NewCreateOrderService().
		Symbol(req.Symbol).
		Side(req.Side).
		Type(req.Type).
		TimeInForce(goBinance.TimeInForceTypeGTC).
		Price(req.Price).
		Quantity(req.Quantity).Do(context.Background())
	b.log("CreateOrder", req, order, err, start)

	if err != nil {
		return nil, err
	}

	return order, err
}

func (b *CliLog) GetOrder(req *orderID) (*goBinance.Order, error) {
	start := time.Now()

	order, err := b.client.NewGetOrderService().
		Symbol(req.Symbol).
		OrderID(req.ID).Do(context.Background())
	b.log("GetOrder", req, order, err, start)

	if err != nil {
		return nil, err
	}

	return order, err
}

func (b *CliLog) CancelOrder(req *orderID) (*goBinance.CancelOrderResponse, error) {
	start := time.Now()

	res, err := b.client.NewCancelOrderService().
		Symbol(req.Symbol).
		OrderID(req.ID).Do(context.Background())
	b.log("CancelOrder", req, res, err, start)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (b *CliLog) GetKlines(req *klinesReq) ([]*goBinance.Kline, error) {
	start := time.Now()

	klines, err := b.client.NewKlinesService().
		Symbol(req.Symbol).
		Interval(req.T).
		StartTime(req.StartTime).Do(context.Background())
	b.log("GetKlines", req, klines, err, start)

	if err != nil {
		return nil, err
	}

	return klines, err
}

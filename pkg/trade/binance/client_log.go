package binance

import (
	"time"

	gobinance "github.com/adshao/go-binance"
	log "github.com/sirupsen/logrus"
)

type CliLog struct {
	client Client
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

func NewClientLog(cli Client) *CliLog {
	return &CliLog{
		client: cli,
	}
}

func (b *CliLog) AccountBalance() ([]gobinance.Balance, error) {
	start := time.Now()

	info, err := b.client.AccountBalance()
	b.log("AccountBalance", nil, info != nil, err, start)

	if err != nil {
		return nil, err
	}

	return info, nil
}

func (b *CliLog) ListPrices(symbol string) ([]*gobinance.SymbolPrice, error) {
	start := time.Now()

	list, err := b.client.ListPrices(symbol)
	b.log("ListPrices", symbol, list, err, start)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *CliLog) CreateOrder(req *orderReq) (*gobinance.CreateOrderResponse, error) {
	start := time.Now()

	order, err := b.client.CreateOrder(req)
	b.log("CreateOrder", req, order, err, start)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (b *CliLog) GetOrder(req *orderID) (*gobinance.Order, error) {
	start := time.Now()

	order, err := b.client.GetOrder(req)
	b.log("GetOrder", req, order, err, start)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (b *CliLog) CancelOrder(req *orderID) (*gobinance.CancelOrderResponse, error) {
	start := time.Now()

	res, err := b.client.CancelOrder(req)
	b.log("CancelOrder", req, res, err, start)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *CliLog) GetKlines(req *klinesReq) ([]*gobinance.Kline, error) {
	start := time.Now()

	klines, err := b.client.GetKlines(req)
	b.log("GetKlines", req, klines, err, start)

	if err != nil {
		return nil, err
	}

	return klines, nil
}

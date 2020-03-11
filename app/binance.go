package app

import (
	"context"
	"github.com/adshao/go-binance"
	"strings"
)

type MyBinance struct {
	client *binance.Client
}

func NewMyBinance(c *binance.Client) *MyBinance {
	return &MyBinance{client : c}
}

func (b *MyBinance) AccountBalance() ([]binance.Balance, error) {
	info, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

func (b *MyBinance) BalanceToUSD(bal *binance.Balance) (float64, error) {
	haveFree := strToFloat64(bal.Free)
	haveLocked := strToFloat64(bal.Locked)
	if bal.Asset == "USDT" {
		return haveFree + haveLocked, nil
	}

	symbolPrice, err := b.client.NewListPricesService().Symbol(bal.Asset + "USDT").Do(context.Background())
	if err != nil {
		return 0, err
	}
	price := strToFloat64(symbolPrice[0].Price)
	haveFree *= price
	haveLocked *= price
	return haveFree + haveLocked, nil
}

func (b *MyBinance) GetRate() (string, error) {
	symbolPrice, err := b.client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		return "", err
	}
	return symbolPrice[0].Price, nil
}

func isEmptyBalance(str string) bool {
	return strings.Trim(str, ".0") == ""
}

package app

import (
	"context"
	"github.com/adshao/go-binance"
)

func binanceAccountBalance(client *binance.Client) ([]binance.Balance, error) {
	info, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

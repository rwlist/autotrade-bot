package app

import (
	binance "github.com/eranyanay/binance-api"
)

func binanceAccountBalance(client *binance.BinanceClient) ([]*binance.Balance, error) {
	info, err := client.Account()
	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

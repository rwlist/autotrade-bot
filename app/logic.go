package app

import (
	"github.com/eranyanay/binance-api"
)

type Logic struct {
	client *binance.BinanceClient
}

func NewLogic(client *binance.BinanceClient) *Logic {
	return &Logic{
		client: client,
	}
}

type Status struct {
	rate     string
	balances []*binance.Balance
}

func (l *Logic) CommandStatus() (*Status, error) {
	rate, err := binanceRateQuery()
	if err != nil {
		return nil, err
	}
	allBalances, err := binanceAccountBalance(l.client)
	if err != nil {
		return nil, err
	}

	var balances []*binance.Balance
	for _, bal := range allBalances {
		if isEmptyBalance(bal.Free) && isEmptyBalance(bal.Locked) {
			continue
		}

		balances = append(balances, bal)
	}

	res := &Status{
		rate:     rate,
		balances: balances,
	}
	return res, err
}

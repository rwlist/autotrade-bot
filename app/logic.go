package app

import (
	"github.com/adshao/go-binance"
)

type Logic struct {
	client *binance.Client
}

func NewLogic(client *binance.Client) *Logic {
	return &Logic{
		client: client,
	}
}

type Status struct {
	total	 string
	rate     string
	balances []binance.Balance
}

func (l *Logic) CommandStatus() (*Status, error) {
	rate, err := binanceRateQuery(l.client)
	if err != nil {
		return nil, err
	}
	allBalances, err := binanceAccountBalance(l.client)
	if err != nil {
		return nil, err
	}

	var balances []binance.Balance
	var total float64
	for _, bal := range allBalances {
		if isEmptyBalance(bal.Free) && isEmptyBalance(bal.Locked) {
			continue
		}

		balUSD, err := balanceToUSD(l.client, &bal)
		if err != nil {
			return &Status{}, err
		}
		total += balUSD

		balances = append(balances, bal)
	}

	res := &Status{
		total:	  fToStr(total),
		rate:     rate,
		balances: balances,
	}
	return res, err
}

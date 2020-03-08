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

type Balance struct {
	usd    string
	asset  string
	free   string
	locked string
}

type Status struct {
	total	 string
	rate     string
	balances []*Balance
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

	var balances []*Balance
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
		resBal := &Balance{
			   usd:    float64ToStr(balUSD),
			   asset:  bal.Asset,
			   free:   bal.Free,
			   locked: bal.Locked,
		}
		balances = append(balances, resBal)
	}

	res := &Status{
		total:	  float64ToStr(total),
		rate:     rate,
		balances: balances,
	}
	return res, err
}

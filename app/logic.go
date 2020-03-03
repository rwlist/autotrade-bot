package app

import "github.com/eranyanay/binance-api"

const (
	prefix = "ovpn_"

	data = "data"
	tcp = "tcp"
	udp = "udp"
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
	rate string
	balances []*binance.Balance
}

func (l *Logic) CommandStatus() (*Status, error) {
	rate, err := binanceRateQuery()
	balances, err := binanceAccountBalance(l.client)
	res := &Status {
		rate: rate,
		balances: balances,
	}
	return res, err
}
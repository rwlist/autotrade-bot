package stat

import (
	"github.com/shopspring/decimal"

	"github.com/rwlist/autotrade-bot/pkg/trade"
)

type exchangeInfo interface {
	AccountBalance() ([]trade.Balance, error)
	BalanceToUSD(bal *trade.Balance) (decimal.Decimal, error)
	GetRate(symbol ...string) (decimal.Decimal, error)
}

type Service struct {
	info exchangeInfo
}

func New(info exchangeInfo) *Service {
	return &Service{
		info: info,
	}
}

func (s *Service) Status() (*Status, error) {
	rate, err := s.info.GetRate()
	if err != nil {
		return nil, err
	}

	allBalances, err := s.info.AccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []Balance
	total := decimal.Zero
	for _, bal := range allBalances {
		bal := bal
		asset := bal.Asset
		free := bal.Free
		locked := bal.Locked

		if free.Equal(decimal.Zero) && locked.Equal(decimal.Zero) {
			continue
		}

		balanceInUSD, err := s.info.BalanceToUSD(&bal)
		if err != nil {
			return &Status{}, err
		}

		total = total.Add(balanceInUSD)

		balances = append(
			balances,
			Balance{
				USD:    balanceInUSD,
				Asset:  asset,
				Free:   free,
				Locked: locked,
			},
		)
	}

	return &Status{
		Total:    total,
		Rate:     rate,
		Balances: balances,
	}, nil
}

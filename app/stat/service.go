package stat

import (
	"github.com/rwlist/autotrade-bot/trade/binance"
	"github.com/shopspring/decimal"
)

type Service struct {
	myBinance *binance.Binance
}

func New(myBinance *binance.Binance) *Service {
	return &Service{
		myBinance: myBinance,
	}
}

func (s *Service) Status() (*Status, error) {
	rateStr, err := s.myBinance.GetRate()
	if err != nil {
		return nil, err
	}

	rate, err := decimal.NewFromString(rateStr)
	if err != nil {
		return nil, err
	}

	allBalances, err := s.myBinance.AccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []Balance
	total := decimal.Zero
	for _, bal := range allBalances {
		bal := bal
		asset := unsafeDecimal(bal.Asset)
		free := unsafeDecimal(bal.Free)
		locked := unsafeDecimal(bal.Locked)

		if free.Equal(decimal.Zero) || locked.Equal(decimal.Zero) {
			continue
		}

		usdFloat, err := s.myBinance.BalanceToUSD(&bal)
		if err != nil {
			return &Status{}, err
		}

		balanceInUSD := decimal.NewFromFloat(usdFloat)
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

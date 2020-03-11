package app

type Logic struct {
	b *MyBinance
}

func NewLogic(b *MyBinance) *Logic {
	return &Logic{
		b: b,
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
	rate, err := l.b.GetRate()
	if err != nil {
		return nil, err
	}
	allBalances, err := l.b.AccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []*Balance
	var total float64
	for _, bal := range allBalances {
		if isEmptyBalance(bal.Free) && isEmptyBalance(bal.Locked) {
			continue
		}

		balUSD, err := l.b.BalanceToUSD(&bal)
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

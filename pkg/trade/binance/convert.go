package binance

import (
	goBinance "github.com/adshao/go-binance"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

func convertBalance(bal *goBinance.Balance) trade.Balance {
	return trade.Balance{
		Asset:  bal.Asset,
		Free:   bal.Free,
		Locked: bal.Locked,
	}
}

func convertBalanceSlice(bal []goBinance.Balance) []trade.Balance {
	var newBal []trade.Balance
	for _, val := range bal {
		val := val
		newBal = append(newBal, convertBalance(&val))
	}
	return newBal
}

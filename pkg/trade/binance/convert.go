package binance

import (
	gobinance "github.com/adshao/go-binance"
	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

func convertBalance(bal *gobinance.Balance) trade.Balance {
	return trade.Balance{
		Asset:  bal.Asset,
		Free:   convert.UnsafeDecimal(bal.Free),
		Locked: convert.UnsafeDecimal(bal.Locked),
	}
}

func convertBalanceSlice(bal []gobinance.Balance) []trade.Balance {
	var newBal []trade.Balance
	for _, val := range bal {
		val := val
		newBal = append(newBal, convertBalance(&val))
	}
	return newBal
}

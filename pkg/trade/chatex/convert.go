package chatex

import (
	"strings"

	chatexsdk "github.com/chatex-com/sdk-go"

	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

func convertBalance(b chatexsdk.Balance) trade.Balance {
	return trade.Balance{
		Asset:  strings.ToUpper(b.Coin),
		Free:   convert.UnsafeDecimal(b.Amount),
		Locked: convert.UnsafeDecimal(b.Held),
	}
}

func convertBalanceSlice(arr []chatexsdk.Balance) []trade.Balance {
	var res []trade.Balance
	for _, v := range arr {
		res = append(res, convertBalance(v))
	}

	return res
}

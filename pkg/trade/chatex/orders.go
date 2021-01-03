package chatex

import (
	"time"

	chatexsdk "github.com/chatex-com/sdk-go"
)

type FetchedOrders struct {
	Timestamp time.Time
	Orders    []chatexsdk.Order
	Pair      string
}

type OrdersSnapshot struct {
	Fetched  map[string]FetchedOrders
	Coins    []chatexsdk.Coin
	Started  time.Time
	Finished time.Time

	// single pair only, fetched now
	IsMomentSnapshot bool
}

func (s OrdersSnapshot) CoinCodes() []string {
	var coins []string
	for _, coin := range s.Coins {
		coins = append(coins, coin.Code)
	}
	return coins
}

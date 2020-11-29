package chatex

import (
	"time"

	chatexsdk "github.com/chatex-com/sdk-go"
)

type FetchedOrders struct {
	Timestamp time.Time
	Orders    []chatexsdk.Order
}

type OrdersSnapshot struct {
	Fetched  map[string]FetchedOrders
	Started  time.Time
	Finished time.Time
}

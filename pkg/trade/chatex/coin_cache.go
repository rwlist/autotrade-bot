package chatex

import (
	"sync"

	chatexsdk "github.com/chatex-com/sdk-go"
)

type CoinCache struct {
	coins map[string]chatexsdk.Coin
	m     sync.RWMutex
}

func NewCoinCache() *CoinCache {
	return &CoinCache{
		coins: make(map[string]chatexsdk.Coin),
	}
}

func (c *CoinCache) Update(coin chatexsdk.Coin) {
	c.m.Lock()
	defer c.m.Unlock()

	c.coins[coin.Code] = coin
}

func (c *CoinCache) Get(code string) chatexsdk.Coin {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.coins[code]
}

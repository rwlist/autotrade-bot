package chatex

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	chatexsdk "github.com/chatex-com/sdk-go"

	"github.com/rwlist/autotrade-bot/pkg/metrics"
	"github.com/rwlist/autotrade-bot/pkg/store/redisdb"
)

type callback func(OrdersSnapshot)

type OrdersCollector struct {
	cli   *chatexsdk.Client
	list  *redisdb.List
	log   logrus.FieldLogger
	opts  *TradeOpts
	coins *CoinCache

	callbacks []callback
	mu        sync.RWMutex
}

func NewOrdersCollector(cli *chatexsdk.Client, list *redisdb.List, opts *TradeOpts) *OrdersCollector {
	return &OrdersCollector{
		cli:   cli,
		list:  list,
		log:   logrus.StandardLogger(),
		opts:  opts,
		coins: NewCoinCache(),
	}
}

func (c *OrdersCollector) getEnabledCoins() ([]chatexsdk.Coin, error) {
	tmp, err := c.cli.GetCoins(context.Background())
	if err != nil {
		return nil, err
	}

	opts, _ := c.opts.GetAll()

	var coins []chatexsdk.Coin
	for _, coin := range tmp {
		c.coins.Update(coin)

		code := coin.Code
		if opts.Get("coins."+code+".disabled") == "true" {
			continue
		}

		coins = append(coins, coin)
	}

	return coins, nil
}

func (c *OrdersCollector) Coin(code string) chatexsdk.Coin {
	return c.coins.Get(code)
}

func (c *OrdersCollector) CollectAll() (*OrdersSnapshot, error) {
	const (
		defaultSleep = time.Second / 2
	)

	started := time.Now()

	coins, err := c.getEnabledCoins()
	if err != nil {
		return nil, err
	}

	n := len(coins)
	pairsCount := n * (n - 1)

	c.log.WithFields(logrus.Fields{
		"coins": n,
		"pairs": pairsCount,
	}).Info("selected all coins")

	result := make(map[string]FetchedOrders)

	for _, coin1 := range coins {
		for _, coin2 := range coins {
			if coin1.Code >= coin2.Code {
				continue
			}

			time.Sleep(defaultSleep)
			fetched1, err := c.FetchOrders(coin1.Code, coin2.Code)
			if err != nil {
				c.log.WithError(err).WithField("pair", fetched1.Pair).Error("failed to get orders")
				return nil, err
			}
			result[fetched1.Pair] = fetched1

			time.Sleep(defaultSleep)
			fetched2, err := c.FetchOrders(coin2.Code, coin1.Code)
			if err != nil {
				c.log.WithError(err).WithField("pair", fetched2.Pair).Error("failed to get orders")
				return nil, err
			}
			result[fetched2.Pair] = fetched2

			momentSnapshot := OrdersSnapshot{
				Fetched: map[string]FetchedOrders{
					fetched1.Pair: fetched1,
					fetched2.Pair: fetched2,
				},
				Coins:            []chatexsdk.Coin{coin1, coin2},
				IsMomentSnapshot: true,
			}

			go c.sendCallbacks(momentSnapshot)
		}
	}

	finished := time.Now()

	return &OrdersSnapshot{
		Fetched:  result,
		Coins:    coins,
		Started:  started,
		Finished: finished,
	}, nil
}

func (c *OrdersCollector) CollectInf(ctx context.Context) error {
	for {
		c.collectAndSave()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.opts.FetchCollectorPeriod()):
			continue
		}
	}
}

func (c *OrdersCollector) collectAndSave() {
	if val, _ := c.opts.GetSingle("chatex.collector.state"); val == "disable" {
		c.log.Info("skipping collectAndSave due to disable config")
		return
	}

	snapshot, err := c.CollectAll()
	if err != nil {
		c.log.WithError(err).Error("failed to collect all")
		metrics.ChatexCollectorErr()
		return
	}

	err = c.list.LPush(snapshot)
	if err != nil {
		c.log.WithError(err).Error("failed to save orders snapshot")
		metrics.ChatexCollectorErr()
		return
	}

	// clean up old records
	err = c.list.LTrim(0, 1000)
	if err != nil {
		c.log.WithError(err).Error("failed to trim snapshot records")
		// no return
	}

	metrics.ChatexCollectorOk()

	c.sendCallbacks(*snapshot)
}

func (c *OrdersCollector) sendCallbacks(snapshot OrdersSnapshot) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, cb := range c.callbacks {
		cb(snapshot)
	}
}

func (c *OrdersCollector) RegisterCallback(cb callback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.callbacks = append(c.callbacks, cb)
}

func (c *OrdersCollector) Last() (*OrdersSnapshot, error) {
	var snapshot OrdersSnapshot
	err := c.list.Left(&snapshot)
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}

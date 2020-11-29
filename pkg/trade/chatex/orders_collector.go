package chatex

import (
	"context"
	"time"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/sirupsen/logrus"

	"github.com/rwlist/autotrade-bot/pkg/metrics"
	"github.com/rwlist/autotrade-bot/pkg/store/redisdb"
)

type OrdersCollector struct {
	cli  *chatexsdk.Client
	list *redisdb.List
	log  logrus.FieldLogger
}

func NewOrdersCollector(cli *chatexsdk.Client, list *redisdb.List) *OrdersCollector {
	return &OrdersCollector{
		cli:  cli,
		list: list,
		log:  logrus.StandardLogger(),
	}
}

func (c *OrdersCollector) CollectAll() (map[string]FetchedOrders, error) {
	const (
		defaultSleep = time.Second / 2
		sleepOnError = time.Second
	)

	coins, err := c.cli.GetCoins(context.Background())
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
			if coin1.Code == coin2.Code {
				continue
			}

			time.Sleep(defaultSleep)

			pair := coin1.Code + "/" + coin2.Code

			const (
				offset = 0
				limit  = 50
			)

			var (
				now    time.Time
				orders []chatexsdk.Order
				err    error
			)

			for retries := 0; retries < 3; retries++ {
				now = time.Now()
				orders, err = c.cli.GetOrders(context.Background(), pair, offset, limit)
				if err == chatexsdk.ErrTooManyRequests {
					time.Sleep(sleepOnError)
					continue
				}
				break
			}
			if err != nil {
				c.log.WithError(err).WithField("pair", pair).Error("failed to get orders")
				return nil, err
			}

			timeAfter := time.Now()

			fetched := FetchedOrders{
				Timestamp: now,
				Orders:    orders,
			}

			c.log.WithFields(logrus.Fields{
				"duration": timeAfter.Sub(now),
				"fetched":  len(fetched.Orders),
			}).Info("fetched orders by pair")

			result[pair] = fetched
		}
	}

	return result, nil
}

func (c *OrdersCollector) CollectInf(ctx context.Context) error {
	const every = time.Minute

	for {
		c.collectAndSave()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(every):
			continue
		}
	}
}

func (c *OrdersCollector) collectAndSave() {
	started := time.Now()

	all, err := c.CollectAll()
	if err != nil {
		c.log.WithError(err).Error("failed to collect all")
		metrics.ChatexCollectorErr()
		return
	}

	finished := time.Now()

	snapshot := OrdersSnapshot{
		Fetched:  all,
		Started:  started,
		Finished: finished,
	}

	err = c.list.LPush(snapshot)
	if err != nil {
		c.log.WithError(err).Error("failed to save orders snapshot")
		metrics.ChatexCollectorErr()
		return
	}

	metrics.ChatexCollectorOk()
}

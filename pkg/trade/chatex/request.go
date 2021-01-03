package chatex

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	chatexsdk "github.com/chatex-com/sdk-go"
)

func (c *OrdersCollector) FetchOrders(code1, code2 string) (FetchedOrders, error) {
	const sleepOnError = time.Second

	pair := code1 + "/" + code2

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
		return FetchedOrders{
			Pair: pair,
		}, err
	}

	timeAfter := time.Now()

	fetched := FetchedOrders{
		Timestamp: now,
		Orders:    orders,
		Pair:      pair,
	}

	c.log.WithFields(logrus.Fields{
		"duration": timeAfter.Sub(now),
		"fetched":  len(fetched.Orders),
	}).Info("fetched orders by pair")

	return fetched, nil
}

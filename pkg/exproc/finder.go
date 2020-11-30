package exproc

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"

	chatexsdk "github.com/chatex-com/sdk-go"

	"github.com/rwlist/autotrade-bot/pkg/money"
	"github.com/rwlist/autotrade-bot/pkg/trade/chatex"
)

type Sender interface {
	Send(text string)
}

type Finder struct {
	cli       *chatexsdk.Client
	collector *chatex.OrdersCollector
	sender    Sender
}

func NewFinder(cli *chatexsdk.Client, collector *chatex.OrdersCollector, sender Sender) *Finder {
	return &Finder{
		cli:       cli,
		collector: collector,
		sender:    sender,
	}
}

func (f *Finder) OnSnapshot(snap chatex.OrdersSnapshot) {
	const logAllPaths = false

	spew.Dump(snap)

	var coins []string
	for _, coin := range snap.Coins {
		coins = append(coins, coin.Code)
	}

	g := make(map[string]map[string]decimal.Decimal)

	for _, coin := range coins {
		g[coin] = make(map[string]decimal.Decimal)
	}

	for k, v := range snap.Fetched {
		tmp := strings.Split(k, "/")
		if len(tmp) != 2 {
			log.WithField("k", k).Error("failed to parse pair")
			continue
		}

		if len(v.Orders) == 0 {
			// no orders
			continue
		}

		order := v.Orders[0]

		reverseRate := order.Rate

		if !reverseRate.IsPositive() {
			log.WithField("k", k).WithField("reverseRate", reverseRate).Error("invalid rate")
			continue
		}

		rate := money.One.DivRound(reverseRate, money.Precision)

		coin1 := tmp[1]
		coin2 := tmp[0]

		g[coin1][coin2] = rate
	}

	var results []string

	for _, start := range coins {
		for _, next := range coins {
			have := money.One

			go1, ok1 := g[start][next]
			go2, ok2 := g[next][start]

			if !ok1 || !ok2 {
				continue
			}

			have = have.Mul(go1).Mul(go2)

			if logAllPaths {
				log.WithFields(log.Fields{
					"start": start,
					"next":  next,
					"have":  have,
				}).Info("check loop")
			}

			if !have.GreaterThan(money.One) {
				// boring
				continue
			}

			places := int32(5)

			order1 := snap.Fetched[next+"/"+start].Orders[0]
			order2 := snap.Fetched[start+"/"+next].Orders[0]

			info := fmt.Sprintf(
				"Found: %s %s → %s %s → %s %s.",
				money.One.Round(places),
				start,
				go1.Round(places),
				next,
				have.Round(places),
				start,
			)

			// next -> start
			midAmount1 := order1.Amount

			// next <- last_start
			midAmount2 := order2.Amount.Mul(order2.Rate)

			// take min
			midAmount := decimal.Min(midAmount1, midAmount2)

			// start
			startAmount := midAmount.Mul(order1.Rate)

			// apply order1
			nextAmount := startAmount.DivRound(order1.Rate, money.Precision)

			buy1 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				nextAmount.Round(places),
				next,
				startAmount.Round(places),
				start,
				order1.ID,
			)

			// apply order2
			lastAmount := nextAmount.DivRound(order2.Rate, money.Precision)

			buy2 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				lastAmount.Round(places),
				start,
				nextAmount.Round(places),
				next,
				order2.ID,
			)

			log.WithFields(log.Fields{
				"order1":      order1,
				"order2":      order2,
				"start":       start,
				"next":        next,
				"info":        info,
				"buy1":        buy1,
				"buy2":        buy2,
				"midAmount1":  midAmount1,
				"midAmount2":  midAmount2,
				"midAmount":   midAmount,
				"startAmount": startAmount,
				"nextAmount":  nextAmount,
				"lastAmount":  lastAmount,
			}).Info("found positive loop")

			info = info + "\n* " + buy1 + "\n* " + buy2
			results = append(results, info)
		}
	}

	if len(results) == 0 {
		// boring
		return
	}

	f.sender.Send(strings.Join(results, "\n\n"))
}

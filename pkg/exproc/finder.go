package exproc

import (
	"fmt"
	"strings"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"

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

		coin1 := tmp[0]
		coin2 := tmp[1]

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

			info := fmt.Sprintf(
				"Found: %s %s → %s %s → %s %s.\tOrders %v and %v",
				money.One.Round(places),
				start,
				go1.Round(places),
				next,
				have,
				start,
				snap.Fetched[next+"/"+start].Orders[0].ID,
				snap.Fetched[start+"/"+next].Orders[0].ID,
			)

			results = append(results, info)
		}
	}

	if len(results) == 0 {
		// boring
		return
	}

	f.sender.Send(strings.Join(results, "\n"))
}

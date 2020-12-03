package exproc

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
	tradeOpts *chatex.TradeOpts

	tradeMutex sync.Mutex
}

func NewFinder(cli *chatexsdk.Client, collector *chatex.OrdersCollector, tradeOpts *chatex.TradeOpts, sender Sender) *Finder {
	return &Finder{
		cli:       cli,
		collector: collector,
		tradeOpts: tradeOpts,
		sender:    sender,
	}
}

func (f *Finder) OnSnapshot(snap chatex.OrdersSnapshot) { //nolint:funlen
	const logAllPaths = false

	log.Info("processing chatex snapshot in finder")

	coins := snap.CoinCodes()

	g := make(map[string]map[string]decimal.Decimal)

	for _, coin := range coins {
		g[coin] = make(map[string]decimal.Decimal)
	}

	for k, v := range snap.Fetched {
		tmp := strings.Split(k, "/")
		const args = 2
		if len(tmp) != args {
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

			const places = 5

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

			calc := OrderCalc{}.CalcTrades(order1, order2)

			buy1 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				calc.NextAmount.Round(places),
				next,
				calc.StartAmount.Round(places),
				start,
				order1.ID,
			)

			buy2 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				calc.LastAmount.Round(places),
				start,
				calc.NextAmount.Round(places),
				next,
				order2.ID,
			)

			log.WithFields(log.Fields{
				"order1": order1,
				"order2": order2,
				"start":  start,
				"next":   next,
				"info":   info,
				"buy1":   buy1,
				"buy2":   buy2,
				"calc":   calc,
			}).Info("found positive loop")

			info = info + "\n* " + buy1 + "\n* " + buy2
			results = append(results, info)

			go f.makeTrades(order1, order2)
		}
	}

	if len(results) == 0 {
		// boring
		return
	}

	f.sender.Send(strings.Join(results, "\n\n"))
}

func (f *Finder) makeTrades(order1, order2 chatexsdk.Order) {
	// trades must not be clashed
	f.tradeMutex.Lock()
	defer f.tradeMutex.Unlock()

	logger := log.WithField("order1", order1).WithField("order2", order2)

	err := f.refreshOrder(&order1)
	if err != nil {
		logger.WithError(err).Error("failed to refresh order1")
		return
	}

	err = f.refreshOrder(&order2)
	if err != nil {
		logger.WithError(err).Error("failed to refresh order2")
		return
	}

	// sleep some time to relax rate limits
	const relaxTime = time.Second / 2
	time.Sleep(relaxTime)
}

func (f *Finder) refreshOrder(ptr *chatexsdk.Order) error {
	const relaxTime = time.Second / 2

	var (
		res *chatexsdk.Order
		err error
	)

	for i := 0; i < 3; i++ {
		res, err = f.cli.GetOrder(context.Background(), uint(ptr.ID))
		if err == chatexsdk.ErrTooManyRequests {
			time.Sleep(relaxTime)
			continue
		}
	}

	if err != nil {
		return err
	}

	*ptr = *res
	return nil
}

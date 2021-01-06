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

	if !snap.IsMomentSnapshot {
		log.Info("ignoring non-moment snapshot")
		return
	}

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

			calc := OrderCalc{
				StartCoin: f.collector.Coin(start),
				NextCoin:  f.collector.Coin(next),
			}.CalcTrades(order1, order2)

			buy1 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				calc.NextAmount.Round(places),
				next,
				calc.StartAmount.Round(places),
				start,
				chatex.OrderLinkMd(order1.ID),
			)

			buy2 := fmt.Sprintf(
				"Buy %s %s for %s %s, orderID = %v",
				calc.LastAmount.Round(places),
				start,
				calc.NextAmount.Round(places),
				next,
				chatex.OrderLinkMd(order2.ID),
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

			go f.makeTrades(snap, order1, order2)
		}
	}

	if len(results) == 0 {
		// boring
		return
	}

	f.sender.Send(strings.Join(results, "\n\n"))
}

func (f *Finder) makeTrades(snap chatex.OrdersSnapshot, order1, order2 chatexsdk.Order) { //nolint:funlen
	const places = 8

	// trades must not be clashed
	f.tradeMutex.Lock()
	defer f.tradeMutex.Unlock()

	logger := log.WithField("order1", order1).WithField("order2", order2)

	if snap.IsMomentSnapshot {
		logger.Info("info from moment snapshot")
		// TODO: invent some way to skip refresh in some cases
	}

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

	log.WithField("order1", order1).WithField("order2", order2).Info("updated orders")

	// sleep some time to relax rate limits
	const relaxTime = time.Second / 2
	time.Sleep(relaxTime)

	opts, err := f.tradeOpts.GetAll()
	if err != nil {
		log.WithError(err).Error("failed to get opts")
		return
	}

	pair1 := pairOf(order1)
	limitOptionName := "limit." + pair1.Sell
	myLimit := opts.Decimal(limitOptionName)

	calc := OrderCalc{
		MaxStartAmount: &myLimit,
		StartCoin:      f.collector.Coin(pair1.Sell),
		NextCoin:       f.collector.Coin(pair1.Buy),
	}.CalcTrades(order1, order2)

	var info []string
	info = append(
		info,
		fmt.Sprintf(
			"Attempt to do trade: %s -> %s -> %s",
			pair1.Sell,
			pair1.Buy,
			pair1.Sell,
		),
		fmt.Sprintf(
			"order1 = %v, order2 = %v",
			chatex.OrderLinkMd(order1.ID),
			chatex.OrderLinkMd(order2.ID),
		),
		fmt.Sprintf(
			"amount1 = %v, amount2 = %v",
			order1.Amount,
			order2.Amount,
		),
		fmt.Sprintf(
			"startAmount = %v, nextAmount = %v, lastAmount = %v",
			calc.StartAmount.RoundBank(places),
			calc.NextAmount.RoundBank(places),
			calc.LastAmount.RoundBank(places),
		),
	)

	if order1.Status != chatexsdk.Active || order2.Status != chatexsdk.Active {
		info = append(info, "Error: one of orders is not active")
		f.sender.Send(strings.Join(info, "\n"))
		return
	}

	if !calc.StartAmount.IsPositive() || !calc.NextAmount.IsPositive() || !calc.LastAmount.IsPositive() {
		info = append(info, "Error: not positive amount.")
		f.sender.Send(strings.Join(info, "\n"))

		if f.tryFixNonPositive(order1, order2) {
			// looks like successful trade, retrying
			f.retryMakeTrades(snap, order1, order2)
		}

		return
	}

	if !calc.LastAmount.GreaterThan(calc.StartAmount) {
		info = append(info, "Error: not positive cycle")
		f.sender.Send(strings.Join(info, "\n"))
		return
	}

	trade1, err := f.makeTrade(order1.ID, chatexsdk.TradeRequest{
		Amount: calc.NextAmount,
		Rate:   order1.Rate,
	})
	if err != nil {
		log.WithError(err).Error("failed to make trade1")
		info = append(info, "Error(trade1): "+err.Error())
		f.sender.Send(strings.Join(info, "\n"))
		return
	}
	logger.WithField("trade1", trade1).Info("made trade")

	// trade1 is finished, so myLimit should be decreased
	myLimit = myLimit.Sub(calc.StartAmount)
	err = f.tradeOpts.SetOption(limitOptionName, myLimit.String())
	if err != nil {
		log.WithError(err).Error("failed to update myLimit")
	}

	info = append(
		info,
		fmt.Sprintf(
			"Updated myLimit = %v",
			myLimit.RoundBank(places),
		),
		fmt.Sprintf(
			"trade1 = %v, received = %v, amount = %v",
			trade1.ID,
			trade1.ReceivedAmount,
			trade1.Amount,
		),
	)

	trade2, err := f.makeTrade(order2.ID, chatexsdk.TradeRequest{
		Amount: calc.LastAmount,
		Rate:   order2.Rate,
	})
	if err != nil {
		log.WithError(err).Error("failed to make trade2")
		info = append(info, "Error(trade2): "+err.Error())
		f.sender.Send(strings.Join(info, "\n"))
		return
	}
	logger.WithField("trade2", trade2).Info("made trade")

	info = append(
		info,
		fmt.Sprintf(
			"trade2 = %v, received = %v, amount = %v",
			trade2.ID,
			trade2.ReceivedAmount,
			trade2.Amount,
		),
		"All ok!",
	)

	f.sender.Send(strings.Join(info, "\n"))

	// successful trade, retrying
	f.retryMakeTrades(snap, order1, order2)
}

func (f *Finder) retryMakeTrades(snap chatex.OrdersSnapshot, order1, order2 chatexsdk.Order) {
	const retrySleep = time.Second
	time.Sleep(retrySleep)

	go f.makeTrades(snap, order1, order2)
}

func (f *Finder) refreshOrder(ptr *chatexsdk.Order) error {
	codes := strings.Split(ptr.Pair, "/")
	if len(codes) != 2 { //nolint:gomnd
		return fmt.Errorf("invalid pair: %s", ptr.Pair)
	}

	res, err := f.collector.FetchOrders(codes[0], codes[1])
	if err != nil {
		return err
	}

	if len(res.Orders) == 0 {
		return fmt.Errorf("pair orderbook is empty, pair=%s", ptr.Pair)
	}

	*ptr = res.Orders[0]
	return nil
}

func (f *Finder) makeTrade(orderID uint64, req chatexsdk.TradeRequest) (*chatexsdk.Trade, error) {
	const relaxTime = time.Second / 2

	var (
		res *chatexsdk.Trade
		err error
	)

	for i := 0; i < 3; i++ {
		res, err = f.cli.CreateTrade(context.Background(), uint(orderID), req)
		if err == chatexsdk.ErrTooManyRequests {
			time.Sleep(relaxTime)
			continue
		}
		break
	}

	return res, err
}

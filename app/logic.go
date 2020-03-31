package app

import (
	"github.com/rwlist/autotrade-bot/draw"
	"github.com/rwlist/autotrade-bot/tostr"

	"time"

	"github.com/rwlist/autotrade-bot/binance"
)

type Logic struct {
	b *binance.MyBinance
}

func NewLogic(b *binance.MyBinance) *Logic {
	return &Logic{
		b: b,
	}
}

type Balance struct {
	usd    string
	asset  string
	free   string
	locked string
}

type Status struct {
	total    string
	rate     string
	balances []*Balance
}

func (l *Logic) CommandStatus() (*Status, error) {
	rate, err := l.b.GetRate()
	if err != nil {
		return nil, err
	}
	allBalances, err := l.b.AccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []*Balance
	var total float64
	for _, bal := range allBalances {
		if binance.IsEmptyBalance(bal.Free) && binance.IsEmptyBalance(bal.Locked) {
			continue
		}

		bal := bal
		balUSD, err := l.b.BalanceToUSD(&bal)
		if err != nil {
			return &Status{}, err
		}
		total += balUSD
		resBal := &Balance{
			usd:    tostr.Float64ToStr(balUSD, 2),
			asset:  bal.Asset,
			free:   bal.Free,
			locked: bal.Locked,
		}
		balances = append(balances, resBal)
	}

	res := &Status{
		total:    tostr.Float64ToStr(total, 2),
		rate:     rate,
		balances: balances,
	}
	return res, err
}

const sleepDur = time.Second

func (l *Logic) CommandBuy(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.BuyAll()
		if err != nil {
			s.Send(errorMessage(err, "Buy-BuyAll"))
			return
		}
		s.Send(startMessage(orderNew))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID())
		if err != nil {
			s.Send(errorMessage(err, "Buy-GetOrder"))
			return
		}
		s.Send(orderStatusMessage(order))
		err = l.b.CancelOrder(order.OrderID())
		if err != nil {
			s.Send(errorMessage(err, "Buy-CancelOrder"))
			return
		}
	}
}

func (l *Logic) CommandSell(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.SellAll()
		if err != nil {
			s.Send(errorMessage(err, "Sell-BuyAll"))
			return
		}
		s.Send(startMessage(orderNew))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID())
		if err != nil {
			s.Send(errorMessage(err, "Sell-GetOrder"))
			return
		}
		s.Send(orderStatusMessage(order))
		err = l.b.CancelOrder(order.OrderID())
		if err != nil {
			s.Send(errorMessage(err, "Sell-CancelOrder"))
			return
		}
	}
}

func (l *Logic) CommandDraw(s *Sender) {
	klines, err := l.b.GetKlines()

	if err != nil {
		s.Send(errorMessage(err, "Draw GetKlines"))
		return
	}

	p, err := draw.NewPlot()
	if err != nil {
		s.Send(errorMessage(err, "Draw in plot.New()"))
		return
	}

	p.DrawEnv()
	p.DrawHelpLines(klines.Last, klines.Min, klines.Max, klines.StartTime)
	err = p.DrawMainGraph(klines)
	if err != nil {
		s.Send(errorMessage(err, "Draw in p.DrawMainGraph(klines)"))
		return
	}

	buffer, err := p.SaveToBuffer()
	if err != nil {
		s.Send(errorMessage(err, "Draw in p.SaveToBuffer()"))
		return
	}
	err = s.SendPhoto("graph.png", buffer.Bytes())
	if err != nil {
		s.Send(errorMessage(err, "Draw in s.SendPhoto(\"graph.png\", buffer.Bytes())"))
		return
	}
}

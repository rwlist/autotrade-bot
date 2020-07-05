package logic

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/trigger"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/tostr"

	"github.com/rwlist/autotrade-bot/pkg/binance"
)

type Logic struct {
	b  binance.Binance
	ft *trigger.FormulaTrigger
}

func NewLogic(b binance.Binance, ft *trigger.FormulaTrigger) *Logic {
	return &Logic{
		b:  b,
		ft: ft,
	}
}

const sleepDur = 500 * time.Millisecond

func (l *Logic) Buy(s Sender) error {
	for i := 0; i < 10; i++ {
		order, done, err := l.b.BuyAll()
		if done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("in binance.BuyAll: %w", err)
		}
		txt := startMessage(order) + "\n" + orderStatusMessage(order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			return fmt.Errorf("in binance.CancelOrder: %w", err)
		}
	}
	return nil
}

func (l *Logic) Sell(s Sender) error {
	for i := 0; i < 10; i++ {
		order, done, err := l.b.SellAll()
		if done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("in binance.SellAll: %w", err)
		}
		txt := startMessage(order) + "\n" + orderStatusMessage(order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			return fmt.Errorf("in binance.CancelOrder: %w", err)
		}
	}
	return nil
}

func (l *Logic) Draw(str string, optF formula.Formula) ([]byte, error) {
	klines, err := l.b.GetKlines()
	if err != nil {
		return nil, fmt.Errorf("logic.Draw in binance.GetKlines: %w", err)
	}

	p := draw.NewPlot()

	p.AddEnv()

	p.AddHelpLines(klines.Last, klines.Min, klines.Max, klines.StartTime)

	p.AddRateGraph(klines)

	yMax := klines.Max + math.Sqrt(klines.Max)
	if optF == nil {
		f, err := formula.NewBasic(str, klines.Last, float64(time.Now().Unix()))
		if err != nil {
			return nil, fmt.Errorf("logic.Draw in formula.NewBasic: %w", err)
		}
		p.AddFunction(f, yMax)
	} else {
		p.AddFunction(optF, yMax)
	}
	buffer, err := p.SaveToBuffer()
	if err != nil {
		return nil, fmt.Errorf("logic.Draw in draw.SaveToBuffer: %w", err)
	}
	return buffer.Bytes(), nil
}

func (l *Logic) Begin(s Sender, str string, isTest bool) error {
	if !isTest {
		err := l.Buy(s)
		if err != nil {
			return fmt.Errorf("logic.Begin in logic.Buy: %w", err)
		}
	}
	rate, err := l.b.GetRate()
	if err != nil {
		return fmt.Errorf("logic.Begin in binance.GetRate: %w", err)
	}
	f, err := formula.NewBasic(str, tostr.StrToFloat64(rate), float64(time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("logic.Begin in formula.NewBasic: %w", err)
	}
	l.ft.Begin(f)
	go l.checkLoop(s, isTest)
	return nil
}

func (l *Logic) checkLoop(s Sender, isTest bool) {
	var cnt, period int64 = 0, 30
	f := l.ft.GetFormula()
	for val := range l.ft.Resp {
		if cnt%period == 0 {
			s.Send(triggerResponseMessage(val))
			b, _ := l.Draw("", f)
			s.SendPhoto("graph.png", b)
		}
		if val.AbsDif < 0 {
			s.Send(triggerResponseMessage(val))
			err := l.End(s, isTest)
			if err != nil {
				s.Send(fmt.Sprintf("command end error: %v", err))
				log.Println(err)
			}
			return
		}
		if val.RelDif < 1 {
			period = 6
		} else {
			period = 30
		}
		cnt++
	}
}

func (l *Logic) End(s Sender, isTest bool) error {
	if !isTest {
		err := l.Sell(s)
		if err != nil {
			return fmt.Errorf("logic.End in logic.Sell: %w", err)
		}
	}
	l.ft.End()
	s.Send("trigger OFF")
	return nil
}

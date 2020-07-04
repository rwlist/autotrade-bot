package logic

import (
	"fmt"
	"math"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/trigger"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/tostr"

	"github.com/rwlist/autotrade-bot/pkg/binance"
)

type Logic struct {
	b  *binance.Binance
	ft trigger.FormulaTrigger
}

func NewLogic(b *binance.Binance) *Logic {
	return &Logic{
		b: b,
	}
}

const sleepDur = 600 * time.Millisecond

func (l *Logic) Buy(s Sender) error {
	for i := 0; i < 10; i++ {
		order, done, err := l.b.BuyAll()
		if done {
			return nil
		}
		if err != nil {
			txt := fmt.Sprintf("Order can't be placed\nError: %v", err)
			s.Send(txt)
			return err
		}
		txt := startMessage(order) + "\n" + orderStatusMessage(order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			txt := fmt.Sprintf("Can't cancel order:\nId: %v\nError:%v", order.OrderID, err)
			s.Send(txt)
			return err
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
			txt := fmt.Sprintf("Order can't be placed\nError: %v", err)
			s.Send(txt)
			return err
		}
		txt := startMessage(order) + "\n" + orderStatusMessage(order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			txt := fmt.Sprintf("Can't cancel order:\nId: %v\nError:%v", order.OrderID, err)
			s.Send(txt)
			return err
		}
	}
	return nil
}

func (l *Logic) Draw(str string, optF formula.Formula) ([]byte, error) {
	klines, err := l.b.GetKlines()
	if err != nil {
		return nil, fmt.Errorf("logic.Draw in binance.GetKlines: %w", err)
	}

	p, err := draw.NewPlot()
	if err != nil {
		return nil, fmt.Errorf("logic.Draw in draw.NewPlot: %w", err)
	}

	p.AddEnv()

	p.AddHelpLines(klines.Last, klines.Min, klines.Max, klines.StartTime)

	err = p.AddRateGraph(klines)
	if err != nil {
		return nil, fmt.Errorf("logic.Draw in draw.AddRateGraph: %w", err)
	}

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
	var err error
	l.ft, err = trigger.NewTrigger(*l.b)
	if err != nil {
		return fmt.Errorf("logic.Begin in trigger.NewTrigger: %w", err)
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
	var cnt, period int64 = 0, 30
	for val := range l.ft.Resp {
		if cnt%period == 0 {
			s.Send(triggerResponseMessage(val))
			b, err := l.Draw("", f)
			if err != nil {
				return fmt.Errorf("logic.Begin in logic.Draw: %w", err)
			}
			err = s.SendPhoto("graph.png", b)
			if err != nil {
				return fmt.Errorf("logic.Begin in logic.SendPhoto: %w", err)
			}
		}
		if val.AbsDif < 0 {
			s.Send(triggerResponseMessage(val))
			err = l.End(s, isTest)
			if err != nil {
				return fmt.Errorf("logic.Begin in logic.End: %w", err)
			}
			return nil
		}
		if val.RelDif < 1 {
			period = 6
		} else {
			period = 30
		}
		cnt++
	}
	return nil
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

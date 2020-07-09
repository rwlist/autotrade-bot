package logic

import (
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/trigger"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/tostr"
)

type Logic struct {
	b  trade.IExchange
	ft *trigger.FormulaTrigger
}

func NewLogic(b trade.IExchange, ft *trigger.FormulaTrigger) *Logic {
	return &Logic{
		b:  b,
		ft: ft,
	}
}

const sleepDur = 500 * time.Millisecond

func (l *Logic) Buy(s Sender) error {
	for i := 0; i < 10; i++ {
		order := l.b.BuyAll()
		if order.Done {
			return nil
		}
		if order.Err != nil {
			return fmt.Errorf("in binance.BuyAll: %w", order.Err)
		}
		txt := startMessage(order.Order) + "\n" + orderStatusMessage(order.Order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err := l.b.CancelOrder(order.Order.OrderID)
		if err != nil {
			return fmt.Errorf("in binance.CancelOrder: %w", err)
		}
	}
	return nil
}

func (l *Logic) Sell(s Sender) error {
	for i := 0; i < 10; i++ {
		order := l.b.SellAll()
		if order.Done {
			return nil
		}
		if order.Err != nil {
			return fmt.Errorf("in binance.SellAll: %w", order.Err)
		}
		txt := startMessage(order.Order) + "\n" + orderStatusMessage(order.Order)
		s.Send(txt)
		time.Sleep(sleepDur)
		err := l.b.CancelOrder(order.Order.OrderID)
		if err != nil {
			return fmt.Errorf("in binance.CancelOrder: %w", err)
		}
	}
	return nil
}

func (l *Logic) Draw(str string, optF formula.Formula) ([]byte, error) {
	klines, err := l.b.GetKlines()
	if err != nil {
		return nil, fmt.Errorf("in binance.GetKlines: %w", err)
	}

	p := draw.NewPlot()

	p.AddEnv()

	p.AddHelpLines(klines.Last, klines.Min, klines.Max, klines.StartTime)

	p.AddRateGraph(klines)

	yMax := klines.Max + math.Sqrt(klines.Max)
	if optF == nil {
		f, err := formula.NewBasic(str, klines.Last, float64(time.Now().Unix()))
		if err != nil {
			return nil, fmt.Errorf("in formula.NewBasic: %w", err)
		}
		p.AddFunction(f, yMax)
	} else {
		p.AddFunction(optF, yMax)
	}
	buffer, err := p.SaveToBuffer()
	if err != nil {
		return nil, fmt.Errorf("in draw.SaveToBuffer: %w", err)
	}
	return buffer.Bytes(), nil
}

func (l *Logic) Begin(s Sender, str string, isTest bool) error {
	if !isTest {
		err := l.Buy(s)
		if err != nil {
			return fmt.Errorf("in logic.Buy: %w", err)
		}
	}
	rate, err := l.b.GetRate()
	if err != nil {
		return fmt.Errorf("in binance.GetRate: %w", err)
	}
	f, err := formula.NewBasic(str, tostr.StrToFloat64(rate), float64(time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("in formula.NewBasic: %w", err)
	}
	l.ft.Begin(f)
	go l.checkLoop(s, isTest)
	return nil
}

type FormulaStatus struct {
	Txt string
	B   []byte
	Err error
}

/*
	Если формула уже есть, то работает с ней. (Если trigger активен, то считается, что она есть)
	Если формул нет, то парсит из строки и обновляет trigger.
*/
func (l *Logic) Fstat(str string) *FormulaStatus {
	f := l.ft.GetFormula()
	if !l.ft.IsActive() {
		if f == nil {
			rate, err := l.b.GetRate()
			if err != nil {
				return &FormulaStatus{
					Txt: "",
					B:   nil,
					Err: fmt.Errorf("in binance.GetRate: %w: ", err),
				}
			}
			f, err = formula.NewBasic(str, tostr.StrToFloat64(rate), float64(time.Now().Unix()))
			if err != nil {
				return &FormulaStatus{
					Txt: "",
					B:   nil,
					Err: fmt.Errorf("in formula.NewBasic: %w", err),
				}
			}
			l.ft.UpdFormula(f)
		}
		l.ft.UpdResponse()
	}
	resp := l.ft.GetResponse()
	b, err := l.Draw("", f)
	if err != nil {
		return &FormulaStatus{
			Txt: "",
			B:   nil,
			Err: fmt.Errorf("in logic.Draw: %w", err),
		}
	}
	return &FormulaStatus{
		Txt: triggerResponseMessage(&resp),
		B:   b,
		Err: nil,
	}
}

const smallPeriod = int64(time.Minute / trigger.TimeSleep)

func (l *Logic) checkLoop(s Sender, isTest bool) {
	var cnt, period int64 = 0, 30
	f := l.ft.GetFormula()
	for range l.ft.Ping {
		resp := l.ft.GetResponse()
		if cnt%period == 0 {
			s.Send(triggerResponseMessage(&resp))
			b, _ := l.Draw("", f)
			s.SendPhoto("graph.png", b)
		}
		if resp.AbsDif < 0 {
			s.Send(triggerResponseMessage(&resp))
			err := l.End(s, isTest)
			if err != nil {
				s.Send(fmt.Sprintf("command end error: %v", err))
				log.WithError(err).Error("command end error")
			}
			return
		}
		if resp.RelDif < 1 {
			period = smallPeriod
		} else {
			period = 5 * smallPeriod
		}
		cnt++
	}
}

func (l *Logic) End(s Sender, isTest bool) error {
	if !isTest {
		err := l.Sell(s)
		if err != nil {
			return fmt.Errorf("in logic.Sell: %w", err)
		}
	}
	l.ft.End()
	s.Send("trigger OFF")
	return nil
}

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
	b      trade.IExchange
	ft     *trigger.FormulaTrigger
	isTest bool
}

func NewLogic(b trade.IExchange, ft *trigger.FormulaTrigger, isTest bool) *Logic {
	return &Logic{
		b:      b,
		ft:     ft,
		isTest: isTest, // default value is false
	}
}

func (l *Logic) SafeTestModeSwitch() bool {
	if !l.ft.IsActive() {
		l.isTest = !l.isTest
	}
	return l.isTest
}

func (l *Logic) SetScale(scale string) {
	l.b.SetScale(scale)
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

	title := fmt.Sprintf("Scale: %v", klines.Scale)
	p.Plot.Title.Text = title

	p.AddHelpLines(klines.Last, klines.Min, klines.Max, klines.StartTime)

	p.AddRateGraph(klines)

	yMax := klines.Max + math.Log(klines.Max)
	xMax := 2*float64(time.Now().Unix()) - klines.StartTime
	if optF == nil {
		f, err := formula.NewBasic(str, klines.Last, float64(time.Now().Unix()))
		if err != nil {
			return nil, fmt.Errorf("in formula.NewBasic: %w", err)
		}
		p.AddFunction(f, yMax, xMax)
	} else {
		p.AddFunction(optF, yMax, xMax)
	}
	buffer, err := p.SaveToBuffer()
	if err != nil {
		return nil, fmt.Errorf("in draw.SaveToBuffer: %w", err)
	}
	return buffer.Bytes(), nil
}

func (l *Logic) Begin(s Sender, str string) error {
	rate, err := l.b.GetRate()
	if err != nil {
		return fmt.Errorf("in binance.GetRate: %w", err)
	}
	f, err := formula.NewBasic(str, tostr.StrToFloat64(rate), float64(time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("in formula.NewBasic: %w", err)
	}

	isTest := l.isTest
	if !isTest {
		err := l.Buy(s)
		if err != nil {
			return fmt.Errorf("in logic.Buy: %w", err)
		}
	}

	l.ft.Begin(f)
	go l.checkLoop(s)
	return nil
}

type FormulaStatus struct {
	Txt string
	B   []byte
	Err error
}

//	Если формула уже есть, то работает с ней. (Если trigger активен, то считается, что она есть)
//	Если формул нет, то парсит из строки и обновляет trigger.
func (l *Logic) Fstat(str string) *FormulaStatus {
	isTest := l.isTest

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
	info := infoToSend{
		resp:   &resp,
		isTest: isTest,
	}
	return &FormulaStatus{
		Txt: triggerResponseMessage(info),
		B:   b,
		Err: nil,
	}
}

const smallPeriod = int64(time.Minute / trigger.TimeSleep)

func (l *Logic) checkLoop(s Sender) {
	isTest := l.isTest

	var cnt, period int64 = 0, 30
	f := l.ft.GetFormula()
	for range l.ft.Ping {
		resp := l.ft.GetResponse()
		if cnt%period == 0 {
			info := infoToSend{
				resp:   &resp,
				isTest: isTest,
			}
			s.Send(triggerResponseMessage(info))
			b, _ := l.Draw("", f)
			s.SendPhoto("graph.png", b)
		}
		if resp.AbsDif < 0 {
			info := infoToSend{
				resp:   &resp,
				isTest: isTest,
			}
			s.Send(triggerResponseMessage(info))
			err := l.End(s)
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

func (l *Logic) End(s Sender) error {
	isTest := l.isTest
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

package trigger

import (
	"sync"
	"time"

	"github.com/rwlist/autotrade-bot/binance"
	"github.com/rwlist/autotrade-bot/formula"
	"github.com/rwlist/autotrade-bot/tostr"
)

type FormulaTrigger struct {
	active  bool
	resp    chan *Response
	quit    chan struct{}
	b       binance.MyBinance
	formula formula.Formula
}

func NewTrigger(b binance.MyBinance) FormulaTrigger {
	var resp chan *Response
	var quit chan struct{}
	return FormulaTrigger{
		active: false,
		resp:   resp,
		quit:   quit,
		b:      b,
	}
}

type Response struct {
	rate       float64
	prediction float64
	err        error
}

const timeSleep = 10 * time.Second

func (ft *FormulaTrigger) CheckLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ft.quit:
			return
		default:
			t := float64(time.Now().Unix())
			rate, err := ft.b.GetRate()
			if err != nil {
				ft.resp <- &Response{
					err: err,
				}
			}
			ft.resp <- &Response{
				rate:       tostr.StrToFloat64(rate),
				prediction: ft.formula.Calc(t),
			}
			time.Sleep(timeSleep)
		}
	}
}

func (ft *FormulaTrigger) Begin(f formula.Formula) {
	ft.active = true
	ft.formula = f
	var wg sync.WaitGroup
	wg.Add(1)
	go ft.CheckLoop(&wg)
	wg.Wait()
}

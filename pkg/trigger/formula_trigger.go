package trigger

import (
	"fmt"
	"sync"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/tostr"
	"github.com/rwlist/autotrade-bot/pkg/trade/binance"
)

type Response struct {
	CurRate     float64
	FormulaRate float64
	AbsDif      float64
	RelDif      float64
	StartRate   float64
	AbsProfit   float64
	RelProfit   float64
	T           time.Time
	Err         error
	Formula     string
}

func (ft *FormulaTrigger) newResponse(curRate, fRate float64) *Response {
	absDif := curRate - fRate
	relDif := 100.0 * absDif / fRate
	d := curRate - ft.formula.Rate()
	relProf := 100.0 * d / ft.formula.Rate()
	absProf := d * ft.haveBTC
	return &Response{
		CurRate:     curRate,
		FormulaRate: fRate,
		AbsDif:      absDif,
		RelDif:      relDif,
		StartRate:   ft.formula.Rate(),
		AbsProfit:   absProf,
		RelProfit:   relProf,
		T:           time.Now(),
		Formula:     ft.formula.String(),
	}
}

type FormulaTrigger struct {
	active  bool
	Resp    *Response
	Ping    chan struct{}
	haveBTC float64
	b       binance.Binance
	formula formula.Formula
	mux     sync.Mutex
}

func NewTrigger(b binance.Binance) FormulaTrigger {
	return FormulaTrigger{
		active:  false,
		Resp:    &Response{},
		Ping:    make(chan struct{}),
		haveBTC: 0,
		b:       b,
	}
}

func (ft *FormulaTrigger) IsActive() bool {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	return ft.active
}

func (ft *FormulaTrigger) updActive(state bool) {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	ft.active = state
}

func (ft *FormulaTrigger) UpdFormula(f formula.Formula) {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	ft.formula = f
}

func (ft *FormulaTrigger) GetFormula() formula.Formula {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	return ft.formula
}

func (ft *FormulaTrigger) updBTC() {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	haveBTC, err := ft.b.AccountSymbolBalance("BTC")
	if err != nil {
		haveBTC = 0
	}
	ft.haveBTC = haveBTC
}

const TimeSleep = 10 * time.Second

func (ft *FormulaTrigger) check() *Response {
	t := time.Now().Unix()
	rate, err := ft.b.GetRate()
	if err != nil {
		return &Response{
			Err:     fmt.Errorf("binance.GetRate error: %w: ", err),
			T:       time.Now(),
			Formula: ft.formula.String(),
		}
	}
	return ft.newResponse(tostr.StrToFloat64(rate), ft.formula.Calc(float64(t)))
}

func (ft *FormulaTrigger) UpdResponse() {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	ft.Resp = ft.check()
}

func (ft *FormulaTrigger) GetResponse() Response {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	return *ft.Resp
}

func (ft *FormulaTrigger) CheckLoop() {
	for {
		if ft.IsActive() {
			ft.UpdResponse()
			ft.Ping <- struct{}{}
			time.Sleep(TimeSleep)
		} else {
			return
		}
	}
}

func (ft *FormulaTrigger) Begin(f formula.Formula) {
	ft.updActive(true)
	ft.updBTC()
	ft.UpdFormula(f)
	go ft.CheckLoop()
}

func (ft *FormulaTrigger) End() {
	ft.updActive(false)
}

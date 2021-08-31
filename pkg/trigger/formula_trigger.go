package trigger

import (
	"fmt"
	"sync"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/shopspring/decimal"

	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/trade/binance"
)

type FormulaTrigger struct {
	active  bool
	Resp    *Response
	Ping    chan struct{}
	haveBTC decimal.Decimal
	b       *binance.Binance
	formula formula.Formula
	mux     sync.Mutex
}

func NewTrigger(b *binance.Binance) FormulaTrigger {
	return FormulaTrigger{
		active:  false,
		Resp:    &Response{},
		Ping:    make(chan struct{}),
		haveBTC: decimal.Zero,
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
		haveBTC = decimal.Zero
	}
	ft.haveBTC = haveBTC
}

const TimeSleep = 10 * time.Second

func (ft *FormulaTrigger) check() *Response {
	t := time.Now()
	rate, err := ft.b.GetRate()
	if err != nil {
		return &Response{
			Err:     fmt.Errorf("binance.GetRate error: %w: ", err),
			T:       time.Now(),
			Formula: ft.formula.String(),
		}
	}
	return ft.newResponse(rate, ft.formula.Calc(t))
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

type Response struct {
	CurRate     decimal.Decimal
	FormulaRate decimal.Decimal
	AbsDif      decimal.Decimal
	RelDif      decimal.Decimal
	StartRate   decimal.Decimal
	AbsProfit   decimal.Decimal
	RelProfit   decimal.Decimal
	T           time.Time
	Err         error
	Formula     string
}

func (ft *FormulaTrigger) newResponse(curRate, fRate decimal.Decimal) *Response {
	absDif := curRate.Sub(fRate)                                   // curRate - fRate
	relDif := absDif.Shift(convert.UsefulShift).Div(fRate)         // 100.0 * absDif / fRate
	d := curRate.Sub(ft.formula.Rate())                            // curRate - ft.formula.Rate()
	relProf := d.Shift(convert.UsefulShift).Div(ft.formula.Rate()) // 100.0 * d / ft.formula.Rate()
	absProf := d.Mul(ft.haveBTC)                                   // d * ft.haveBTC
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

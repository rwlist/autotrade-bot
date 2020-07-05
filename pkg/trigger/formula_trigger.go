package trigger

import (
	"sync"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/binance"
	"github.com/rwlist/autotrade-bot/pkg/formula"
	"github.com/rwlist/autotrade-bot/pkg/tostr"
)

type FormulaTrigger struct {
	active  bool
	Resp    chan *Response
	b       binance.Binance
	formula formula.Formula
	haveBTC float64
	mux     sync.Mutex
}

func NewTrigger(b binance.Binance) (FormulaTrigger, error) {
	haveBTC, err := b.AccountSymbolBalance("BTC")
	if err != nil {
		haveBTC = 0
	}
	return FormulaTrigger{
		active:  false,
		Resp:    make(chan *Response),
		b:       b,
		haveBTC: haveBTC,
	}, err
}

type Response struct {
	CurRate     float64
	FormulaRate float64
	AbsDif      float64
	RelDif      float64
	StartRate   float64
	AbsProfit   float64
	RelProfit   float64
	err         error
}

func (ft *FormulaTrigger) isActive() bool {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	return ft.active
}

func (ft *FormulaTrigger) changeActive(state bool) {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	ft.active = state
}

func (ft *FormulaTrigger) changeFormula(f formula.Formula) {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	ft.formula = f
}

func (ft *FormulaTrigger) GetFormula() formula.Formula {
	ft.mux.Lock()
	defer ft.mux.Unlock()
	return ft.formula
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
	}
}

const timeSleep = 10 * time.Second

func (ft *FormulaTrigger) CheckLoop() {
	for {
		if ft.isActive() {
			t := time.Now().Unix()
			rate, err := ft.b.GetRate()
			if err != nil {
				ft.Resp <- &Response{
					err: err,
				}
			}
			ft.Resp <- ft.newResponse(tostr.StrToFloat64(rate), ft.formula.Calc(float64(t)))
			time.Sleep(timeSleep)
		} else {
			return
		}
	}
}

func (ft *FormulaTrigger) Begin(f formula.Formula) {
	ft.changeActive(true)
	ft.changeFormula(f)
	go ft.CheckLoop()
}

func (ft *FormulaTrigger) End() {
	ft.changeActive(false)
}

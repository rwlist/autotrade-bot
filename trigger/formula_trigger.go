package trigger

import (
	"fmt"
	"time"

	"github.com/rwlist/autotrade-bot/binance"
	"github.com/rwlist/autotrade-bot/formula"
	"github.com/rwlist/autotrade-bot/tostr"
)

type FormulaTrigger struct {
	active  bool
	Resp    chan *Response
	quit    chan struct{}
	b       binance.MyBinance
	formula formula.Formula
	haveBTC float64
}

func NewTrigger(b binance.MyBinance, haveBTC float64) FormulaTrigger {
	return FormulaTrigger{
		active:  false,
		Resp:    make(chan *Response),
		quit:    make(chan struct{}),
		b:       b,
		haveBTC: haveBTC,
	}
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

func (ft *FormulaTrigger) newResponse(curRate, fRate float64) *Response {
	absDif := curRate - fRate
	relDif := absDif / fRate
	d := curRate - ft.formula.Rate()
	relProf := d / ft.formula.Rate()
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
		select {
		case <-ft.quit:
			return

		default:
			t := float64(time.Now().Unix())
			rate, err := ft.b.GetRate()
			if err != nil {
				ft.Resp <- &Response{
					err: err,
				}
			}
			ft.Resp <- ft.newResponse(tostr.StrToFloat64(rate), ft.formula.Calc(t))
			time.Sleep(timeSleep)
		}
	}
}

func (ft *FormulaTrigger) Begin(f formula.Formula) {
	ft.active = true
	ft.formula = f
	go ft.CheckLoop()
	fmt.Println("pzdc")
}

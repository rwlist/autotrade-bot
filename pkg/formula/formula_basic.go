package formula

import (
	"fmt"

	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/shopspring/decimal"
)

const cntCoef = 3

//	Базовая функция удовлетворяющая интерфейсу Formula
//	Имеет вид rate-10+0.0002*(now-start)^1.2
//	Парсится через regexp patternBasic
type Basic struct {
	rate  decimal.Decimal
	start int64
	coef  []decimal.Decimal // Числовые коэффициенты
}

func (f *Basic) String() string {
	return fmt.Sprintf("rate-%v+%v*(now-start)^%v", f.coef[0].String(), f.coef[1].String(), f.coef[2].String())
}

// Calc(now) Вычисляет значение в точке now
func (f *Basic) Calc(now float64) float64 {
	return convert.Float64(f.CalcDec(int64(now)))
}

func (f *Basic) CalcDec(now int64) decimal.Decimal {
	brackets := decimal.NewFromInt(now - f.Start())
	tmp := f.coef[1].Mul(convert.Pow(brackets, f.coef[2])) // Нужен Pow получше
	return f.Rate().
		Sub(f.coef[0]).
		Add(tmp)
}

func (f *Basic) Start() int64 {
	return f.start
}

func (f *Basic) Rate() decimal.Decimal {
	return f.rate
}

//	По заданной строке определяется является ли она формулой этого вида
//	Создаёт и возвращает указатель на структуру, если да
func NewBasic(s string, rate decimal.Decimal, start int64) (*Basic, error) {
	coef, err := parseBasic(s)
	if err != nil {
		return nil, err
	}
	return &Basic{
		rate:  rate,
		start: start,
		coef:  coef,
	}, nil
}

// Меняет формулу сохранняя значения rate и start прежними
func (f *Basic) Alter(s string) error {
	coef, err := parseBasic(s)
	if err != nil {
		return err
	}
	f.coef = coef
	return nil
}

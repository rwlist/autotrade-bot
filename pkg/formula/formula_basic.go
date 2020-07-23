package formula

import (
	"fmt"
	"math"
)

const cntCoef = 3

//	Базовая функция удовлетворяющая интерфейсу Formula
//	Имеет вид rate-10+0.0002*(now-start)^1.2
//	Парсится через regexp patternBasic
type Basic struct {
	rate, start float64   // значения rate и start
	coef        []float64 // Числовые коэффициенты
}

func (f *Basic) String() string {
	return fmt.Sprintf("rate-%f+%f*(now-start)^%f", f.coef[0], f.coef[1], f.coef[2])
}

// Calc(now) Вычисляет значение в точке now
func (f *Basic) Calc(now float64) float64 {
	return f.Rate() - f.coef[0] + f.coef[1]*math.Pow(now-f.Start(), f.coef[2])
}

func (f *Basic) Start() float64 {
	return f.start
}

func (f *Basic) Rate() float64 {
	return f.rate
}

//	По заданной строке определяется является ли она формулой этого вида
//	Создаёт и возвращает указатель на структуру, если да
func NewBasic(s string, rate, start float64) (*Basic, error) {
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

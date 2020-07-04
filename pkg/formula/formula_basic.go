package formula

import (
	"errors"
	"math"
	"regexp"

	"github.com/rwlist/autotrade-bot/pkg/tostr"
)

const patternFloat = `[0-9]+\.?[0-9]*`
const patternBasic = `(rate)-[0-9]+\.?[0-9]*\+[0-9]+\.?[0-9]*\*\((now)-(start)\)\^[0-9]+\.?[0-9]*`

/*
	Базовая функция удовлетворяющая интерфейсу Formula
	Имеет вид rate-10+0.0002*(now-start)^1.2
	Парсится через regexp patternBasic
*/
type Basic struct {
	rate, start float64   // значения rate и start
	coef        []float64 // Числовые коэффициенты
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

const cntCoef = 3

/*
	По заданной строке определяется является ли она формулой этого вида
	Создаёт и возвращает указатель на структуру, если да
*/
func NewBasic(s string, rate, start float64) (*Basic, error) {
	re := regexp.MustCompile(patternBasic)
	s = re.FindString(s)
	if s == "" {
		return nil, errors.New("invalid formula format")
	}
	re = regexp.MustCompile(patternFloat)
	nums := re.FindAllString(s, -1)
	if len(nums) != cntCoef {
		return nil, errors.New("invalid formula format")
	}
	var coef []float64
	for _, val := range nums {
		coef = append(coef, tostr.StrToFloat64(val))
	}
	return &Basic{
		rate:  rate,
		start: start,
		coef:  coef,
	}, nil
}

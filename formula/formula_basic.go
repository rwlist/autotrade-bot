package formula

import (
	"errors"
	"math"
	"regexp"

	"github.com/rwlist/autotrade-bot/tostr"
)

const patternFloat string = `[0-9]+\.?[0-9]*`
const patternBasic string = `(rate)-[0-9]+\.?[0-9]*\+[0-9]+\.?[0-9]*\*\((now)-(start)\)\^[0-9]+\.?[0-9]*`

type Basic struct {
	rate, start float64
	Coef        []float64
}

func (f *Basic) Calc(now float64) float64 {
	return f.Rate() - f.Coef[0] + f.Coef[1]*math.Pow(now-f.Start(), f.Coef[2])
}

func (f *Basic) Start() float64 {
	return f.start
}

func (f *Basic) Rate() float64 {
	return f.rate
}

const cntCoef = 3

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
		Coef:  coef,
	}, nil
}

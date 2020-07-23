package formula

import (
	"errors"
	"regexp"

	"github.com/rwlist/autotrade-bot/pkg/tostr"
)

const patternFloat = `[0-9]+\.?[0-9]*`
const patternBasic = `(rate)-[0-9]+\.?[0-9]*\+[0-9]+\.?[0-9]*\*\((now)-(start)\)\^[0-9]+\.?[0-9]*`

func parseBasic(s string) ([]float64, error) {
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
	return coef, nil
}

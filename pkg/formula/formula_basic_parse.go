package formula

import (
	"errors"
	"regexp"

	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/shopspring/decimal"
)

const patternFloat = `[0-9]+\.?[0-9]*`
const patternBasic = `(rate)-[0-9]+\.?[0-9]*\+[0-9]+\.?[0-9]*\*\((now)-(start)\)\^[0-9]+\.?[0-9]*`

func parseBasic(s string) ([]decimal.Decimal, error) {
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
	var coef []decimal.Decimal
	for _, val := range nums {
		coef = append(coef, convert.UnsafeDecimal(val))
	}
	return coef, nil
}

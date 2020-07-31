package formula

import (
	"errors"
	"regexp"

	"github.com/rwlist/autotrade-bot/pkg/convert"
	"github.com/shopspring/decimal"
)

const patternBasic = `rate([-\+][0-9]+\.?[0-9]*)([-\+][0-9]+\.?[0-9]*)\*\(now-start\)\^([0-9]+\.?[0-9]*)`

func parseBasic(s string) ([]decimal.Decimal, error) {
	re := regexp.MustCompile(patternBasic)
	all := re.FindStringSubmatch(s)
	if all == nil {
		return nil, errors.New("invalid formula format")
	}

	all = all[1:]
	var coef []decimal.Decimal
	for _, val := range all {
		coef = append(coef, convert.UnsafeDecimal(val))
	}
	return coef, nil
}

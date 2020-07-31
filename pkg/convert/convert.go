package convert

import (
	"math"

	log "github.com/sirupsen/logrus"

	"github.com/shopspring/decimal"
)

const UsefulShift = 2
const MoneyTrunc = 6

func UnsafeDecimal(str string) decimal.Decimal {
	dec, err := decimal.NewFromString(str)
	if err != nil {
		log.WithField("str", str).WithError(err).Error("in UnsafeDecimal")
		return decimal.Zero
	}
	return dec
}

func Float64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func Pow(x, y decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(math.Pow(Float64(x), Float64(y)))
}

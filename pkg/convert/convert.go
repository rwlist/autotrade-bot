package convert

import (
	"math"

	"github.com/shopspring/decimal"
)

const UsefulShift = 2
const MoneyTrunc = 6

func UnsafeDecimal(str string) decimal.Decimal {
	dec, _ := decimal.NewFromString(str)
	return dec
}

func Float64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func Sum(str1, str2 string) decimal.Decimal {
	return decimal.Sum(UnsafeDecimal(str1), UnsafeDecimal(str2))
}

func Pow(x, y decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(math.Pow(Float64(x), Float64(y)))
}

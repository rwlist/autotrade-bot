package stat

import (
	"github.com/shopspring/decimal"
)

func unsafeDecimal(str string) decimal.Decimal {
	dec, _ := decimal.NewFromString(str)
	return dec
}

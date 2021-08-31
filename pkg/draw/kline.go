package draw

import (
	"time"

	"github.com/shopspring/decimal"
)

type KlineTOHLCV struct {
	T int64
	O float64
	H float64
	L float64
	C float64
	V float64
}

type Klines struct {
	Klines    []KlineTOHLCV
	Last      decimal.Decimal
	Min       decimal.Decimal
	Max       decimal.Decimal
	StartTime time.Time
	Scale     string
}

func (k *Klines) Len() int {
	return len(k.Klines)
}

func (k *Klines) TOHLCV(i int) (t, o, h, l, c, v float64) { //nolint:gocritic
	return float64(k.Klines[i].T), k.Klines[i].O, k.Klines[i].H, k.Klines[i].L, k.Klines[i].C, k.Klines[i].V
}

type KlinesOpts struct {
	Symbol string
	T      string
}

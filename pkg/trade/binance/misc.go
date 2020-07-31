package binance

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/rwlist/autotrade-bot/pkg/draw"
)

const timeShift = 1000
const hday = 24
const dweek = 7
const dmonth = 30
const klinesCount = 96

func calcStartTime(s string) int64 {
	t := time.Now()
	var d time.Duration
	switch s[len(s)-1] {
	case 'm':
		d = -time.Minute

	case 'H':
		d = -time.Hour

	case 'D':
		d = -time.Hour * hday

	case 'W':
		d = -time.Hour * hday * dweek

	case 'M':
		d = -time.Hour * hday * dmonth
	}
	d *= time.Duration(klinesCount * convert.Int(s[:len(s)-1]))
	return int64(timeShift) * t.Add(d).Unix()
}

func (b *Binance) klinesOpts() draw.KlinesOpts {
	return draw.KlinesOpts{
		Symbol: b.opts.Symbol,
		T:      b.opts.Scale,
	}
}

func sum(str1, str2 string) decimal.Decimal {
	return decimal.Sum(convert.UnsafeDecimal(str1), convert.UnsafeDecimal(str2))
}

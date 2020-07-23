package binance

import (
	"time"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/tostr"
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
	d *= time.Duration(klinesCount * tostr.Int(s[:len(s)-1]))
	return int64(timeShift) * t.Add(d).Unix()
}

func sum(str1, str2 string) float64 {
	return tostr.StrToFloat64(str1) + tostr.StrToFloat64(str2)
}

func (b *Binance) klinesOpts() draw.KlinesOpts {
	return draw.KlinesOpts{
		Symbol: b.opts.Symbol,
		T:      b.opts.Scale,
	}
}

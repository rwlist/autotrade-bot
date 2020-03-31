package draw

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
	Last      float64
	Min       float64
	Max       float64
	StartTime float64
}

func (k Klines) Len() int {
	return len(k.Klines)
}

func (k Klines) TOHLCV(i int) (float64, float64, float64, float64, float64, float64) {
	return float64(k.Klines[i].T), k.Klines[i].O, k.Klines[i].H, k.Klines[i].L, k.Klines[i].C, k.Klines[i].V
}

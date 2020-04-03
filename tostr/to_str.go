package tostr

import (
	"fmt"
	"math"
	"strconv"
)

func Str(id int) string {
	return strconv.Itoa(id)
}

func Str64(id int64) string {
	return fmt.Sprintf("%v", id)
}

func Float64ToStr(f float64, d int) string {
	return fmt.Sprintf("%f", Round(f, d))
}

func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func Round(f float64, d int) float64 {
	return float64(int(f*math.Pow10(d))) / math.Pow10(d)
}

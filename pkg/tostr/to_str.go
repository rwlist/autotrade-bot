package tostr

import (
	"fmt"
	"math"
	"strconv"
)

func Str(id int) string {
	return strconv.Itoa(id)
}

func Int(s string) int64 {
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

func Str64(id int64) string {
	return fmt.Sprintf("%v", id)
}

func Float64ToStr(f float64, decimals int) string {
	return fmt.Sprintf("%f", Round(f, decimals))
}

func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func Round(f float64, decimals int) float64 {
	return float64(int(f*math.Pow10(decimals))) / math.Pow10(decimals)
}

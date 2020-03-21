package app

import (
	"fmt"
	"math"
	"strconv"
)

func str(id int) string {
	return strconv.Itoa(id)
}

func float64ToStr(f float64, d int) string {
	return fmt.Sprintf("%f", Round(f, d))
}

func strToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func Round(f float64, d int) float64 {
	return float64(int(f * math.Pow10(d))) / math.Pow10(d)
}
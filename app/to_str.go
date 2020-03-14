package app

import (
	"fmt"
	"strconv"
)

func str(id int) string {
	return strconv.Itoa(id)
}

func float64ToStr(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func float64ToStrLong(f float64) string {
	return fmt.Sprintf("%f", f)
}

func strToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
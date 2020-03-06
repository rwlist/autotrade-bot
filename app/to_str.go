package app

import (
	"fmt"
	"strconv"
)

func str(id int) string {
	return strconv.Itoa(id)
}

func fToStr(f float64) string {
	return fmt.Sprintf("%v", f)
}

func strToF(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
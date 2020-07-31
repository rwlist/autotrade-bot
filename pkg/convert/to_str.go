package convert

import (
	"fmt"
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

func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

package float

import (
	"math"
	"strconv"
)

// Round 将浮点数四舍五入到指定的小数位数
func Round(v, precision float64) float64 {
	return math.Round(v*math.Pow(10, precision)) / math.Pow(10, precision)
}

// Cut 将浮点数截取到指定的小数位数
func Cut(v float64, precision int) float64 {
	f, err := strconv.ParseFloat(strconv.FormatFloat(v, 'f', precision, 64), 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

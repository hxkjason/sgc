package utils

import (
	"math"
	"strconv"
)

// DecimalRoundOffFour 四舍五入 保留几位小数
func DecimalRoundOffFour(num float64, precision float64) float64 {
	return math.Round(num*math.Pow(10, precision)) / math.Pow(10, precision)
}

// DecimalFloat 保留几位小数
func DecimalFloat(value float64, precision int) float64 {
	value, _ = strconv.ParseFloat(strconv.FormatFloat(value, 'f', precision, 64), 64)
	return value
}

// FormatFloatToStr 强制转化为带有几位小数的字符串
func FormatFloatToStr(v float64, precision int) string {
	if precision > 10 {
		precision = 10
	}
	return strconv.FormatFloat(v, 'f', precision, 64)
}

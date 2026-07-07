package xmath

import "math"

// RoundFloat64 按指定小数位四舍五入。
//
// 参数：
//   - value: 待处理值。
//   - precision: 保留小数位数；负数表示按十、百、千等位数处理。
//
// 返回值：
//   - float64: 四舍五入后的值。
func RoundFloat64(value float64, precision int) float64 {
	if precision == 0 {
		return math.Round(value)
	}
	if precision > 0 {
		scale := math.Pow10(precision)
		return math.Round(value*scale) / scale
	}

	scale := math.Pow10(-precision)
	return math.Round(value/scale) * scale
}

// FloorFloat64 按指定小数位向下取整。
//
// 参数：
//   - value: 待处理值。
//   - precision: 保留小数位数；负数表示按十、百、千等位数处理。
//
// 返回值：
//   - float64: 向下取整后的值。
func FloorFloat64(value float64, precision int) float64 {
	if precision == 0 {
		return math.Floor(value)
	}
	if precision > 0 {
		scale := math.Pow10(precision)
		return math.Floor(value*scale) / scale
	}

	scale := math.Pow10(-precision)
	return math.Floor(value/scale) * scale
}

// CeilFloat64 按指定小数位向上取整。
//
// 参数：
//   - value: 待处理值。
//   - precision: 保留小数位数；负数表示按十、百、千等位数处理。
//
// 返回值：
//   - float64: 向上取整后的值。
func CeilFloat64(value float64, precision int) float64 {
	if precision == 0 {
		return math.Ceil(value)
	}
	if precision > 0 {
		scale := math.Pow10(precision)
		return math.Ceil(value*scale) / scale
	}

	scale := math.Pow10(-precision)
	return math.Ceil(value/scale) * scale
}

// TruncFloat64 按指定小数位截断。
//
// 参数：
//   - value: 待处理值。
//   - precision: 保留小数位数；负数表示按十、百、千等位数处理。
//
// 返回值：
//   - float64: 截断后的值。
func TruncFloat64(value float64, precision int) float64 {
	if precision == 0 {
		return math.Trunc(value)
	}
	if precision > 0 {
		scale := math.Pow10(precision)
		return math.Trunc(value*scale) / scale
	}

	scale := math.Pow10(-precision)
	return math.Trunc(value/scale) * scale
}

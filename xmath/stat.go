package xmath

// SumInt 返回 int 切片求和值。
//
// 参数：
//   - values: 待求和数据。
//
// 返回值：
//   - int: values 中所有元素的和。
func SumInt(values []int) int {
	var total int
	for _, value := range values {
		total += value
	}
	return total
}

// SumInt64 返回 int64 切片求和值。
//
// 参数：
//   - values: 待求和数据。
//
// 返回值：
//   - int64: values 中所有元素的和。
func SumInt64(values []int64) int64 {
	var total int64
	for _, value := range values {
		total += value
	}
	return total
}

// SumFloat64 返回 float64 切片求和值。
//
// 参数：
//   - values: 待求和数据。
//
// 返回值：
//   - float64: values 中所有元素的和。
func SumFloat64(values []float64) float64 {
	var total float64
	for _, value := range values {
		total += value
	}
	return total
}

// AvgInt 返回 int 切片平均值。
//
// 参数：
//   - values: 待计算数据。
//
// 返回值：
//   - float64: 平均值。
//   - bool: values 非空时返回 true。
func AvgInt(values []int) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	return float64(SumInt(values)) / float64(len(values)), true
}

// AvgInt64 返回 int64 切片平均值。
//
// 参数：
//   - values: 待计算数据。
//
// 返回值：
//   - float64: 平均值。
//   - bool: values 非空时返回 true。
func AvgInt64(values []int64) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	return float64(SumInt64(values)) / float64(len(values)), true
}

// AvgFloat64 返回 float64 切片平均值。
//
// 参数：
//   - values: 待计算数据。
//
// 返回值：
//   - float64: 平均值。
//   - bool: values 非空时返回 true。
func AvgFloat64(values []float64) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	return SumFloat64(values) / float64(len(values)), true
}

// Percent 返回 part / total * 100 的百分比值。
//
// 参数：
//   - part: 分子。
//   - total: 分母。
//
// 返回值：
//   - float64: 百分比值。
//   - bool: total 不为 0 时返回 true。
func Percent(part, total float64) (float64, bool) {
	if total == 0 {
		return 0, false
	}
	return part / total * 100, true
}

// Ratio 返回 part / total 的比值。
//
// 参数：
//   - part: 分子。
//   - total: 分母。
//
// 返回值：
//   - float64: 比值。
//   - bool: total 不为 0 时返回 true。
func Ratio(part, total float64) (float64, bool) {
	if total == 0 {
		return 0, false
	}
	return part / total, true
}

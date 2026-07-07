package xmath

// MinFloat64 返回两个 float64 中较小的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - float64: a 和 b 中较小的值；相等时返回 a。
func MinFloat64(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}

// MaxFloat64 返回两个 float64 中较大的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - float64: a 和 b 中较大的值；相等时返回 a。
func MaxFloat64(a, b float64) float64 {
	if a >= b {
		return a
	}
	return b
}

// ClampFloat64 将 value 限制在闭区间 [min, max] 内。
//
// 参数：
//   - value: 待限制值。
//   - min: 区间下限。
//   - max: 区间上限。
//
// 返回值：
//   - float64: 限制后的值。
//
// min > max 时会自动交换。
func ClampFloat64(value, min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// InRangeFloat64 判断 value 是否在闭区间 [min, max] 内。
//
// 参数：
//   - value: 待判断值。
//   - min: 区间下限。
//   - max: 区间上限。
//
// 返回值：
//   - bool: value 在闭区间内时返回 true。
//
// min > max 时会自动交换。
func InRangeFloat64(value, min, max float64) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

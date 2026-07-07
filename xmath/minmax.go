package xmath

// MinInt 返回两个 int 中较小的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - int: a 和 b 中较小的值；相等时返回 a。
func MinInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// MinInt64 返回两个 int64 中较小的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - int64: a 和 b 中较小的值；相等时返回 a。
func MinInt64(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

// MaxInt 返回两个 int 中较大的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - int: a 和 b 中较大的值；相等时返回 a。
func MaxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// MaxInt64 返回两个 int64 中较大的值。
//
// 参数：
//   - a: 待比较值。
//   - b: 待比较值。
//
// 返回值：
//   - int64: a 和 b 中较大的值；相等时返回 a。
func MaxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

// AbsInt 返回 int 的绝对值。
//
// 参数：
//   - value: 待处理值。
//
// 返回值：
//   - int: value 的绝对值。
//
// 当 value 为最小 int 时无法表示对应正数，会原样返回。
func AbsInt(value int) int {
	if value >= 0 {
		return value
	}
	result := -value
	if result < 0 {
		return value
	}
	return result
}

// AbsInt64 返回 int64 的绝对值。
//
// 参数：
//   - value: 待处理值。
//
// 返回值：
//   - int64: value 的绝对值。
//
// 当 value 为 math.MinInt64 时无法表示对应正数，会原样返回。
func AbsInt64(value int64) int64 {
	if value >= 0 {
		return value
	}
	result := -value
	if result < 0 {
		return value
	}
	return result
}

// ClampInt 将 value 限制在闭区间 [min, max] 内。
//
// 参数：
//   - value: 待限制值。
//   - min: 区间下限。
//   - max: 区间上限。
//
// 返回值：
//   - int: 限制后的值。
//
// min > max 时会自动交换。
func ClampInt(value, min, max int) int {
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

// ClampInt64 将 value 限制在闭区间 [min, max] 内。
//
// 参数：
//   - value: 待限制值。
//   - min: 区间下限。
//   - max: 区间上限。
//
// 返回值：
//   - int64: 限制后的值。
//
// min > max 时会自动交换。
func ClampInt64(value, min, max int64) int64 {
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

// InRangeInt 判断 value 是否在闭区间 [min, max] 内。
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
func InRangeInt(value, min, max int) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

// InRangeInt64 判断 value 是否在闭区间 [min, max] 内。
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
func InRangeInt64(value, min, max int64) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

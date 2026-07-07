package xrand

import (
	"math/rand"
)

// RangeInt 返回闭区间 [min, max] 内的随机 int。
//
// 参数：
//   - min: 随机范围下限。
//   - max: 随机范围上限。
//
// 返回值：
//   - int: 闭区间 [min, max] 内的随机值。
//
// min > max 时会自动交换；min == max 时直接返回该值。
func RangeInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}

	span := uint(max) - uint(min) + 1
	return int(uint(min) + uint(rangeUint64(uint64(span))))
}

// RangeInt64 返回闭区间 [min, max] 内的随机 int64。
//
// 参数：
//   - min: 随机范围下限。
//   - max: 随机范围上限。
//
// 返回值：
//   - int64: 闭区间 [min, max] 内的随机值。
//
// min > max 时会自动交换；min == max 时直接返回该值。
func RangeInt64(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}

	span := uint64(max) - uint64(min) + 1
	return int64(uint64(min) + rangeUint64(span))
}

// RangeFloat64 返回半开区间 [min, max) 内的随机 float64。
//
// 参数：
//   - min: 随机范围下限。
//   - max: 随机范围上限。
//
// 返回值：
//   - float64: 半开区间 [min, max) 内的随机值。
//
// min > max 时会自动交换；min == max 时直接返回该值。
func RangeFloat64(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

// Bool 返回随机布尔值。
//
// 返回值：
//   - bool: 随机 true 或 false。
func Bool() bool {
	return rand.Intn(2) == 1
}

// Pick 从 src 中随机选择 total 个元素返回。
//
// 参数：
//   - src: 原始数据列表。
//   - total: 需要返回的元素数量。
//
// 返回值：
//   - []T: 随机挑选结果；total <= 0 或 src 为空时返回空切片。
//
// 不修改 src；total > len(src) 时会循环补充，允许重复顺序片段。
func Pick[T any](src []T, total int) []T {
	if total <= 0 || len(src) == 0 {
		return []T{}
	}

	shuffled := make([]T, len(src))
	copy(shuffled, src)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := make([]T, 0, total)
	for len(result) < total {
		need := total - len(result)
		if need >= len(shuffled) {
			result = append(result, shuffled...)
			continue
		}
		result = append(result, shuffled[:need]...)
	}
	return result
}

// PickOne 从 src 中随机选择一个元素返回。
//
// 参数：
//   - src: 原始数据列表。
//
// 返回值：
//   - T: 随机选中的元素。
//   - bool: src 非空时返回 true。
func PickOne[T any](src []T) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}
	return src[rand.Intn(len(src))], true
}

// PickUnique 从 src 中随机选择不重复的 total 个元素返回。
//
// 参数：
//   - src: 原始数据列表。
//   - total: 需要返回的元素数量。
//
// 返回值：
//   - []T: 随机挑选结果；total <= 0 或 src 为空时返回空切片。
//
// 不修改 src；total > len(src) 时最多返回 len(src) 个元素。
func PickUnique[T any](src []T, total int) []T {
	if total <= 0 || len(src) == 0 {
		return []T{}
	}
	if total > len(src) {
		total = len(src)
	}

	shuffled := Shuffle(src)
	return shuffled[:total]
}

// Shuffle 返回 src 的随机打乱副本。
//
// 参数：
//   - src: 原始数据列表。
//
// 返回值：
//   - []T: 随机打乱后的副本。
//
// 不修改 src。
func Shuffle[T any](src []T) []T {
	shuffled := make([]T, len(src))
	copy(shuffled, src)
	ShuffleInPlace(shuffled)
	return shuffled
}

// ShuffleInPlace 原地随机打乱 src。
//
// 参数：
//   - src: 待打乱的数据列表。
func ShuffleInPlace[T any](src []T) {
	rand.Shuffle(len(src), func(i, j int) {
		src[i], src[j] = src[j], src[i]
	})
}

func rangeUint64(span uint64) uint64 {
	if span == 0 {
		return rand.Uint64()
	}

	maxUint64 := ^uint64(0)
	limit := maxUint64 - (maxUint64 % span)
	for {
		n := rand.Uint64()
		if n < limit {
			return n % span
		}
	}
}

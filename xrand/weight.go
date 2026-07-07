package xrand

// WeightedIndex 根据权重随机返回索引。
//
// 参数：
//   - weights: 权重列表，权重大于 0 的项才会参与随机。
//
// 返回值：
//   - int: 命中的索引。
//   - bool: 存在可用权重时返回 true。
func WeightedIndex(weights []int) (int, bool) {
	var total int64
	for _, weight := range weights {
		if weight > 0 {
			total += int64(weight)
		}
	}
	if total <= 0 {
		return 0, false
	}

	hit := RangeInt64(1, total)
	var current int64
	for i, weight := range weights {
		if weight <= 0 {
			continue
		}
		current += int64(weight)
		if hit <= current {
			return i, true
		}
	}
	return 0, false
}

// WeightedPick 根据权重随机选择一个元素。
//
// 参数：
//   - values: 候选数据列表。
//   - weights: 权重列表，权重大于 0 的项才会参与随机。
//
// 返回值：
//   - T: 命中的元素。
//   - bool: 存在可用权重且命中索引在 values 范围内时返回 true。
//
// 当 weights 长度大于 values 时，多余权重会被忽略；当 weights 长度小于 values 时，多余元素不会参与随机。
func WeightedPick[T any](values []T, weights []int) (T, bool) {
	var zero T
	if len(values) == 0 || len(weights) == 0 {
		return zero, false
	}
	if len(weights) > len(values) {
		weights = weights[:len(values)]
	}

	index, ok := WeightedIndex(weights)
	if !ok || index >= len(values) {
		return zero, false
	}
	return values[index], true
}

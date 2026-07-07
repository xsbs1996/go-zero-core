package xrand

import "math/rand"

// Chance 按概率返回是否命中。
//
// 参数：
//   - probability: 命中概率，取值范围通常为 [0, 1]。
//
// 返回值：
//   - bool: 命中时返回 true。
//
// probability <= 0 时固定返回 false；probability >= 1 时固定返回 true。
func Chance(probability float64) bool {
	if probability <= 0 {
		return false
	}
	if probability >= 1 {
		return true
	}
	return rand.Float64() < probability
}

// ChancePercent 按百分比概率返回是否命中。
//
// 参数：
//   - percent: 命中百分比，取值范围通常为 [0, 100]。
//
// 返回值：
//   - bool: 命中时返回 true。
//
// percent <= 0 时固定返回 false；percent >= 100 时固定返回 true。
func ChancePercent(percent float64) bool {
	return Chance(percent / 100)
}

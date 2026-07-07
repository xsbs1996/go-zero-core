package xrand

import "math/rand"

// Permutation 返回 [0, n) 的随机排列。
//
// 参数：
//   - n: 排列长度。
//
// 返回值：
//   - []int: 随机排列；n <= 0 时返回空切片。
func Permutation(n int) []int {
	if n <= 0 {
		return []int{}
	}
	return rand.Perm(n)
}

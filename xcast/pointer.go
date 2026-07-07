package xcast

// Ptr 返回 v 的指针。
//
// 参数：
//   - v: 任意值。
//
// 返回值：
//   - *T: 指向 v 副本的指针。
func Ptr[T any](v T) *T {
	return &v
}

// Value 返回指针指向的值，指针为 nil 时返回零值。
//
// 参数：
//   - v: 待取值指针。
//
// 返回值：
//   - T: 指针指向的值；v 为 nil 时返回 T 的零值。
func Value[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// ValueOr 返回指针指向的值，指针为 nil 时返回默认值。
//
// 参数：
//   - v: 待取值指针。
//   - defaultValue: v 为 nil 时返回的默认值。
//
// 返回值：
//   - T: 指针指向的值；v 为 nil 时返回 defaultValue。
func ValueOr[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

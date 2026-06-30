package convert

// Ptr 返回 v 的指针。
func Ptr[T any](v T) *T {
	return &v
}

// Value 返回指针指向的值，指针为 nil 时返回零值。
func Value[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// ValueOr 返回指针指向的值，指针为 nil 时返回默认值。
func ValueOr[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

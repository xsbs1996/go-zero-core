package convert

import "strconv"

// BoolToString 将 bool 转换为 "true" 或 "false"。
func BoolToString(v bool) string {
	return strconv.FormatBool(v)
}

// BoolToInt 将 true 转换为 1，将 false 转换为 0。
func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// IntToBool 将非 0 整数转换为 true。
func IntToBool(v int) bool {
	return v != 0
}

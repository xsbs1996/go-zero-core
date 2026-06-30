package convert

import "strconv"

// IntToString 将 int 转换为十进制字符串。
func IntToString(v int) string {
	return strconv.Itoa(v)
}

// Int8ToString 将 int8 转换为十进制字符串。
func Int8ToString(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int16ToString 将 int16 转换为十进制字符串。
func Int16ToString(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int32ToString 将 int32 转换为十进制字符串。
func Int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int64ToString 将 int64 转换为十进制字符串。
func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

// UintToString 将 uint 转换为十进制字符串。
func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint8ToString 将 uint8 转换为十进制字符串。
func Uint8ToString(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint16ToString 将 uint16 转换为十进制字符串。
func Uint16ToString(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint32ToString 将 uint32 转换为十进制字符串。
func Uint32ToString(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint64ToString 将 uint64 转换为十进制字符串。
func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// Float32ToString 将 float32 转换为字符串，并去掉不必要的末尾 0。
func Float32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// Float64ToString 将 float64 转换为字符串，并去掉不必要的末尾 0。
func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// IntToInt64 将 int 转换为 int64。
func IntToInt64(v int) int64 {
	return int64(v)
}

// Int64ToInt 将 int64 转换为 int。
func Int64ToInt(v int64) int {
	return int(v)
}

// IntToFloat64 将 int 转换为 float64。
func IntToFloat64(v int) float64 {
	return float64(v)
}

// Int64ToFloat64 将 int64 转换为 float64。
func Int64ToFloat64(v int64) float64 {
	return float64(v)
}

// Float64ToInt 将 float64 转换为 int，小数部分会被截断。
func Float64ToInt(v float64) int {
	return int(v)
}

// Float64ToInt64 将 float64 转换为 int64，小数部分会被截断。
func Float64ToInt64(v float64) int64 {
	return int64(v)
}

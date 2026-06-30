package convert

import (
	"strconv"
	"strings"
	"time"
)

// StringToInt 去除首尾空格后将字符串转换为 int。
func StringToInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// StringToInt8 去除首尾空格后将字符串转换为 int8。
func StringToInt8(s string) (int8, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 8)
	return int8(v), err
}

// StringToInt16 去除首尾空格后将字符串转换为 int16。
func StringToInt16(s string) (int16, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 16)
	return int16(v), err
}

// StringToInt32 去除首尾空格后将字符串转换为 int32。
func StringToInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int32(v), err
}

// StringToInt64 去除首尾空格后将字符串转换为 int64。
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

// StringToUint 去除首尾空格后将字符串转换为 uint。
func StringToUint(s string) (uint, error) {
	v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 0)
	return uint(v), err
}

// StringToUint8 去除首尾空格后将字符串转换为 uint8。
func StringToUint8(s string) (uint8, error) {
	v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 8)
	return uint8(v), err
}

// StringToUint16 去除首尾空格后将字符串转换为 uint16。
func StringToUint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	return uint16(v), err
}

// StringToUint32 去除首尾空格后将字符串转换为 uint32。
func StringToUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	return uint32(v), err
}

// StringToUint64 去除首尾空格后将字符串转换为 uint64。
func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 10, 64)
}

// StringToFloat32 去除首尾空格后将字符串转换为 float32。
func StringToFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	return float32(v), err
}

// StringToFloat64 去除首尾空格后将字符串转换为 float64。
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// StringToBool 去除首尾空格后将字符串转换为 bool。
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(s))
}

// StringToTime 去除首尾空格后，按指定布局和本地时区解析时间字符串。
func StringToTime(s string, layout string) (time.Time, error) {
	return StringToTimeInLocation(s, layout, time.Local)
}

// StringToTimeInLocation 去除首尾空格后，按指定布局和指定时区解析时间字符串。
func StringToTimeInLocation(s string, layout string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, strings.TrimSpace(s), locationOrLocal(loc))
}

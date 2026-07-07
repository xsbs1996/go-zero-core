package xcast

import "encoding/json"

// ToJSON 将值序列化为紧凑 JSON 字符串。
//
// 参数：
//   - v: 待序列化值。
//
// 返回值：
//   - string: 紧凑 JSON 字符串。
//   - error: JSON 序列化失败时返回错误。
func ToJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToPrettyJSON 将值序列化为带缩进的 JSON 字符串。
//
// 参数：
//   - v: 待序列化值。
//
// 返回值：
//   - string: 带缩进的 JSON 字符串。
//   - error: JSON 序列化失败时返回错误。
func ToPrettyJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON 将 JSON 字符串反序列化到 v。
//
// 参数：
//   - s: JSON 字符串。
//   - v: 目标值，通常为指针。
//
// 返回值：
//   - error: JSON 反序列化失败时返回错误。
func FromJSON(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

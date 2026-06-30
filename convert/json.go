package convert

import "encoding/json"

// ToJSON 将值序列化为紧凑 JSON 字符串。
func ToJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToPrettyJSON 将值序列化为带缩进的 JSON 字符串。
func ToPrettyJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON 将 JSON 字符串反序列化到 v。
func FromJSON(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

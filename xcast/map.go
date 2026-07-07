package xcast

import "encoding/json"

// StructToMap 通过 JSON 标签将结构体转换为 map[string]any。
//
// 参数：
//   - v: 待转换结构体或可 JSON 序列化的值。
//
// 返回值：
//   - map[string]any: 转换后的 map。
//   - error: JSON 序列化或反序列化失败时返回错误。
func StructToMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// MapToStruct 通过 JSON 标签将 map[string]any 转换为结构体。
//
// 参数：
//   - m: 待转换 map。
//   - v: 目标结构体指针。
//
// 返回值：
//   - error: JSON 序列化或反序列化失败时返回错误。
func MapToStruct(m map[string]any, v any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

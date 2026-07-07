package xcast

import (
	"sort"
	"strings"
)

// JoinInt64List 将 int64 列表序列化为 JSON 数组字符串。
//
// 参数：
//   - values: 待序列化 int64 列表。
//
// 返回值：
//   - string: JSON 数组字符串；过滤后为空或序列化失败时返回空字符串。
//
// 会过滤 <=0、去重并升序排序；最终为空时返回空字符串。
func JoinInt64List(values []int64) string {
	normalized := NormalizePositiveUniqueSortedInt64(values)
	if len(normalized) == 0 {
		return ""
	}
	encoded, err := ToJSON(normalized)
	if err != nil {
		return ""
	}
	return encoded
}

// ParseJSONInt64List 将 JSON 数组字符串解析为 int64 列表。
//
// 参数：
//   - raw: JSON 数组字符串。
//
// 返回值：
//   - []int64: 过滤、去重、排序后的正整数列表；空字符串或非法 JSON 返回 nil。
//
// 空字符串或非法 JSON 返回 nil；返回值会过滤 <=0、去重并升序排序。
func ParseJSONInt64List(raw string) []int64 {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	var values []int64
	if err := FromJSON(trimmed, &values); err != nil {
		return nil
	}
	return NormalizePositiveUniqueSortedInt64(values)
}

// ParseInt64List 将字符串解析为 int64 列表。
//
// 参数：
//   - raw: JSON 数组字符串或逗号分隔字符串。
//
// 返回值：
//   - []int64: 过滤、去重、排序后的正整数列表；输入为空或无有效数字时返回 nil。
//
// 兼容 JSON 数组字符串（如 "[11,12,13]"）和逗号分隔字符串（如 "11,12,13"）。
func ParseInt64List(raw string) []int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		if list := ParseJSONInt64List(raw); len(list) > 0 {
			return list
		}
	}

	parts := strings.Split(raw, ",")
	out := make([]int64, 0, len(parts))
	for _, part := range parts {
		value, err := StringToInt64(part)
		if err != nil {
			continue
		}
		out = append(out, value)
	}
	return NormalizePositiveUniqueSortedInt64(out)
}

// NormalizePositiveUniqueSortedInt64 过滤 <=0、去重并升序排序。
//
// 参数：
//   - values: 原始 int64 列表。
//
// 返回值：
//   - []int64: 只包含正数、去重且升序排序的列表；结果为空时返回 nil。
func NormalizePositiveUniqueSortedInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	out := make([]int64, 0, len(values))
	seen := make(map[int64]struct{}, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	if len(out) == 0 {
		return nil
	}
	return out
}

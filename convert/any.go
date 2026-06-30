package convert

import "fmt"

// AnyToString 将常见标量类型转换为字符串。
func AnyToString(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case []byte:
		return string(value), nil
	case bool:
		return BoolToString(value), nil
	case int:
		return IntToString(value), nil
	case int8:
		return Int8ToString(value), nil
	case int16:
		return Int16ToString(value), nil
	case int32:
		return Int32ToString(value), nil
	case int64:
		return Int64ToString(value), nil
	case uint:
		return UintToString(value), nil
	case uint8:
		return Uint8ToString(value), nil
	case uint16:
		return Uint16ToString(value), nil
	case uint32:
		return Uint32ToString(value), nil
	case uint64:
		return Uint64ToString(value), nil
	case float32:
		return Float32ToString(value), nil
	case float64:
		return Float64ToString(value), nil
	default:
		return "", fmt.Errorf("unsupported type %T for string conversion", v)
	}
}

// AnyToInt 将常见标量类型转换为 int。
func AnyToInt(v any) (int, error) {
	switch value := v.(type) {
	case int:
		return value, nil
	case int8:
		return int(value), nil
	case int16:
		return int(value), nil
	case int32:
		return int(value), nil
	case int64:
		return int(value), nil
	case uint:
		return int(value), nil
	case uint8:
		return int(value), nil
	case uint16:
		return int(value), nil
	case uint32:
		return int(value), nil
	case uint64:
		return int(value), nil
	case float32:
		return int(value), nil
	case float64:
		return int(value), nil
	case bool:
		return BoolToInt(value), nil
	case string:
		return StringToInt(value)
	case []byte:
		return StringToInt(string(value))
	default:
		return 0, fmt.Errorf("unsupported type %T for int conversion", v)
	}
}

// AnyToInt64 将常见标量类型转换为 int64。
func AnyToInt64(v any) (int64, error) {
	switch value := v.(type) {
	case int:
		return int64(value), nil
	case int8:
		return int64(value), nil
	case int16:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int64:
		return value, nil
	case uint:
		return int64(value), nil
	case uint8:
		return int64(value), nil
	case uint16:
		return int64(value), nil
	case uint32:
		return int64(value), nil
	case uint64:
		return int64(value), nil
	case float32:
		return int64(value), nil
	case float64:
		return int64(value), nil
	case bool:
		return int64(BoolToInt(value)), nil
	case string:
		return StringToInt64(value)
	case []byte:
		return StringToInt64(string(value))
	default:
		return 0, fmt.Errorf("unsupported type %T for int64 conversion", v)
	}
}

// AnyToFloat64 将常见标量类型转换为 float64。
func AnyToFloat64(v any) (float64, error) {
	switch value := v.(type) {
	case int:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case uint8:
		return float64(value), nil
	case uint16:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case bool:
		return float64(BoolToInt(value)), nil
	case string:
		return StringToFloat64(value)
	case []byte:
		return StringToFloat64(string(value))
	default:
		return 0, fmt.Errorf("unsupported type %T for float64 conversion", v)
	}
}

// AnyToBool 将常见标量类型转换为 bool。
func AnyToBool(v any) (bool, error) {
	switch value := v.(type) {
	case bool:
		return value, nil
	case int:
		return IntToBool(value), nil
	case int8:
		return value != 0, nil
	case int16:
		return value != 0, nil
	case int32:
		return value != 0, nil
	case int64:
		return value != 0, nil
	case uint:
		return value != 0, nil
	case uint8:
		return value != 0, nil
	case uint16:
		return value != 0, nil
	case uint32:
		return value != 0, nil
	case uint64:
		return value != 0, nil
	case float32:
		return value != 0, nil
	case float64:
		return value != 0, nil
	case string:
		return StringToBool(value)
	case []byte:
		return StringToBool(string(value))
	default:
		return false, fmt.Errorf("unsupported type %T for bool conversion", v)
	}
}

// AnyToStringOrDefault 将 v 转换为字符串，转换失败时返回默认值。
func AnyToStringOrDefault(v any, defaultValue string) string {
	value, err := AnyToString(v)
	if err != nil {
		return defaultValue
	}
	return value
}

// AnyToIntOrDefault 将 v 转换为 int，转换失败时返回默认值。
func AnyToIntOrDefault(v any, defaultValue int) int {
	value, err := AnyToInt(v)
	if err != nil {
		return defaultValue
	}
	return value
}

// AnyToInt64OrDefault 将 v 转换为 int64，转换失败时返回默认值。
func AnyToInt64OrDefault(v any, defaultValue int64) int64 {
	value, err := AnyToInt64(v)
	if err != nil {
		return defaultValue
	}
	return value
}

// AnyToFloat64OrDefault 将 v 转换为 float64，转换失败时返回默认值。
func AnyToFloat64OrDefault(v any, defaultValue float64) float64 {
	value, err := AnyToFloat64(v)
	if err != nil {
		return defaultValue
	}
	return value
}

// AnyToBoolOrDefault 将 v 转换为 bool，转换失败时返回默认值。
func AnyToBoolOrDefault(v any, defaultValue bool) bool {
	value, err := AnyToBool(v)
	if err != nil {
		return defaultValue
	}
	return value
}

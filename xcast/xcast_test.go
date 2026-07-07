package xcast

import (
	"reflect"
	"testing"
	"time"
)

// TestScalarConversions 验证基础标量类型转换。
func TestScalarConversions(t *testing.T) {
	if got := BoolToString(true); got != "true" {
		t.Fatalf("BoolToString(true) = %q", got)
	}
	if got := BoolToInt(true); got != 1 {
		t.Fatalf("BoolToInt(true) = %d", got)
	}
	if got, err := StringToInt("42"); err != nil || got != 42 {
		t.Fatalf("StringToInt() = %d, %v", got, err)
	}
	if got, err := StringToBool("true"); err != nil || !got {
		t.Fatalf("StringToBool() = %v, %v", got, err)
	}
	if got := Float64ToInt64(12.9); got != 12 {
		t.Fatalf("Float64ToInt64() = %d", got)
	}
}

// TestAnyConversions 验证 any 到常用标量类型的转换和默认值兜底。
func TestAnyConversions(t *testing.T) {
	if got, err := AnyToString([]byte("abc")); err != nil || got != "abc" {
		t.Fatalf("AnyToString() = %q, %v", got, err)
	}
	if got, err := AnyToInt("12"); err != nil || got != 12 {
		t.Fatalf("AnyToInt() = %d, %v", got, err)
	}
	if got, err := AnyToUint64(float64(1001)); err != nil || got != 1001 {
		t.Fatalf("AnyToUint64() = %d, %v", got, err)
	}
	if got, err := AnyToFloat64(true); err != nil || got != 1 {
		t.Fatalf("AnyToFloat64() = %f, %v", got, err)
	}
	if got := AnyToBoolOrDefault(struct{}{}, true); !got {
		t.Fatal("AnyToBoolOrDefault() should return default value")
	}
}

// TestJSONAndMapConversions 验证 JSON 编解码和 struct/map 互转。
func TestJSONAndMapConversions(t *testing.T) {
	type item struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	text, err := ToJSON(item{Name: "alice", Age: 18})
	if err != nil || text != `{"name":"alice","age":18}` {
		t.Fatalf("ToJSON() = %q, %v", text, err)
	}

	var decoded item
	if err := FromJSON(text, &decoded); err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}
	if decoded.Name != "alice" || decoded.Age != 18 {
		t.Fatalf("decoded mismatch: %#v", decoded)
	}

	m, err := StructToMap(decoded)
	if err != nil {
		t.Fatalf("StructToMap() error = %v", err)
	}
	if !reflect.DeepEqual(m, map[string]any{"name": "alice", "age": float64(18)}) {
		t.Fatalf("StructToMap() = %#v", m)
	}
}

// TestInt64ListConversions 验证 int64 列表解析、规范化和序列化。
func TestInt64ListConversions(t *testing.T) {
	want := []int64{1, 2, 3}
	if got := ParseJSONInt64List("[3,1,2,2,0,-1]"); !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseJSONInt64List() = %#v", got)
	}
	if got := ParseInt64List("3,1,2,2,0,-1,bad"); !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseInt64List() = %#v", got)
	}
	if got := JoinInt64List([]int64{3, 1, 2, 2, 0, -1}); got != "[1,2,3]" {
		t.Fatalf("JoinInt64List() = %q", got)
	}
}

// TestPointerAndTimeConversions 验证指针辅助函数和时间戳转换。
func TestPointerAndTimeConversions(t *testing.T) {
	value := Ptr("ok")
	if got := Value(value); got != "ok" {
		t.Fatalf("Value() = %q", got)
	}
	if got := ValueOr((*string)(nil), "fallback"); got != "fallback" {
		t.Fatalf("ValueOr() = %q", got)
	}

	ts := time.Date(2026, 7, 6, 1, 2, 3, 4*int(time.Millisecond), time.UTC)
	if got := TimeToDateTimeString(ts); got != "2026-07-06 01:02:03" {
		t.Fatalf("TimeToDateTimeString() = %q", got)
	}
	if got := UnixMilliToTime(TimeToUnixMilli(ts)); !got.Equal(ts) {
		t.Fatalf("UnixMilli roundtrip = %v, want %v", got, ts)
	}
}

package xrand

import (
	"strings"
	"testing"
)

// TestRandomHelpers 验证随机字节、十六进制和字符集随机串生成。
func TestRandomHelpers(t *testing.T) {
	if got, err := Bytes(8); err != nil || len(got) != 8 {
		t.Fatalf("Bytes() len = %d, err = %v", len(got), err)
	}
	if got, err := Hex(4); err != nil || len(got) != 8 {
		t.Fatalf("Hex() = %q, %v", got, err)
	}
	if got, err := Number(6); err != nil || len(got) != 6 || containsOutside(got, NumberLetters) {
		t.Fatalf("Number() = %q, %v", got, err)
	}
	if _, err := String(1, ""); err == nil {
		t.Fatal("String(empty letters) error = nil")
	}
	if _, err := Bytes(-1); err == nil {
		t.Fatal("Bytes(negative) error = nil")
	}
}

// containsOutside 判断字符串中是否存在不属于指定字符集的字符。
func containsOutside(value string, letters string) bool {
	for _, r := range value {
		if !strings.ContainsRune(letters, r) {
			return true
		}
	}
	return false
}

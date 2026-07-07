package xuuid

import "testing"

// TestUUID 验证 UUID 生成、解析和合法性判断。
func TestUUID(t *testing.T) {
	id := New()
	if !IsValid(id) {
		t.Fatalf("New() invalid UUID: %q", id)
	}
	if _, err := Parse(id); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	withoutDash := NewWithoutDash()
	if len(withoutDash) != 32 {
		t.Fatalf("NewWithoutDash() length = %d", len(withoutDash))
	}
	if !IsValid(withoutDash) {
		t.Fatalf("NewWithoutDash() invalid UUID: %q", withoutDash)
	}
}

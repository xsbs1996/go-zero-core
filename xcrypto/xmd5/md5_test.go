package xmd5

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMD5 验证字符串、字节切片和文件 MD5 摘要。
func TestMD5(t *testing.T) {
	const want = "5d41402abc4b2a76b9719d911017c592"
	if got := Sum("hello"); got != want {
		t.Fatalf("Sum() = %q", got)
	}
	if got := SumBytes([]byte("hello")); got != want {
		t.Fatalf("SumBytes() = %q", got)
	}

	path := filepath.Join(t.TempDir(), "data.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if got, err := SumFile(path); err != nil || got != want {
		t.Fatalf("SumFile() = %q, %v", got, err)
	}
}

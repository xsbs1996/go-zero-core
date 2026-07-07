package xsha

import "testing"

// TestSHA 验证 SHA 摘要函数输出。
func TestSHA(t *testing.T) {
	if got := SHA1("hello"); got != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Fatalf("SHA1() = %q", got)
	}
	if got := SHA256("hello"); got != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("SHA256() = %q", got)
	}
	if SHA384("hello") == "" || SHA512("hello") == "" {
		t.Fatal("SHA384/SHA512 returned empty string")
	}
}

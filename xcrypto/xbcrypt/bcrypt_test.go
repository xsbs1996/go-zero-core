package xbcrypt

import "testing"

// TestHashAndCompare 验证 bcrypt 哈希和密码比对。
func TestHashAndCompare(t *testing.T) {
	hash, err := Hash("secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if !Compare("secret", hash) {
		t.Fatal("Compare(correct) = false")
	}
	if Compare("wrong", hash) {
		t.Fatal("Compare(wrong) = true")
	}
}

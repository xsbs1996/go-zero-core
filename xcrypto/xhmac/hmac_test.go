package xhmac

import "testing"

// TestHMAC 验证 HMAC 签名生成和验签。
func TestHMAC(t *testing.T) {
	sign := SHA256Hex("data", "key")
	if sign == "" {
		t.Fatal("SHA256Hex() returned empty string")
	}
	if !VerifySHA256Hex("data", "key", sign) {
		t.Fatal("VerifySHA256Hex(correct) = false")
	}
	if VerifySHA256Hex("data", "wrong", sign) {
		t.Fatal("VerifySHA256Hex(wrong) = true")
	}
	if SHA1Hex("data", "key") == "" || SHA512Hex("data", "key") == "" || SHA256Base64("data", "key") == "" {
		t.Fatal("HMAC helpers returned empty string")
	}
}

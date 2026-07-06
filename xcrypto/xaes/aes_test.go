package xaes

import (
	"bytes"
	"errors"
	"testing"
)

// TestGCMRoundTrip 验证 AES-GCM 加解密可逆。
func TestGCMRoundTrip(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plain := []byte("hello")

	cipherText, err := EncryptGCM(plain, key)
	if err != nil {
		t.Fatalf("EncryptGCM() error = %v", err)
	}
	got, err := DecryptGCM(cipherText, key)
	if err != nil {
		t.Fatalf("DecryptGCM() error = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("DecryptGCM() = %q", got)
	}
}

// TestCBCRoundTripAndErrors 验证 AES-CBC 加解密可逆和错误分支。
func TestCBCRoundTripAndErrors(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("abcdefghijklmnop")
	plain := []byte("hello")

	cipherText, err := EncryptCBC(plain, key, iv)
	if err != nil {
		t.Fatalf("EncryptCBC() error = %v", err)
	}
	got, err := DecryptCBC(cipherText, key, iv)
	if err != nil {
		t.Fatalf("DecryptCBC() error = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("DecryptCBC() = %q", got)
	}
	if _, err := EncryptCBC(plain, []byte("bad"), iv); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("EncryptCBC(invalid key) error = %v", err)
	}
	if _, err := EncryptCBC(plain, key, []byte("bad")); !errors.Is(err, ErrInvalidIV) {
		t.Fatalf("EncryptCBC(invalid iv) error = %v", err)
	}
}

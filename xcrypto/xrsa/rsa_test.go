package xrsa

import (
	"bytes"
	"errors"
	"testing"
)

// TestRSARoundTripAndSignature 验证 RSA 密钥编解码、OAEP 加解密和 PSS 签名验签。
func TestRSARoundTripAndSignature(t *testing.T) {
	privateKey, err := GenerateKey(1024)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	privatePEM := PrivateKeyToPEM(privateKey)
	parsedPrivate, err := ParsePrivateKey(privatePEM)
	if err != nil {
		t.Fatalf("ParsePrivateKey() error = %v", err)
	}

	publicPEM, err := PublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("PublicKeyToPEM() error = %v", err)
	}
	parsedPublic, err := ParsePublicKey(publicPEM)
	if err != nil {
		t.Fatalf("ParsePublicKey() error = %v", err)
	}

	cipherText, err := EncryptOAEP([]byte("hello"), parsedPublic)
	if err != nil {
		t.Fatalf("EncryptOAEP() error = %v", err)
	}
	plain, err := DecryptOAEP(cipherText, parsedPrivate)
	if err != nil {
		t.Fatalf("DecryptOAEP() error = %v", err)
	}
	if !bytes.Equal(plain, []byte("hello")) {
		t.Fatalf("DecryptOAEP() = %q", plain)
	}

	signature, err := SignPSS([]byte("hello"), parsedPrivate)
	if err != nil {
		t.Fatalf("SignPSS() error = %v", err)
	}
	if err := VerifyPSS([]byte("hello"), signature, parsedPublic); err != nil {
		t.Fatalf("VerifyPSS() error = %v", err)
	}
}

// TestParseInvalidKeys 验证非法 RSA PEM 会返回统一错误。
func TestParseInvalidKeys(t *testing.T) {
	if _, err := ParsePrivateKey("bad"); !errors.Is(err, ErrInvalidPrivateKey) {
		t.Fatalf("ParsePrivateKey() error = %v", err)
	}
	if _, err := ParsePublicKey("bad"); !errors.Is(err, ErrInvalidPublicKey) {
		t.Fatalf("ParsePublicKey() error = %v", err)
	}
}

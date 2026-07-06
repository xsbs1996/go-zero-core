package xbase64

import (
	"bytes"
	"testing"
)

// TestBaseEncodingsRoundTrip 验证 Base64、Base64URL、Base62、Base58 编解码可逆。
func TestBaseEncodingsRoundTrip(t *testing.T) {
	data := []byte{0, 0, 1, 2, 3, 255}

	base64Text := Base64Encode(data)
	if got, err := Base64Decode(base64Text); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("Base64 roundtrip = %v, %v", got, err)
	}

	base64URLText := Base64URLEncode(data)
	if got, err := Base64URLDecode(base64URLText); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("Base64URL roundtrip = %v, %v", got, err)
	}

	base62Text := Base62Encode(data)
	if got, err := Base62Decode(base62Text); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("Base62 roundtrip = %v, %v", got, err)
	}

	base58Text := Base58Encode(data)
	if got, err := Base58Decode(base58Text); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("Base58 roundtrip = %v, %v", got, err)
	}
}

// TestBaseDecodeInvalidCharacter 验证 Base58/Base62 非法字符会返回错误。
func TestBaseDecodeInvalidCharacter(t *testing.T) {
	if _, err := Base58Decode("0"); err == nil {
		t.Fatal("Base58Decode(invalid) error = nil")
	}
	if _, err := Base62Decode("*"); err == nil {
		t.Fatal("Base62Decode(invalid) error = nil")
	}
}

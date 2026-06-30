package base64x

import (
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

const (
	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Base64Encode 进行 base64 编码。
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 进行 base64 解码。
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64URLEncode 进行 URL 安全 base64 编码。
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URLDecode 进行 URL 安全 base64 解码。
func Base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// Base62Encode 进行 base62 编码。
func Base62Encode(data []byte) string {
	return encodeBig(data, base62Alphabet)
}

// Base62Decode 进行 base62 解码。
func Base62Decode(s string) ([]byte, error) {
	return decodeBig(s, base62Alphabet)
}

// Base58Encode 进行 base58 编码。
func Base58Encode(data []byte) string {
	return encodeBig(data, base58Alphabet)
}

// Base58Decode 进行 base58 解码。
func Base58Decode(s string) ([]byte, error) {
	return decodeBig(s, base58Alphabet)
}

// encodeBig 使用指定字符表编码字节切片。
func encodeBig(data []byte, alphabet string) string {
	if len(data) == 0 {
		return ""
	}

	x := new(big.Int).SetBytes(data)
	base := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)
	mod := new(big.Int)
	var result []byte

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}
	for _, b := range data {
		if b != 0 {
			break
		}
		result = append(result, alphabet[0])
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

// decodeBig 使用指定字符表解码字符串。
func decodeBig(s string, alphabet string) ([]byte, error) {
	if s == "" {
		return []byte{}, nil
	}

	result := big.NewInt(0)
	base := big.NewInt(int64(len(alphabet)))
	for _, r := range s {
		index := strings.IndexRune(alphabet, r)
		if index < 0 {
			return nil, errors.New("base64x: invalid character")
		}
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(index)))
	}

	data := result.Bytes()
	for _, r := range s {
		if byte(r) != alphabet[0] {
			break
		}
		data = append([]byte{0}, data...)
	}
	return data, nil
}

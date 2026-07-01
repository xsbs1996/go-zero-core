package xhmac

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// SHA1Hex 计算 HMAC-SHA1 并返回十六进制字符串。
func SHA1Hex(data string, key string) string {
	return sumHex(data, key, sha1.New)
}

// SHA256Hex 计算 HMAC-SHA256 并返回十六进制字符串。
func SHA256Hex(data string, key string) string {
	return sumHex(data, key, sha256.New)
}

// SHA512Hex 计算 HMAC-SHA512 并返回十六进制字符串。
func SHA512Hex(data string, key string) string {
	return sumHex(data, key, sha512.New)
}

// SHA256Base64 计算 HMAC-SHA256 并返回 base64 字符串。
func SHA256Base64(data string, key string) string {
	return sumBase64(data, key, sha256.New)
}

// VerifySHA256Hex 校验 HMAC-SHA256 十六进制签名。
func VerifySHA256Hex(data string, key string, sign string) bool {
	return hmac.Equal([]byte(SHA256Hex(data, key)), []byte(sign))
}

// sumHex 计算 HMAC 并返回十六进制字符串。
func sumHex(data string, key string, fn func() hash.Hash) string {
	mac := hmac.New(fn, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// sumBase64 计算 HMAC 并返回 base64 字符串。
func sumBase64(data string, key string, fn func() hash.Hash) string {
	mac := hmac.New(fn, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

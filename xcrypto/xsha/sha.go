package xsha

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// SHA1 计算 SHA1 摘要。
//
// 参数：
//   - s: 待计算字符串。
//
// 返回值：
//   - string: 十六进制 SHA1 摘要。
func SHA1(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// SHA256 计算 SHA256 摘要。
//
// 参数：
//   - s: 待计算字符串。
//
// 返回值：
//   - string: 十六进制 SHA256 摘要。
func SHA256(s string) string {
	return SHA256Bytes([]byte(s))
}

// SHA256Bytes 计算字节切片 SHA256 摘要。
//
// 参数：
//   - data: 待计算字节切片。
//
// 返回值：
//   - string: 十六进制 SHA256 摘要。
func SHA256Bytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// SHA384 计算 SHA384 摘要。
//
// 参数：
//   - s: 待计算字符串。
//
// 返回值：
//   - string: 十六进制 SHA384 摘要。
func SHA384(s string) string {
	sum := sha512.Sum384([]byte(s))
	return hex.EncodeToString(sum[:])
}

// SHA512 计算 SHA512 摘要。
//
// 参数：
//   - s: 待计算字符串。
//
// 返回值：
//   - string: 十六进制 SHA512 摘要。
func SHA512(s string) string {
	sum := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sum[:])
}

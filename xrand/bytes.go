package xrand

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// Bytes 生成指定长度的加密安全随机字节。
//
// 参数：
//   - length: 随机字节长度，必须大于等于 0。
//
// 返回值：
//   - []byte: 随机字节切片。
//   - error: length 非法或随机源读取失败时返回错误。
func Bytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("xrand: invalid length")
	}
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

// Hex 生成指定字节长度的十六进制随机字符串。
//
// 参数：
//   - length: 随机字节长度，最终字符串长度为 length * 2。
//
// 返回值：
//   - string: 十六进制随机字符串。
//   - error: length 非法或随机源读取失败时返回错误。
func Hex(length int) (string, error) {
	buf, err := Bytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

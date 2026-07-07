package xrand

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	NumberLetters = "0123456789"                                                     // NumberLetters 数字字符集。
	AlphaLetters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"           // AlphaLetters 字母字符集。
	MixedLetters  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // MixedLetters 数字字母字符集。
)

// String 使用指定字符集生成加密安全随机字符串。
//
// 参数：
//   - length: 随机字符串长度，必须大于等于 0。
//   - letters: 随机字符集，不能为空。
//
// 返回值：
//   - string: 随机字符串。
//   - error: length 非法、字符集为空或随机源读取失败时返回错误。
func String(length int, letters string) (string, error) {
	if length < 0 {
		return "", errors.New("xrand: invalid length")
	}
	if letters == "" {
		return "", errors.New("xrand: empty letters")
	}

	result := make([]byte, length)
	newInt := big.NewInt(int64(len(letters)))
	for i := range result {
		n, err := rand.Int(rand.Reader, newInt)
		if err != nil {
			return "", err
		}
		result[i] = letters[n.Int64()]
	}
	return string(result), nil
}

// Number 生成数字随机字符串。
//
// 参数：
//   - length: 随机字符串长度。
//
// 返回值：
//   - string: 仅包含数字的随机字符串。
//   - error: length 非法或随机源读取失败时返回错误。
func Number(length int) (string, error) {
	return String(length, NumberLetters)
}

// Alpha 生成字母随机字符串。
//
// 参数：
//   - length: 随机字符串长度。
//
// 返回值：
//   - string: 仅包含大小写字母的随机字符串。
//   - error: length 非法或随机源读取失败时返回错误。
func Alpha(length int) (string, error) {
	return String(length, AlphaLetters)
}

// Mixed 生成数字字母随机字符串。
//
// 参数：
//   - length: 随机字符串长度。
//
// 返回值：
//   - string: 包含数字和大小写字母的随机字符串。
//   - error: length 非法或随机源读取失败时返回错误。
func Mixed(length int) (string, error) {
	return String(length, MixedLetters)
}

package xrand

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
)

const (
	NumberLetters = "0123456789"                                                     // NumberLetters 数字字符集。
	AlphaLetters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"           // AlphaLetters 字母字符集。
	MixedLetters  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // MixedLetters 数字字母字符集。
)

// Bytes 生成指定长度随机字节。
func Bytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("xrand: invalid length")
	}
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

// Hex 生成指定字节长度的十六进制随机字符串。
func Hex(length int) (string, error) {
	buf, err := Bytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// String 使用指定字符集生成随机字符串。
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
func Number(length int) (string, error) {
	return String(length, NumberLetters)
}

// Alpha 生成字母随机字符串。
func Alpha(length int) (string, error) {
	return String(length, AlphaLetters)
}

// Mixed 生成数字字母随机字符串。
func Mixed(length int) (string, error) {
	return String(length, MixedLetters)
}

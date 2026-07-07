package xmd5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// Sum 计算字符串 MD5 摘要。
//
// 参数：
//   - s: 待计算字符串。
//
// 返回值：
//   - string: 十六进制 MD5 摘要。
func Sum(s string) string {
	return SumBytes([]byte(s))
}

// SumBytes 计算字节切片 MD5 摘要。
//
// 参数：
//   - data: 待计算字节切片。
//
// 返回值：
//   - string: 十六进制 MD5 摘要。
func SumBytes(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// SumFile 计算文件 MD5 摘要。
//
// 参数：
//   - path: 文件路径。
//
// 返回值：
//   - string: 十六进制 MD5 摘要。
//   - error: 文件打开或读取失败时返回错误。
func SumFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

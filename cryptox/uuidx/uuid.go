package uuidx

import "github.com/google/uuid"

// New 生成 UUID 字符串。
func New() string {
	return uuid.NewString()
}

// NewWithoutDash 生成不带横线的 UUID 字符串。
func NewWithoutDash() string {
	id := uuid.New()
	buf := make([]byte, 32)
	hex := "0123456789abcdef"
	for i, b := range id {
		buf[i*2] = hex[b>>4]
		buf[i*2+1] = hex[b&0x0f]
	}
	return string(buf)
}

// Parse 解析 UUID 字符串。
func Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValid 判断字符串是否为合法 UUID。
func IsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

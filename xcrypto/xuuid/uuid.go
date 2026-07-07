package xuuid

import "github.com/google/uuid"

// New 生成 UUID 字符串。
//
// 返回值：
//   - string: 标准带横线 UUID 字符串。
func New() string {
	return uuid.NewString()
}

// NewWithoutDash 生成不带横线的 UUID 字符串。
//
// 返回值：
//   - string: 32 位不带横线 UUID 字符串。
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
//
// 参数：
//   - s: UUID 字符串，支持带横线和不带横线格式。
//
// 返回值：
//   - uuid.UUID: 解析后的 UUID。
//   - error: 输入不是合法 UUID 时返回错误。
func Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValid 判断字符串是否为合法 UUID。
//
// 参数：
//   - s: 待校验字符串。
//
// 返回值：
//   - bool: true 表示字符串是合法 UUID。
func IsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

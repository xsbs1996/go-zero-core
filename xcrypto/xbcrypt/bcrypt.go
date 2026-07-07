package xbcrypt

import "golang.org/x/crypto/bcrypt"

// Hash 对密码进行 bcrypt 哈希。
//
// 参数：
//   - password: 明文密码。
//
// 返回值：
//   - string: bcrypt 哈希字符串。
//   - error: 哈希生成失败时返回错误。
func Hash(password string) (string, error) {
	return HashWithCost(password, bcrypt.DefaultCost)
}

// HashWithCost 使用指定 cost 对密码进行 bcrypt 哈希。
//
// 参数：
//   - password: 明文密码。
//   - cost: bcrypt cost，必须处于 bcrypt 支持范围内。
//
// 返回值：
//   - string: bcrypt 哈希字符串。
//   - error: cost 非法或哈希生成失败时返回错误。
func HashWithCost(password string, cost int) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Compare 校验密码和 bcrypt 哈希是否匹配。
//
// 参数：
//   - password: 明文密码。
//   - hashed: bcrypt 哈希字符串。
//
// 返回值：
//   - bool: true 表示密码匹配，false 表示不匹配或哈希非法。
func Compare(password string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

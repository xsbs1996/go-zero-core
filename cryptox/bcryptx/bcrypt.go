package bcryptx

import "golang.org/x/crypto/bcrypt"

// Hash 对密码进行 bcrypt 哈希。
func Hash(password string) (string, error) {
	return HashWithCost(password, bcrypt.DefaultCost)
}

// HashWithCost 使用指定 cost 对密码进行 bcrypt 哈希。
func HashWithCost(password string, cost int) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Compare 校验密码和 bcrypt 哈希是否匹配。
func Compare(password string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

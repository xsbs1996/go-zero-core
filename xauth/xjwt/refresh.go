package xjwt

// Refresh 刷新 JWT token。
func Refresh(conf Config, tokenString string) (string, error) {
	claims, err := Parse(conf, tokenString)
	if err != nil {
		return "", err
	}
	return Generate(conf, claims.Data)
}

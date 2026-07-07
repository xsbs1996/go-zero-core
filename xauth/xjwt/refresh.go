package xjwt

// Refresh 刷新 JWT token。
//
// 参数：
//   - conf: JWT 解析和重新签发配置。
//   - tokenString: 需要刷新的 JWT token。
//
// 返回值：
//   - string: 重新签发后的 JWT token。
//   - error: 原 token 解析失败或重新签名失败时返回错误。
func Refresh(conf Config, tokenString string) (string, error) {
	claims, err := Parse(conf, tokenString)
	if err != nil {
		return "", err
	}
	return Generate(conf, claims.Data)
}

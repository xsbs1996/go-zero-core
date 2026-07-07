package xjwt

import "github.com/golang-jwt/jwt/v5"

// Parse 解析 JWT token。
//
// 参数：
//   - conf: JWT 校验配置，必须包含签名密钥；Issuer/Audience 非空时会参与校验。
//   - tokenString: 待解析的 JWT token 字符串。
//
// 返回值：
//   - *Claims: 解析成功后的 JWT 载荷。
//   - error: 密钥为空、签名算法不匹配、token 无效或校验失败时返回错误。
func Parse(conf Config, tokenString string) (*Claims, error) {
	if conf.Secret == "" {
		return nil, ErrMissingSecret
	}

	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(conf.Secret), nil
	}, jwt.WithAudience(conf.Audience...), jwt.WithIssuer(conf.Issuer))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

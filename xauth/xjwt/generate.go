package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generate 生成 JWT token。
//
// 参数：
//   - conf: JWT 签名、签发者、受众和过期时间配置。
//   - data: 业务载荷，会写入 Claims.Data。
//
// 返回值：
//   - string: 签名后的 JWT token。
//   - error: 密钥为空或签名失败时返回错误。
func Generate(conf Config, data map[string]any) (string, error) {
	if conf.Secret == "" {
		return "", ErrMissingSecret
	}

	now := time.Now()
	claims := Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    conf.Issuer,
			Subject:   conf.Subject,
			Audience:  conf.Audience,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	if conf.Expire > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(conf.Expire))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.Secret))
}

package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generate 生成 JWT token。
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

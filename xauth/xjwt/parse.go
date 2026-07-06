package xjwt

import "github.com/golang-jwt/jwt/v5"

// Parse 解析 JWT token。
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

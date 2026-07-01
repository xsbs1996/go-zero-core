package jwtx

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingSecret = errors.New("jwtx: missing secret") // ErrMissingSecret 表示 JWT 密钥为空。
	ErrInvalidToken  = errors.New("jwtx: invalid token")  // ErrInvalidToken 表示 JWT token 不合法。
)

// Config JWT 配置。
type Config struct {
	Secret        string        `json:"secret" yaml:"secret"`                        // Secret JWT 签名密钥。
	Issuer        string        `json:"issuer,optional" yaml:"issuer"`               // Issuer 签发者。
	Subject       string        `json:"subject,optional" yaml:"subject"`             // Subject 主题。
	Audience      []string      `json:"audience,optional" yaml:"audience"`           // Audience 受众。
	Expire        time.Duration `json:"expire,optional" yaml:"expire"`               // Expire 有效期。
	RefreshExpire time.Duration `json:"refreshExpire,optional" yaml:"refreshExpire"` // RefreshExpire 刷新有效期。
}

// Claims JWT 载荷。
type Claims struct {
	Data map[string]any `json:"data"` // Data 业务数据。
	jwt.RegisteredClaims
}

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

// Refresh 刷新 JWT token。
func Refresh(conf Config, tokenString string) (string, error) {
	claims, err := Parse(conf, tokenString)
	if err != nil {
		return "", err
	}
	return Generate(conf, claims.Data)
}

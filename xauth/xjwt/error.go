package xjwt

import "errors"

var (
	ErrMissingSecret = errors.New("xjwt: missing secret") // ErrMissingSecret 表示 JWT 密钥为空。
	ErrInvalidToken  = errors.New("xjwt: invalid token")  // ErrInvalidToken 表示 JWT token 不合法。
)

package wsx

import "errors"

var (
	ErrInvalidCode     = errors.New("wsx: invalid code")            // ErrInvalidCode 表示连接编码无效
	ErrMaxConnExceeded = errors.New("wsx: max connection exceeded") // ErrMaxConnExceeded 表示连接数超过上限
	ErrSessionNotFound = errors.New("wsx: session not found")       // ErrSessionNotFound 表示会话不存在
)

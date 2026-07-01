package xws

import "errors"

var (
	ErrInvalidCode     = errors.New("xws: invalid code")            // ErrInvalidCode 表示连接编码无效
	ErrMaxConnExceeded = errors.New("xws: max connection exceeded") // ErrMaxConnExceeded 表示连接数超过上限
	ErrSessionNotFound = errors.New("xws: session not found")       // ErrSessionNotFound 表示会话不存在
)

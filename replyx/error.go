package replyx

import "errors"

var (
	ErrReservedCode  = errors.New("replyx: reserved code")  // ErrReservedCode 表示注册了保留区间内的错误码
	ErrDuplicateCode = errors.New("replyx: duplicate code") // ErrDuplicateCode 表示注册了已存在的错误码
	ErrEmptyMsg      = errors.New("replyx: empty msg")      // ErrEmptyMsg 表示注册错误码时传入了空消息
)

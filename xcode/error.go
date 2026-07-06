package xcode

import "errors"

var (
	ErrReservedCode  = errors.New("xcode: reserved code")  // ErrReservedCode 表示注册了保留区间内的错误码。
	ErrDuplicateCode = errors.New("xcode: duplicate code") // ErrDuplicateCode 表示注册了已存在的错误码。
	ErrEmptyMsg      = errors.New("xcode: empty msg")      // ErrEmptyMsg 表示注册错误码时传入了空消息。
)

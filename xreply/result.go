package xreply

// Result 表示统一 API 响应结构。
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// NewResult 创建统一 API 响应结构。
func NewResult[T any](code int, msg string, data T) Result[T] {
	return Result[T]{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

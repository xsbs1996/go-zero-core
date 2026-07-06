package xreply

import (
	"net/http"

	"github.com/xsbs1996/go-zero-core/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Success 输出成功 JSON 响应。
func Success[T any](w http.ResponseWriter, data T, vars ...xcode.Vars) {
	httpx.OkJson(w, NewResult(xcode.CodeSuccess, xcode.Msg(xcode.CodeSuccess, vars...), data))
}

// SuccessMsg 输出带自定义消息的成功 JSON 响应。
func SuccessMsg[T any](w http.ResponseWriter, code int, customMsg string, data T) {
	if customMsg == "" {
		customMsg = xcode.Msg(code)
	}
	httpx.OkJson(w, NewResult(code, customMsg, data))
}

// Fail 根据业务错误码输出失败 JSON 响应，HTTP 状态码固定为 200。
func Fail(w http.ResponseWriter, code int, vars ...xcode.Vars) {
	httpx.OkJson(w, NewResult[any](code, xcode.Msg(code, vars...), nil))
}

// FailStatus 根据 HTTP 状态码和业务错误码输出失败 JSON 响应。
func FailStatus(w http.ResponseWriter, status int, code int, vars ...xcode.Vars) {
	httpx.WriteJson(w, status, NewResult[any](code, xcode.Msg(code, vars...), nil))
}

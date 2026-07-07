package xreply

import (
	"net/http"

	"github.com/xsbs1996/go-zero-core/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Success 输出成功 JSON 响应。
//
// 参数：
//   - w: HTTP 响应写入器。
//   - data: 成功响应数据。
//   - vars: 可选文案模板变量，通常不需要传。
func Success[T any](w http.ResponseWriter, data T, vars ...xcode.Vars) {
	httpx.OkJson(w, NewResult(xcode.CodeSuccess, xcode.Msg(xcode.CodeSuccess, vars...), data))
}

// SuccessMsg 输出带自定义消息的成功 JSON 响应。
//
// 参数：
//   - w: HTTP 响应写入器。
//   - code: 业务错误码，允许使用非 CodeSuccess 的成功类业务码。
//   - customMsg: 自定义响应文案；为空时使用 xcode.Msg(code)。
//   - data: 响应数据。
func SuccessMsg[T any](w http.ResponseWriter, code int, customMsg string, data T) {
	if customMsg == "" {
		customMsg = xcode.Msg(code)
	}
	httpx.OkJson(w, NewResult(code, customMsg, data))
}

// Fail 根据业务错误码输出失败 JSON 响应，HTTP 状态码固定为 200。
//
// 参数：
//   - w: HTTP 响应写入器。
//   - code: 业务错误码。
//   - vars: 可选文案模板变量。
func Fail(w http.ResponseWriter, code int, vars ...xcode.Vars) {
	httpx.OkJson(w, NewResult[any](code, xcode.Msg(code, vars...), nil))
}

// FailStatus 根据 HTTP 状态码和业务错误码输出失败 JSON 响应。
//
// 参数：
//   - w: HTTP 响应写入器。
//   - status: HTTP 状态码。
//   - code: 业务错误码。
//   - vars: 可选文案模板变量。
func FailStatus(w http.ResponseWriter, status int, code int, vars ...xcode.Vars) {
	httpx.WriteJson(w, status, NewResult[any](code, xcode.Msg(code, vars...), nil))
}

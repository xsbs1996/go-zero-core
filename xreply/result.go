package xreply

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type page[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page,omitempty"`
	PageSize int   `json:"pageSize,omitempty"`
}

func newResult[T any](code int, msg string, data T) result[T] {
	return result[T]{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Success 输出成功 JSON 响应
func Success[T any](w http.ResponseWriter, data T, vars ...Vars) {
	httpx.OkJson(w, newResult(CodeSuccess, msg(CodeSuccess, vars...), data))
}

// Fail 根据业务错误码输出失败 JSON 响应，HTTP 状态码固定为 200
func Fail(w http.ResponseWriter, code int, vars ...Vars) {
	httpx.OkJson(w, newResult[any](code, msg(code, vars...), nil))
}

// FailStatus 根据 HTTP 状态码和业务错误码输出失败 JSON 响应
func FailStatus(w http.ResponseWriter, status int, code int, vars ...Vars) {
	httpx.WriteJson(w, status, newResult[any](code, msg(code, vars...), nil))
}

// SuccessPage 输出分页成功 JSON 响应
func SuccessPage[T any](w http.ResponseWriter, list []T, total int64, pageNo int, pageSize int) {
	httpx.OkJson(w, successPageResult(list, total, pageNo, pageSize))
}

// SuccessMsg 输出带自定义消息的成功 JSON 响应
func SuccessMsg[T any](w http.ResponseWriter, code int, customMsg string, data T) {
	if customMsg == "" {
		customMsg = msg(code)
	}
	httpx.OkJson(w, newResult(code, customMsg, data))
}

func successPageResult[T any](list []T, total int64, pageNo int, pageSize int) result[page[T]] {
	return newResult(CodeSuccess, MsgSuccess, page[T]{
		List:     list,
		Total:    total,
		Page:     pageNo,
		PageSize: pageSize,
	})
}

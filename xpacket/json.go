package xpacket

import (
	"encoding/json"
	"fmt"

	"github.com/xsbs1996/go-zero-core/xcode"
	"github.com/xsbs1996/go-zero-core/xreply"
)

// EncodeJson 将 action 和 JSON 载荷编码为默认二进制封包格式。
//
// v 为 nil 时会编码 JSON null。该函数只负责普通 JSON 载荷，不主动追加 code/msg。
// 如果需要和 xreply 一致的 {code,msg,data} 响应结构，请使用 EncodeJsonResult。
func EncodeJson(action uint16, v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("xpacket: marshal json: %w", err)
	}
	return EncodeBody(action, body)
}

// DecodeJson 将默认二进制封包格式解码到 JSON 目标值，并返回包头中的 action。
//
// v 必须是可写入的指针。解码完成后，如果 v 实现 JsonBusinessValidator，
// 会继续执行 ValidateBusinessJson 做 JSON 载荷业务校验。
func DecodeJson(data []byte, v any) (action uint16, err error) {
	if v == nil {
		return 0, ErrNilMessage
	}

	action, body, err := DecodeBody(data)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return action, fmt.Errorf("xpacket: unmarshal json: %w", err)
	}

	if err := validateBusinessJson(v); err != nil {
		return action, err
	}

	return action, nil
}

// EncodeJsonResult 将 action、业务码和 data 编码为 xreply 统一响应结构的 JSON 封包。
//
// body JSON 结构与 xreply.Success/xreply.Fail 保持一致：
//
//	{"code":0,"msg":"success","data":...}
//
// msg 由 xcode.Msg 根据 code 和 vars 渲染，错误码仍由 xcode 统一管理。
func EncodeJsonResult[T any](action uint16, code int, data T, vars ...xcode.Vars) ([]byte, error) {
	return EncodeJson(action, xreply.NewResult(code, xcode.Msg(code, vars...), data))
}

// DecodeJsonResult 将默认二进制封包格式解码为 xreply 统一响应结构。
//
// 该函数适合接收 EncodeJsonResult 生成的包；错误码和文案字段仍沿用
// xreply.Result 的 code/msg/data 结构。
func DecodeJsonResult[T any](data []byte) (action uint16, result xreply.Result[T], err error) {
	action, err = DecodeJson(data, &result)
	if err != nil {
		return 0, result, err
	}
	return action, result, nil
}

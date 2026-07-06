// Package xcode 提供跨 HTTP、WebSocket 等传输层复用的统一业务错误码。
package xcode

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const minCustomCode = 100

var unresolvedPlaceholderPattern = regexp.MustCompile(`\s*[:：,，;；-]?\s*\{[^{}]+\}`)

const (
	CodeSuccess            = 0  // CodeSuccess 表示请求处理成功。
	CodeError              = 1  // CodeError 表示通用失败。
	CodeInvalidParam       = 2  // CodeInvalidParam 表示请求参数无效，默认文案中的 {field} 可替换为具体字段名。
	CodeUnauthorized       = 3  // CodeUnauthorized 表示请求未认证。
	CodeForbidden          = 4  // CodeForbidden 表示请求已认证但无访问权限。
	CodeNotFound           = 5  // CodeNotFound 表示请求的资源不存在。
	CodeConflict           = 6  // CodeConflict 表示请求与当前资源状态冲突。
	CodeTooManyRequests    = 7  // CodeTooManyRequests 表示请求过于频繁。
	CodeInvalidToken       = 8  // CodeInvalidToken 表示 token 无效。
	CodeTokenExpired       = 9  // CodeTokenExpired 表示 token 已过期。
	CodeInvalidSign        = 10 // CodeInvalidSign 表示签名无效。
	CodeTimeout            = 11 // CodeTimeout 表示请求处理超时。
	CodeServiceBusy        = 12 // CodeServiceBusy 表示服务繁忙。
	CodeServiceUnavailable = 13 // CodeServiceUnavailable 表示服务暂不可用。
	CodeInternal           = 99 // CodeInternal 表示内部错误。
)

const (
	CodeInsufficientBalance = 100 // CodeInsufficientBalance 表示账户余额不足。
)

const (
	MsgSuccess             = "success"                // MsgSuccess 表示成功默认文案。
	MsgError               = "error"                  // MsgError 表示通用失败默认文案。
	MsgInvalidParam        = "invalid param: {field}" // MsgInvalidParam 表示参数错误默认文案，{field} 可替换为字段名。
	MsgUnauthorized        = "unauthorized"           // MsgUnauthorized 表示未认证默认文案。
	MsgForbidden           = "forbidden"              // MsgForbidden 表示无权限默认文案。
	MsgNotFound            = "not found"              // MsgNotFound 表示资源不存在默认文案。
	MsgConflict            = "conflict"               // MsgConflict 表示资源冲突默认文案。
	MsgTooManyRequests     = "too many requests"      // MsgTooManyRequests 表示请求过于频繁默认文案。
	MsgInvalidToken        = "invalid token"          // MsgInvalidToken 表示 token 无效默认文案。
	MsgTokenExpired        = "token expired"          // MsgTokenExpired 表示 token 过期默认文案。
	MsgInvalidSign         = "invalid sign"           // MsgInvalidSign 表示签名无效默认文案。
	MsgTimeout             = "timeout"                // MsgTimeout 表示超时默认文案。
	MsgServiceBusy         = "service busy"           // MsgServiceBusy 表示服务繁忙默认文案。
	MsgServiceUnavailable  = "service unavailable"    // MsgServiceUnavailable 表示服务不可用默认文案。
	MsgInternal            = "internal error"         // MsgInternal 表示内部错误默认文案。
	MsgInsufficientBalance = "insufficient balance"   // MsgInsufficientBalance 表示余额不足默认文案。
	msgUnknownCode         = "unknown error"
)

var (
	codeMu sync.RWMutex
	codes  = map[int]string{
		CodeSuccess:             MsgSuccess,
		CodeError:               MsgError,
		CodeInvalidParam:        MsgInvalidParam,
		CodeUnauthorized:        MsgUnauthorized,
		CodeForbidden:           MsgForbidden,
		CodeNotFound:            MsgNotFound,
		CodeConflict:            MsgConflict,
		CodeTooManyRequests:     MsgTooManyRequests,
		CodeInvalidToken:        MsgInvalidToken,
		CodeTokenExpired:        MsgTokenExpired,
		CodeInvalidSign:         MsgInvalidSign,
		CodeTimeout:             MsgTimeout,
		CodeServiceBusy:         MsgServiceBusy,
		CodeServiceUnavailable:  MsgServiceUnavailable,
		CodeInternal:            MsgInternal,
		CodeInsufficientBalance: MsgInsufficientBalance,
	}
)

// Vars 表示用于渲染错误文案模板的命名变量。
type Vars map[string]any

// RegisterCodes 将业务项目自定义错误码合并到统一错误码表。
//
// 0-99 为基础库通用错误码保留区间；100 及以上可用于业务错误码。
// 如果业务错误码与库内已登记错误码冲突，会返回 ErrDuplicateCode。
func RegisterCodes(items map[int]string) error {
	codeMu.Lock()
	defer codeMu.Unlock()

	for code, msg := range items {
		if code < minCustomCode {
			return fmt.Errorf("%w: %d", ErrReservedCode, code)
		}
		if strings.TrimSpace(msg) == "" {
			return fmt.Errorf("%w: %d", ErrEmptyMsg, code)
		}
		if _, ok := codes[code]; ok {
			return fmt.Errorf("%w: %d", ErrDuplicateCode, code)
		}
	}

	for code, msg := range items {
		codes[code] = msg
	}
	return nil
}

// Msg 返回错误码对应的文案，并使用 vars 渲染模板变量。
func Msg(code int, vars ...Vars) string {
	codeMu.RLock()
	template, ok := codes[code]
	codeMu.RUnlock()
	if !ok {
		template = msgUnknownCode
	}
	return renderMsg(template, vars...)
}

// CodeMsgMap 返回当前错误码和文案映射的副本。
func CodeMsgMap() map[int]string {
	codeMu.RLock()
	defer codeMu.RUnlock()

	items := make(map[int]string, len(codes))
	for code, msg := range codes {
		items[code] = msg
	}
	return items
}

func renderMsg(template string, vars ...Vars) string {
	if len(vars) == 0 {
		return cleanMsg(template)
	}

	for _, values := range vars {
		for key, value := range values {
			if key == "" {
				continue
			}
			template = strings.ReplaceAll(template, "{"+key+"}", fmt.Sprint(value))
		}
	}
	return cleanMsg(template)
}

func cleanMsg(msg string) string {
	msg = unresolvedPlaceholderPattern.ReplaceAllString(msg, "")
	return strings.TrimSpace(msg)
}

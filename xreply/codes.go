package xreply

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const minCustomCode = 100

var unresolvedPlaceholderPattern = regexp.MustCompile(`\s*[:：,，;；-]?\s*\{[^{}]+\}`)

const (
	CodeSuccess            = 0  // CodeSuccess 表示请求处理成功
	CodeError              = 1  // CodeError 表示通用失败
	CodeInvalidParam       = 2  // CodeInvalidParam 表示请求参数无效，默认文案中的 {field} 可替换为具体字段名
	CodeUnauthorized       = 3  // CodeUnauthorized 表示请求未认证
	CodeForbidden          = 4  // CodeForbidden 表示请求已认证但无访问权限
	CodeNotFound           = 5  // CodeNotFound 表示请求的资源不存在
	CodeConflict           = 6  // CodeConflict 表示请求与当前资源状态冲突
	CodeTooManyRequests    = 7  // CodeTooManyRequests 表示请求过于频繁
	CodeInvalidToken       = 8  // CodeInvalidToken 表示 token 无效
	CodeTokenExpired       = 9  // CodeTokenExpired 表示 token 已过期
	CodeInvalidSign        = 10 // CodeInvalidSign 表示签名无效
	CodeTimeout            = 11 // CodeTimeout 表示请求处理超时
	CodeServiceBusy        = 12 // CodeServiceBusy 表示服务繁忙
	CodeServiceUnavailable = 13 // CodeServiceUnavailable 表示服务暂不可用
	CodeInternal           = 99 // CodeInternal 表示内部错误
)

const (
	MsgSuccess            = "success"
	MsgError              = "error"
	MsgInvalidParam       = "invalid param: {field}"
	MsgUnauthorized       = "unauthorized"
	MsgForbidden          = "forbidden"
	MsgNotFound           = "not found"
	MsgConflict           = "conflict"
	MsgTooManyRequests    = "too many requests"
	MsgInvalidToken       = "invalid token"
	MsgTokenExpired       = "token expired"
	MsgInvalidSign        = "invalid sign"
	MsgTimeout            = "timeout"
	MsgServiceBusy        = "service busy"
	MsgServiceUnavailable = "service unavailable"
	MsgInternal           = "internal error"
	msgUnknownCode        = "unknown error"
)

var (
	codeMu sync.RWMutex
	codes  = map[int]string{
		CodeSuccess:            MsgSuccess,
		CodeError:              MsgError,
		CodeInvalidParam:       MsgInvalidParam,
		CodeUnauthorized:       MsgUnauthorized,
		CodeForbidden:          MsgForbidden,
		CodeNotFound:           MsgNotFound,
		CodeConflict:           MsgConflict,
		CodeTooManyRequests:    MsgTooManyRequests,
		CodeInvalidToken:       MsgInvalidToken,
		CodeTokenExpired:       MsgTokenExpired,
		CodeInvalidSign:        MsgInvalidSign,
		CodeTimeout:            MsgTimeout,
		CodeServiceBusy:        MsgServiceBusy,
		CodeServiceUnavailable: MsgServiceUnavailable,
		CodeInternal:           MsgInternal,
	}
)

// Vars 表示用于渲染 msg 模板的命名变量
type Vars map[string]any

// RegisterCodes 将业务项目自定义错误码合并到默认错误码表
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

func msg(code int, vars ...Vars) string {
	codeMu.RLock()
	template, ok := codes[code]
	codeMu.RUnlock()
	if !ok {
		template = msgUnknownCode
	}
	return renderMsg(template, vars...)
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

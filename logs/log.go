package logs

import (
	"context"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
)

// Debug 输出 debug 日志。
func Debug(ctx context.Context, msg string, content any, err any) {
	if IsDisabled(ctx) {
		return
	}
	DebugContent(ctx, Content{Msg: msg, Content: content, Error: err})
}

// DebugContent 输出 debug 日志。
func DebugContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Debug(BodyContent(content))
}

// Info 输出 info 日志。
func Info(ctx context.Context, msg string, content any) {
	if IsDisabled(ctx) {
		return
	}
	InfoContent(ctx, Content{Msg: msg, Content: content})
}

// InfoContent 输出 info 日志。
func InfoContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Info(BodyContent(content))
}

// Warn 输出 warn 日志。
func Warn(ctx context.Context, msg string, content any, err any) {
	if IsDisabled(ctx) {
		return
	}
	WarnContent(ctx, Content{Msg: msg, Content: content, Error: err})
}

// WarnContent 输出 warn 日志。
func WarnContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Sloww(BodyContent(content), logx.Field("severity", "warn"))
}

// Error 输出 error 日志。
func Error(ctx context.Context, msg string, content any, err any) {
	if IsDisabled(ctx) {
		return
	}
	ErrorContent(ctx, Content{Msg: msg, Content: content, Error: err})
}

// ErrorContent 输出 error 日志。
func ErrorContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Error(BodyContent(content))
}

// ErrorStack 输出带堆栈的 error 日志。
func ErrorStack(ctx context.Context, msg string, content any, err any) {
	if IsDisabled(ctx) {
		return
	}
	ErrorStackContent(ctx, Content{Msg: msg, Content: content, Error: err})
}

// ErrorStackContent 输出带堆栈的 error 日志。
func ErrorStackContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	body := content
	body.Content = map[string]any{
		"data":  normalizeBody(content.Content),
		"stack": string(debug.Stack()),
	}

	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Error(BodyContent(body))
}

// Severe 输出 severe 日志。
func Severe(ctx context.Context, msg string, content any, err any) {
	if IsDisabled(ctx) {
		return
	}
	SevereContent(ctx, Content{Msg: msg, Content: content, Error: err})
}

// SevereContent 输出 severe 日志。
func SevereContent(ctx context.Context, content Content) {
	if IsDisabled(ctx) {
		return
	}
	logx.WithContext(ContextWithTrace(ctx)).WithCallerSkip(2).Errorw(BodyContent(content), logx.Field("severity", "severe"))
}

package xlog

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"go.opentelemetry.io/otel/trace"
)

var errGenerateTraceFailed = errors.New("xlog: generate trace failed")

// ContextWithTrace 复用上下文已有链路追踪信息，不存在时生成新的链路追踪信息。
//
// 参数：
//   - ctx: 父上下文；nil 会回退为 context.Background。
//
// 返回值：
//   - context.Context: 带 trace/span 信息的上下文。
func ContextWithTrace(ctx context.Context) context.Context {
	ctx = safeContext(ctx)

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() && spanCtx.HasSpanID() {
		return ctx
	}

	return contextWithNewTrace(ctx)
}

// ContextWithForceTrace 强制生成新的链路追踪信息并覆盖上下文。
//
// 参数：
//   - ctx: 父上下文；nil 会回退为 context.Background。
//
// 返回值：
//   - context.Context: 写入新 trace/span 信息的上下文。
func ContextWithForceTrace(ctx context.Context) context.Context {
	return contextWithNewTrace(safeContext(ctx))
}

// contextWithNewTrace 生成新的链路追踪信息并写入上下文。
func contextWithNewTrace(ctx context.Context) context.Context {
	traceID, err := randomTraceID()
	if err != nil {
		panic(err)
	}

	spanID, err := randomSpanID()
	if err != nil {
		panic(err)
	}

	newSpanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})

	return trace.ContextWithRemoteSpanContext(ctx, newSpanCtx)
}

// safeContext 返回安全的上下文，空上下文会回退为后台上下文。
func safeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

// randomTraceID 生成符合 OpenTelemetry 规范的 traceID。
func randomTraceID() (trace.TraceID, error) {
	id, err := randomHex(16)
	if err != nil {
		return trace.TraceID{}, err
	}

	return trace.TraceIDFromHex(id)
}

// randomSpanID 生成符合 OpenTelemetry 规范的 spanID。
func randomSpanID() (trace.SpanID, error) {
	id, err := randomHex(8)
	if err != nil {
		return trace.SpanID{}, err
	}

	return trace.SpanIDFromHex(id)
}

// randomHex 生成指定字节数的非零随机十六进制字符串。
func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	for {
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return "", errGenerateTraceFailed
		}
		if !allZero(buf) {
			return hex.EncodeToString(buf), nil
		}
	}
}

// allZero 判断字节切片是否全部为零。
func allZero(buf []byte) bool {
	for _, b := range buf {
		if b != 0 {
			return false
		}
	}

	return true
}

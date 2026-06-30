package logs

import "context"

type disableLogKey struct{}

// ContextWithDisable 返回禁用日志输出的上下文，仅影响当前上下文链路。
func ContextWithDisable(ctx context.Context) context.Context {
	return context.WithValue(safeContext(ctx), disableLogKey{}, true)
}

// IsDisabled 返回当前上下文是否已禁用日志输出。
func IsDisabled(ctx context.Context) bool {
	disabled, ok := safeContext(ctx).Value(disableLogKey{}).(bool)
	return ok && disabled
}

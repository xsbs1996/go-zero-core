package xmid

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

type ipContextKey struct{}

// IPConfig IP 地址中间件配置。
type IPConfig struct {
	Headers []string `json:"headers,optional" yaml:"headers"` // Headers 允许读取的代理 IP 请求头。
}

// IP 创建 go-zero REST IP 地址中间件。
//
// 参数：
//   - conf: IP 解析配置，主要包含代理 IP 请求头列表。
//
// 返回值：
//   - rest.Middleware: 可通过 go-zero server.Use 注册的客户端 IP 中间件。
func IP(conf IPConfig) rest.Middleware {
	conf = conf.WithDefault()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := ClientIP(r, conf.Headers)
			ctx := context.WithValue(r.Context(), ipContextKey{}, clientIP)
			next(w, r.WithContext(ctx))
		}
	}
}

// WithDefault 返回补齐默认值后的 IP 配置。
//
// 返回值：
//   - IPConfig: 已补齐默认代理请求头的配置。
func (c IPConfig) WithDefault() IPConfig {
	if len(c.Headers) == 0 {
		c.Headers = []string{"X-Forwarded-For", "X-Real-IP"}
	}
	return c
}

// IPFromContext 从上下文读取客户端 IP。
//
// 参数：
//   - ctx: 请求上下文。
//
// 返回值：
//   - string: IP 中间件写入的客户端 IP。
//   - bool: true 表示上下文中存在客户端 IP。
func IPFromContext(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(ipContextKey{}).(string)
	return ip, ok && ip != ""
}

// ClientIP 从请求中解析客户端 IP。
//
// 参数：
//   - r: HTTP 请求对象。
//   - headers: 优先读取的代理 IP 请求头列表，例如 X-Forwarded-For。
//
// 返回值：
//   - string: 解析到的客户端 IP；未命中代理头时回退为 RemoteAddr。
func ClientIP(r *http.Request, headers []string) string {
	for _, header := range headers {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}
		if header == "X-Forwarded-For" {
			parts := strings.Split(value, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}
		return value
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

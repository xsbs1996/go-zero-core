package xmid

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

type authContextKey struct{}

// AuthVerifyFunc 鉴权函数，返回写入上下文的认证信息。
type AuthVerifyFunc func(r *http.Request, token string) (any, error)

// AuthUnauthorizedFunc 鉴权失败响应函数。
type AuthUnauthorizedFunc func(w http.ResponseWriter, r *http.Request, err error)

// AuthConfig 鉴权中间件配置。
type AuthConfig struct {
	Header       string               `json:"header,default=Authorization" yaml:"header"` // Header token 所在请求头。
	Prefix       string               `json:"prefix,default=Bearer" yaml:"prefix"`        // Prefix token 前缀，例如 Bearer。
	SkipPaths    []string             `json:"skipPaths,optional" yaml:"skipPaths"`        // SkipPaths 跳过鉴权的路径。
	Verify       AuthVerifyFunc       `json:"-" yaml:"-"`                                 // Verify 鉴权函数。
	Unauthorized AuthUnauthorizedFunc `json:"-" yaml:"-"`                                 // Unauthorized 鉴权失败响应函数。
}

// Auth 创建 go-zero REST 鉴权中间件。
func Auth(conf AuthConfig) rest.Middleware {
	conf = conf.WithDefault()
	skipPaths := toSet(conf.SkipPaths)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := skipPaths[r.URL.Path]; ok {
				next(w, r)
				return
			}

			token := tokenFromHeader(r.Header.Get(conf.Header), conf.Prefix)
			if token == "" {
				conf.Unauthorized(w, r, ErrMissingToken)
				return
			}

			authInfo, err := conf.Verify(r, token)
			if err != nil {
				conf.Unauthorized(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), authContextKey{}, authInfo)
			next(w, r.WithContext(ctx))
		}
	}
}

// WithDefault 返回补齐默认值后的鉴权配置。
func (c AuthConfig) WithDefault() AuthConfig {
	if c.Header == "" {
		c.Header = "Authorization"
	}
	if c.Prefix == "" {
		c.Prefix = "Bearer"
	}
	if c.Unauthorized == nil {
		c.Unauthorized = defaultUnauthorized
	}
	return c
}

// AuthInfo 从上下文读取认证信息。
func AuthInfo(ctx context.Context) (any, bool) {
	info := ctx.Value(authContextKey{})
	return info, info != nil
}

// tokenFromHeader 从请求头提取 token。
func tokenFromHeader(value string, prefix string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if prefix == "" {
		return value
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], prefix) {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// toSet 将字符串切片转换为集合。
func toSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

package xmid

import (
	"net/http"
	"strings"

	"go-zero-core/xcast"

	"github.com/zeromicro/go-zero/rest"
)

// CorsConfig 跨域中间件配置。
type CorsConfig struct {
	AllowedOrigins   []string `json:"allowedOrigins,optional" yaml:"allowedOrigins"`     // AllowedOrigins 允许的来源。
	AllowedMethods   []string `json:"allowedMethods,optional" yaml:"allowedMethods"`     // AllowedMethods 允许的请求方法。
	AllowedHeaders   []string `json:"allowedHeaders,optional" yaml:"allowedHeaders"`     // AllowedHeaders 允许的请求头。
	ExposedHeaders   []string `json:"exposedHeaders,optional" yaml:"exposedHeaders"`     // ExposedHeaders 允许浏览器读取的响应头。
	AllowCredentials bool     `json:"allowCredentials,optional" yaml:"allowCredentials"` // AllowCredentials 是否允许携带凭证。
	MaxAge           int      `json:"maxAge,optional" yaml:"maxAge"`                     // MaxAge 预检请求缓存秒数。
}

// Cors 创建 go-zero REST 跨域中间件。
func Cors(conf CorsConfig) rest.Middleware {
	conf = conf.WithDefault()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				writeCorsHeaders(w.Header(), conf, origin)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}

// WithDefault 返回补齐默认值后的跨域配置。
func (c CorsConfig) WithDefault() CorsConfig {
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"*"}
	}
	if len(c.AllowedMethods) == 0 {
		c.AllowedMethods = []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}
	}
	if len(c.AllowedHeaders) == 0 {
		c.AllowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	return c
}

// writeCorsHeaders 写入跨域响应头。
func writeCorsHeaders(header http.Header, conf CorsConfig, origin string) {
	if allowOrigin := allowedOrigin(conf.AllowedOrigins, origin, conf.AllowCredentials); allowOrigin != "" {
		header.Set("Access-Control-Allow-Origin", allowOrigin)
	}
	header.Set("Access-Control-Allow-Methods", strings.Join(conf.AllowedMethods, ", "))
	header.Set("Access-Control-Allow-Headers", strings.Join(conf.AllowedHeaders, ", "))
	if len(conf.ExposedHeaders) > 0 {
		header.Set("Access-Control-Expose-Headers", strings.Join(conf.ExposedHeaders, ", "))
	}
	if conf.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if conf.MaxAge > 0 {
		header.Set("Access-Control-Max-Age", xcast.IntToString(conf.MaxAge))
	}
	header.Add("Vary", "Origin")
}

// allowedOrigin 返回允许写入响应头的来源。
func allowedOrigin(origins []string, origin string, allowCredentials bool) string {
	for _, item := range origins {
		if item == "*" {
			if allowCredentials {
				return origin
			}
			return "*"
		}
		if item == origin {
			return origin
		}
	}
	return ""
}

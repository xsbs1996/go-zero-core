package xmid

import (
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

// RateLimitKeyFunc 限流键生成函数。
type RateLimitKeyFunc func(r *http.Request) string

// RateLimitExceededFunc 限流响应函数。
type RateLimitExceededFunc func(w http.ResponseWriter, r *http.Request)

// RateLimitConfig 限流中间件配置。
type RateLimitConfig struct {
	Limit     int                   `json:"limit,default=100" yaml:"limit"`      // Limit 时间窗口内允许的最大请求数。
	Window    time.Duration         `json:"window,default=1m" yaml:"window"`     // Window 限流时间窗口，配置值示例：1m。
	KeyFunc   RateLimitKeyFunc      `json:"-" yaml:"-"`                          // KeyFunc 限流键生成函数。
	Exceeded  RateLimitExceededFunc `json:"-" yaml:"-"`                          // Exceeded 限流响应函数。
	Headers   []string              `json:"headers,optional" yaml:"headers"`     // Headers 解析客户端 IP 的请求头。
	SkipPaths []string              `json:"skipPaths,optional" yaml:"skipPaths"` // SkipPaths 跳过限流的路径。
}

type rateLimitItem struct {
	count     int
	expiresAt time.Time
}

// RateLimit 创建 go-zero REST 限流中间件。
func RateLimit(conf RateLimitConfig) rest.Middleware {
	conf = conf.WithDefault()
	store := make(map[string]rateLimitItem)
	var mu sync.Mutex
	skipPaths := toSet(conf.SkipPaths)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := skipPaths[r.URL.Path]; ok {
				next(w, r)
				return
			}

			key := conf.KeyFunc(r)
			now := time.Now()

			mu.Lock()
			item := store[key]
			if now.After(item.expiresAt) {
				item = rateLimitItem{expiresAt: now.Add(conf.Window)}
			}
			item.count++
			store[key] = item
			limited := item.count > conf.Limit
			cleanExpired(store, now)
			mu.Unlock()

			if limited {
				conf.Exceeded(w, r)
				return
			}

			next(w, r)
		}
	}
}

// WithDefault 返回补齐默认值后的限流配置。
func (c RateLimitConfig) WithDefault() RateLimitConfig {
	if c.Limit <= 0 {
		c.Limit = 100
	}
	if c.Window <= 0 {
		c.Window = time.Minute
	}
	if c.Exceeded == nil {
		c.Exceeded = defaultRateLimitExceeded
	}
	if c.KeyFunc == nil {
		headers := c.Headers
		c.KeyFunc = func(r *http.Request) string {
			return ClientIP(r, headers)
		}
	}
	return c
}

// cleanExpired 清理过期限流记录。
func cleanExpired(store map[string]rateLimitItem, now time.Time) {
	for key, item := range store {
		if now.After(item.expiresAt) {
			delete(store, key)
		}
	}
}

// defaultRateLimitExceeded 输出默认限流响应。
func defaultRateLimitExceeded(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

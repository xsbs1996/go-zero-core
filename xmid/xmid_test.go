package xmid

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestAuthMiddleware 验证鉴权中间件能提取 token 并写入认证上下文。
func TestAuthMiddleware(t *testing.T) {
	mw := Auth(AuthConfig{
		Verify: func(_ *http.Request, token string) (any, error) {
			if token != "abc" {
				return nil, errors.New("bad token")
			}
			return "user", nil
		},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer abc")
	mw(func(w http.ResponseWriter, r *http.Request) {
		info, ok := AuthInfo(r.Context())
		if !ok || info != "user" {
			t.Fatalf("AuthInfo() = %#v, %v", info, ok)
		}
		w.WriteHeader(http.StatusAccepted)
	})(w, r)
	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d", w.Code)
	}
}

// TestCorsMiddleware 验证 CORS 中间件对预检请求写入跨域响应头。
func TestCorsMiddleware(t *testing.T) {
	mw := Cors(CorsConfig{AllowCredentials: true})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/", nil)
	r.Header.Set("Origin", "https://example.com")

	mw(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next should not be called for OPTIONS")
	})(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("allow origin = %q", got)
	}
}

// TestIPMiddlewareAndClientIP 验证客户端 IP 解析和上下文写入。
func TestIPMiddlewareAndClientIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Forwarded-For", "1.1.1.1, 2.2.2.2")
	if got := ClientIP(r, []string{"X-Forwarded-For"}); got != "1.1.1.1" {
		t.Fatalf("ClientIP() = %q", got)
	}

	w := httptest.NewRecorder()
	IP(IPConfig{})(func(_ http.ResponseWriter, r *http.Request) {
		if ip, ok := IPFromContext(r.Context()); !ok || ip == "" {
			t.Fatalf("IPFromContext() = %q, %v", ip, ok)
		}
	})(w, r)
}

// TestRateLimitMiddleware 验证限流中间件在超过窗口次数后返回 429。
func TestRateLimitMiddleware(t *testing.T) {
	mw := RateLimit(RateLimitConfig{
		Limit:  1,
		Window: time.Minute,
		KeyFunc: func(*http.Request) string {
			return "same"
		},
	})
	handler := mw(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	w := httptest.NewRecorder()
	handler(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}

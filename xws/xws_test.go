package xws

import (
	"net/http"
	"testing"
	"time"
)

// TestDefaultAndNormalizeConfig 验证 WebSocket 默认配置和配置归一化。
func TestDefaultAndNormalizeConfig(t *testing.T) {
	conf := DefaultConfig()
	if conf.ReadBufferSize == 0 || conf.WriteBufferSize == 0 {
		t.Fatalf("DefaultConfig() = %#v", conf)
	}
	if conf.ReadDeadline != defaultReadDeadline || conf.WriteDeadline != defaultWriteDeadline {
		t.Fatalf("deadlines = %v/%v", conf.ReadDeadline, conf.WriteDeadline)
	}

	normalized := normalizeConfig(Config{
		ReadBufferSize:  1,
		WriteBufferSize: 2,
		ReadDeadline:    time.Second,
		CheckOrigin: func(*http.Request) bool {
			return false
		},
	})
	if normalized.ReadBufferSize != 1 || normalized.WriteBufferSize != 2 {
		t.Fatalf("normalizeConfig() = %#v", normalized)
	}
	if normalized.WriteDeadline != defaultWriteDeadline {
		t.Fatalf("WriteDeadline = %v", normalized.WriteDeadline)
	}
	if normalized.CheckOrigin(nil) {
		t.Fatal("custom CheckOrigin should be preserved")
	}
}

// TestManagerBasics 验证 WebSocket manager 的空状态行为。
func TestManagerBasics(t *testing.T) {
	manager := NewManager()
	if got := manager.Count(); got != 0 {
		t.Fatalf("Count() = %d", got)
	}
	if _, ok := manager.Get("missing"); ok {
		t.Fatal("Get(missing) ok = true")
	}
	if err := manager.CloseConn("missing"); err != ErrSessionNotFound {
		t.Fatalf("CloseConn() error = %v", err)
	}
}

// TestUserCode 验证用户 ID 和 WebSocket 会话编码的转换。
func TestUserCode(t *testing.T) {
	if got := UserCode(1001); got != "user:1001" {
		t.Fatalf("UserCode() = %q", got)
	}
	if got := UserCode(0); got != "" {
		t.Fatalf("UserCode(0) = %q", got)
	}
}

// TestUserID 验证从会话 code 中解析用户 ID，并兼容纯数字 code。
func TestUserID(t *testing.T) {
	if got := UserID(nil); got != 0 {
		t.Fatalf("UserID(nil) = %d", got)
	}

	if got := UserID(&Session{code: "user:1001"}); got != 1001 {
		t.Fatalf("UserID(user code) = %d", got)
	}
	if got := UserID(&Session{code: "1002"}); got != 1002 {
		t.Fatalf("UserID(raw code) = %d", got)
	}
	if got := UserID(&Session{code: "device:1003"}); got != 0 {
		t.Fatalf("UserID(non-user code) = %d", got)
	}
}

// TestManagerUserHelpers 验证用户 ID 便捷方法复用底层 code 查询和关闭逻辑。
func TestManagerUserHelpers(t *testing.T) {
	manager := NewManager()
	if _, ok := manager.GetUser(1001); ok {
		t.Fatal("GetUser(missing) ok = true")
	}
	if err := manager.CloseUserConn(1001); err != ErrSessionNotFound {
		t.Fatalf("CloseUserConn() error = %v", err)
	}
}

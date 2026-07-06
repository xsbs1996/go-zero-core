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

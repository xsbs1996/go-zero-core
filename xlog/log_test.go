package xlog

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// TestBodyNormalizesErrors 验证日志正文会把 error 规范化为字符串。
func TestBodyNormalizesErrors(t *testing.T) {
	text := Body("msg", map[string]any{"id": 1}, errors.New("failed"))

	var content Content
	if err := json.Unmarshal([]byte(text), &content); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if content.Msg != "msg" || content.Error != "failed" {
		t.Fatalf("content = %#v", content)
	}
}

// TestConfigWithDefaultAndToLogConf 验证日志配置默认值和 go-zero 配置转换。
func TestConfigWithDefaultAndToLogConf(t *testing.T) {
	conf := Config{}.WithDefault()
	if conf.Mode != ModeConsole || conf.Encoding != EncodingJSON || conf.Level != LevelInfo {
		t.Fatalf("WithDefault() = %#v", conf)
	}
	if conf.Stat == nil || !*conf.Stat {
		t.Fatalf("Stat default = %#v", conf.Stat)
	}

	logConf, err := conf.ToLogConf()
	if err != nil {
		t.Fatalf("ToLogConf() error = %v", err)
	}
	if logConf.Mode != ModeConsole || logConf.Encoding != EncodingJSON {
		t.Fatalf("LogConf = %#v", logConf)
	}
}

// TestContextHelpers 验证日志上下文辅助函数。
func TestContextHelpers(t *testing.T) {
	ctx := ContextWithDisable(context.Background())
	if !IsDisabled(ctx) {
		t.Fatal("IsDisabled() = false")
	}
	if safeContext(nil) == nil {
		t.Fatal("safeContext(nil) returned nil")
	}
}

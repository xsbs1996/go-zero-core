package xmysql

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

// TestConfigDefaultsValidateAndDSN 验证 MySQL 配置默认值、校验和 DSN 生成。
func TestConfigDefaultsValidateAndDSN(t *testing.T) {
	conf := Config{
		Host:      "127.0.0.1",
		Port:      3306,
		Username:  "root",
		Password:  "pass",
		Database:  "app",
		ParseTime: true,
	}
	if err := conf.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	withDefault := conf.WithDefault()
	if withDefault.Charset != defaultCharset || withDefault.Loc != defaultLoc {
		t.Fatalf("WithDefault() = %#v", withDefault)
	}
	if withDefault.Timeout != 10*time.Second || withDefault.ReadTimeout != 30*time.Second {
		t.Fatalf("timeouts = %v/%v", withDefault.Timeout, withDefault.ReadTimeout)
	}

	dsn := conf.DSN()
	for _, item := range []string{"root:pass@", "tcp(127.0.0.1:3306)", "/app", "charset=utf8mb4", "parseTime=true"} {
		if !strings.Contains(dsn, item) {
			t.Fatalf("DSN() = %q, missing %q", dsn, item)
		}
	}
}

// TestConfigValidateErrorsAndLogLevel 验证 MySQL 配置错误和日志级别转换。
func TestConfigValidateErrorsAndLogLevel(t *testing.T) {
	if err := (Config{}).Validate(); !errors.Is(err, ErrMissingHost) {
		t.Fatalf("Validate(empty) error = %v", err)
	}
	if err := (Config{Host: "127.0.0.1"}).Validate(); !errors.Is(err, ErrMissingUser) {
		t.Fatalf("Validate(missing user) error = %v", err)
	}
	if err := (Config{Host: "127.0.0.1", User: "root"}).Validate(); !errors.Is(err, ErrMissingDBName) {
		t.Fatalf("Validate(missing db) error = %v", err)
	}
	if got := (Config{LogLevel: "info"}).GormLogLevel(); got != logger.Info {
		t.Fatalf("GormLogLevel(info) = %v", got)
	}
}

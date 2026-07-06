package xpostgres

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

// TestConfigDefaultsValidateAndDSN 验证 PostgreSQL 配置默认值、校验和 DSN 生成。
func TestConfigDefaultsValidateAndDSN(t *testing.T) {
	conf := Config{
		Host:     "127.0.0.1",
		Port:     5432,
		Username: "postgres",
		Password: "pass",
		Database: "app",
		TimeZone: "Asia/Shanghai",
	}
	if err := conf.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	withDefault := conf.WithDefault()
	if withDefault.SSLMode != defaultSSLMode || withDefault.ConnectTimeout != 10*time.Second {
		t.Fatalf("WithDefault() = %#v", withDefault)
	}

	dsn := conf.DSN()
	for _, item := range []string{"postgres://postgres:pass@", "127.0.0.1:5432", "/app", "sslmode=disable", "connect_timeout=10", "TimeZone=Asia%2FShanghai"} {
		if !strings.Contains(dsn, item) {
			t.Fatalf("DSN() = %q, missing %q", dsn, item)
		}
	}
}

// TestConfigValidateErrorsAndLogLevel 验证 PostgreSQL 配置错误和日志级别转换。
func TestConfigValidateErrorsAndLogLevel(t *testing.T) {
	if err := (Config{}).Validate(); !errors.Is(err, ErrMissingHost) {
		t.Fatalf("Validate(empty) error = %v", err)
	}
	if err := (Config{Host: "127.0.0.1"}).Validate(); !errors.Is(err, ErrMissingUser) {
		t.Fatalf("Validate(missing user) error = %v", err)
	}
	if err := (Config{Host: "127.0.0.1", User: "postgres"}).Validate(); !errors.Is(err, ErrMissingDBName) {
		t.Fatalf("Validate(missing db) error = %v", err)
	}
	if got := (Config{LogLevel: "error"}).GormLogLevel(); got != logger.Error {
		t.Fatalf("GormLogLevel(error) = %v", got)
	}
}

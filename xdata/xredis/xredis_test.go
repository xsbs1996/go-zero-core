package xredis

import (
	"errors"
	"testing"
	"time"
)

// TestConfigDefaultsValidateAndOptions 验证 Redis 配置默认值、校验和原生 options 转换。
func TestConfigDefaultsValidateAndOptions(t *testing.T) {
	conf := Config{Addr: "127.0.0.1:6379", Password: "pass", DB: 2}
	if err := conf.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	withDefault := conf.WithDefault()
	if withDefault.PoolSize != defaultPoolSize || withDefault.DialTimeout != 5*time.Second {
		t.Fatalf("WithDefault() = %#v", withDefault)
	}

	options := redisOptions(conf)
	if options.Addr != conf.Addr || options.Password != conf.Password || options.DB != conf.DB {
		t.Fatalf("redisOptions() = %#v", options)
	}
	if options.PoolSize != defaultPoolSize {
		t.Fatalf("PoolSize = %d", options.PoolSize)
	}
}

// TestConfigValidateError 验证 Redis 地址为空时返回配置错误。
func TestConfigValidateError(t *testing.T) {
	if err := (Config{}).Validate(); !errors.Is(err, ErrMissingAddr) {
		t.Fatalf("Validate() error = %v", err)
	}
}

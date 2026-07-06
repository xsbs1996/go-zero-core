package xrabbitmq

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestConfigDefaultsValidateAndDSN 验证 RabbitMQ 配置默认值、校验、DSN 和 AMQP 配置转换。
func TestConfigDefaultsValidateAndDSN(t *testing.T) {
	conf := Config{
		Host:     "127.0.0.1",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		VHost:    "app",
		TLS:      true,
	}
	if err := conf.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	withDefault := conf.WithDefault()
	if withDefault.Locale != defaultLocale || withDefault.Heartbeat != 10*time.Second {
		t.Fatalf("WithDefault() = %#v", withDefault)
	}

	dsn := conf.DSN()
	for _, item := range []string{"amqps://guest:guest@", "127.0.0.1:5672", "/app"} {
		if !strings.Contains(dsn, item) {
			t.Fatalf("DSN() = %q, missing %q", dsn, item)
		}
	}

	amqpConf := conf.AMQPConfig()
	if amqpConf.Heartbeat != defaultHeartbeat || amqpConf.Locale != defaultLocale || amqpConf.Dial == nil {
		t.Fatalf("AMQPConfig() = %#v", amqpConf)
	}
}

// TestConfigValidateErrorsAndManagerEmptyState 验证 RabbitMQ 配置错误和空 manager 状态。
func TestConfigValidateErrorsAndManagerEmptyState(t *testing.T) {
	if err := (Config{}).Validate(); !errors.Is(err, ErrMissingHost) {
		t.Fatalf("Validate(empty) error = %v", err)
	}
	if err := (Config{Host: "127.0.0.1"}).Validate(); !errors.Is(err, ErrMissingUsername) {
		t.Fatalf("Validate(missing user) error = %v", err)
	}
	if err := (Config{URL: "amqp://guest:guest@localhost:5672/"}).Validate(); err != nil {
		t.Fatalf("Validate(URL) error = %v", err)
	}

	manager := NewManager()
	if manager.IsProducerRegistered("missing") || manager.IsConsumerRegistered("missing") {
		t.Fatal("empty manager should not have registered endpoints")
	}
	if err := manager.PublishBatch(nil, "missing"); !errors.Is(err, ErrInvalidPublishBatch) {
		t.Fatalf("PublishBatch(empty) error = %v", err)
	}
}

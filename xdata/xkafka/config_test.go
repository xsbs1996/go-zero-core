package xkafka

import (
	"errors"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

// TestConfigDefaultsAndValidation 验证 Kafka 配置默认值和生产/消费配置校验。
func TestConfigDefaultsAndValidation(t *testing.T) {
	conf := Config{Brokers: []string{"127.0.0.1:9092"}}
	if err := conf.ValidateProducer(); err != nil {
		t.Fatalf("ValidateProducer() error = %v", err)
	}
	if err := conf.ValidateConsumer(); err != nil {
		t.Fatalf("ValidateConsumer() error = %v", err)
	}

	withDefault := conf.WithDefault()
	if withDefault.BatchSize != 100 || withDefault.ProducerBatchTimeout() != time.Second {
		t.Fatalf("WithDefault() = %#v", withDefault)
	}
}

// TestProducerConsumerConstructionWithoutPing 验证跳过 ping 时可本地构造生产者和消费者。
func TestProducerConsumerConstructionWithoutPing(t *testing.T) {
	conf := Config{Brokers: []string{"127.0.0.1:9092"}}
	writer, err := NewProducer(" topic ", conf, WithoutProducerPing())
	if err != nil {
		t.Fatalf("NewProducer() error = %v", err)
	}
	if writer.Topic != "topic" {
		t.Fatalf("writer topic = %q", writer.Topic)
	}
	_ = writer.Close()

	reader, err := NewConsumer(" topic ", " group ", conf, WithoutConsumerPing())
	if err != nil {
		t.Fatalf("NewConsumer() error = %v", err)
	}
	_ = reader.Close()
}

// TestKafkaErrorsAndManagerEmptyState 验证 Kafka 错误分支和空 manager 状态。
func TestKafkaErrorsAndManagerEmptyState(t *testing.T) {
	if err := (Config{}).ValidateProducer(); !errors.Is(err, ErrMissingBrokers) {
		t.Fatalf("ValidateProducer(empty) error = %v", err)
	}
	if _, err := NewProducer("", Config{Brokers: []string{"b"}}, WithWriter(&kafka.Writer{})); !errors.Is(err, ErrMissingTopic) {
		t.Fatalf("NewProducer(empty topic) error = %v", err)
	}

	manager := NewManager()
	if _, err := manager.GetProducer("missing"); !errors.Is(err, ErrProducerNotRegistered) {
		t.Fatalf("GetProducer() error = %v", err)
	}
	if manager.IsProducerRegistered("missing") {
		t.Fatal("IsProducerRegistered(missing) = true")
	}
}

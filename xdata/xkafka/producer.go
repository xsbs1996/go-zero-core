package xkafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNilWriter         = errors.New("xkafka: writer is nil")       // ErrNilWriter 表示 Kafka writer 为空。
	ErrInvalidWriteBatch = errors.New("xkafka: empty message batch") // ErrInvalidWriteBatch 表示批量消息为空。
)

// writer 构建 Kafka 生产者配置。
func writer(conf Config, topic string) *kafka.Writer {
	conf = conf.WithDefault()
	return &kafka.Writer{
		Addr:         kafka.TCP(conf.Brokers...),
		Topic:        topic,
		Balancer:     balancer(),
		BatchSize:    conf.BatchSize,
		BatchTimeout: conf.ProducerBatchTimeout(),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		RequiredAcks: conf.RequiredAcks,
		Async:        conf.Async,
	}
}

// NewProducer 根据 topic 和配置创建 Kafka 生产者。
func NewProducer(topic string, conf Config, opts ...ProducerOption) (*kafka.Writer, error) {
	if err := conf.ValidateProducer(); err != nil {
		return nil, err
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return nil, err
	}

	conf = conf.WithDefault()
	options := newProducerOptions(conf, topic, opts...)
	return options.writer, nil
}

// MustNewProducer 根据 topic 和配置创建 Kafka 生产者，失败时直接 panic。
func MustNewProducer(topic string, conf Config, opts ...ProducerOption) *kafka.Writer {
	writer, err := NewProducer(topic, conf, opts...)
	if err != nil {
		panic(err)
	}
	return writer
}

// closeProducer 关闭 Kafka 生产者。
func closeProducer(writer *kafka.Writer) error {
	if writer == nil {
		return ErrNilProducer
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close kafka producer failed: %w", err)
	}
	return nil
}

// writeMessage 使用 Kafka writer 写入单条消息。
func writeMessage(ctx context.Context, writer *kafka.Writer, msg kafka.Message) error {
	if writer == nil {
		return ErrNilWriter
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("produce kafka message failed: %w", err)
	}
	return nil
}

// writeMessages 使用 Kafka writer 批量写入消息。
func writeMessages(ctx context.Context, writer *kafka.Writer, msgs ...kafka.Message) error {
	if writer == nil {
		return ErrNilWriter
	}
	if len(msgs) == 0 {
		return ErrInvalidWriteBatch
	}

	if err := writer.WriteMessages(ctx, msgs...); err != nil {
		return fmt.Errorf("produce kafka messages failed: %w", err)
	}
	return nil
}

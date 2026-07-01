package kafkax

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Client 定义 Kafka 工具箱对业务暴露的核心能力。
type Client interface {
	// RegisterProducer 注册指定 topic 的生产者，默认会先检查 broker 连通性，生产者按 topic 缓存复用。
	RegisterProducer(topic string, conf Config, opts ...ProducerOption) error
	// RegisterConsumer 注册指定 topic 和 group 的消费者，默认会先检查 broker 连通性，消费者按 topic+group 缓存复用。
	RegisterConsumer(topic, group string, conf Config, opts ...ConsumerOption) error
	// Produce 向指定 topic 发送单条消息。
	Produce(ctx context.Context, topic string, msg kafka.Message) error
	// ProduceBatch 向指定 topic 批量发送消息。
	ProduceBatch(ctx context.Context, topic string, msgs ...kafka.Message) error
	// Consume 注册指定 topic 和 group 的单条消息消费处理函数。
	Consume(ctx context.Context, topic, group string, handler func(context.Context, kafka.Message) error) error
	// ConsumeBatch 注册指定 topic 和 group 的批量消息消费处理函数。
	ConsumeBatch(ctx context.Context, topic, group string, batchSize int, batchTimeout time.Duration, handler func(context.Context, []kafka.Message) error) error
	// CloseProducer 关闭指定 topic 的生产者。
	CloseProducer(topic string) error
	// CloseConsumer 关闭指定 topic 和 group 的消费者。
	CloseConsumer(topic, group string) error
	// Close 关闭全部生产者和消费者。
	Close() error
	// IsProducerRegistered 返回指定 topic 的生产者是否已注册。
	IsProducerRegistered(topic string) bool
	// IsConsumerRegistered 返回指定 topic 和 group 的消费者是否已注册。
	IsConsumerRegistered(topic, group string) bool
}

// NewClient 创建一个 Kafka 客户端接口实例。
func NewClient() Client {
	return NewManager()
}

// DefaultClient 返回包级默认 Kafka 客户端接口。
func DefaultClient() Client {
	return defaultManager
}

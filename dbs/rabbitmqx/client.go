package rabbitmqx

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client 定义 RabbitMQ 工具箱对业务暴露的核心能力。
type Client interface {
	// RegisterProducer 注册生产者。
	RegisterProducer(name string, conf ProducerConfig, opts ...ProducerOption) error
	// RegisterConsumer 注册消费者。
	RegisterConsumer(name string, conf ConsumerConfig, opts ...ConsumerOption) error
	// Publish 通过指定生产者发布单条消息。
	Publish(ctx context.Context, name string, msg amqp.Publishing) error
	// PublishBatch 通过指定生产者批量发布消息。
	PublishBatch(ctx context.Context, name string, msgs ...amqp.Publishing) error
	// Consume 通过指定消费者注册消费处理函数。
	Consume(ctx context.Context, name string, handler Handler) error
	// CloseProducer 关闭指定生产者。
	CloseProducer(name string) error
	// CloseConsumer 关闭指定消费者。
	CloseConsumer(name string) error
	// Close 关闭全部生产者和消费者。
	Close() error
	// IsProducerRegistered 返回生产者是否已注册。
	IsProducerRegistered(name string) bool
	// IsConsumerRegistered 返回消费者是否已注册。
	IsConsumerRegistered(name string) bool
}

var defaultClient = NewManager() // defaultClient 默认 RabbitMQ 客户端。

// NewClient 创建一个 RabbitMQ 客户端接口实例。
func NewClient() Client {
	return NewManager()
}

// DefaultClient 返回包级默认 RabbitMQ 客户端接口。
func DefaultClient() Client {
	return defaultClient
}

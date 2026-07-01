package rabbitmqx

import amqp "github.com/rabbitmq/amqp091-go"

type producerOptions struct {
	channel *amqp.Channel // channel RabbitMQ 原生 channel。
	ping    bool          // ping 是否在注册时检查连接。
}

type consumerOptions struct {
	channel *amqp.Channel // channel RabbitMQ 原生 channel。
	ping    bool          // ping 是否在注册时检查连接。
}

// ProducerOption RabbitMQ 生产者注册的可选配置。
type ProducerOption func(*producerOptions)

// ConsumerOption RabbitMQ 消费者注册的可选配置。
type ConsumerOption func(*consumerOptions)

// WithProducerChannel 使用 RabbitMQ 原生 channel 作为生产者 channel。
func WithProducerChannel(channel *amqp.Channel) ProducerOption {
	return func(o *producerOptions) {
		if channel != nil {
			o.channel = channel
		}
	}
}

// WithConsumerChannel 使用 RabbitMQ 原生 channel 作为消费者 channel。
func WithConsumerChannel(channel *amqp.Channel) ConsumerOption {
	return func(o *consumerOptions) {
		if channel != nil {
			o.channel = channel
		}
	}
}

// WithoutProducerPing 跳过生产者注册时的连接检查。
func WithoutProducerPing() ProducerOption {
	return func(o *producerOptions) {
		o.ping = false
	}
}

// WithoutConsumerPing 跳过消费者注册时的连接检查。
func WithoutConsumerPing() ConsumerOption {
	return func(o *consumerOptions) {
		o.ping = false
	}
}

func newProducerOptions(opts ...ProducerOption) producerOptions {
	options := producerOptions{ping: true}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

func newConsumerOptions(opts ...ConsumerOption) consumerOptions {
	options := consumerOptions{ping: true}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

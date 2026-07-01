package kafkax

import "github.com/segmentio/kafka-go"

type producerOptions struct {
	writer *kafka.Writer // writer Kafka 原生生产者。
	ping   bool          // ping 是否在注册时检查 broker 连通性。
}

type consumerOptions struct {
	readerConfig *kafka.ReaderConfig // readerConfig Kafka 原生消费者配置。
	ping         bool                // ping 是否在注册时检查 broker 连通性。
}

// ProducerOption Kafka 生产者初始化的可选配置。
type ProducerOption func(*producerOptions)

// ConsumerOption Kafka 消费者初始化的可选配置。
type ConsumerOption func(*consumerOptions)

// WithWriter 使用 kafka-go 原生 Writer。
func WithWriter(writer *kafka.Writer) ProducerOption {
	return func(o *producerOptions) {
		if writer != nil {
			o.writer = writer
		}
	}
}

// WithReaderConfig 使用 kafka-go 原生 ReaderConfig。
func WithReaderConfig(conf *kafka.ReaderConfig) ConsumerOption {
	return func(o *consumerOptions) {
		if conf != nil {
			o.readerConfig = conf
		}
	}
}

// WithoutProducerPing 跳过生产者注册时的 broker 连通性检查。
func WithoutProducerPing() ProducerOption {
	return func(o *producerOptions) {
		o.ping = false
	}
}

// WithoutConsumerPing 跳过消费者注册时的 broker 连通性检查。
func WithoutConsumerPing() ConsumerOption {
	return func(o *consumerOptions) {
		o.ping = false
	}
}

// newProducerOptions 构建生产者可选配置。
func newProducerOptions(conf Config, topic string, opts ...ProducerOption) producerOptions {
	options := producerOptions{
		writer: writer(conf, topic),
		ping:   true,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

// newConsumerOptions 构建消费者可选配置。
func newConsumerOptions(conf Config, topic, group string, opts ...ConsumerOption) consumerOptions {
	options := consumerOptions{
		readerConfig: readerConfig(conf, topic, group),
		ping:         true,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

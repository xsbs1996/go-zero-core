package kafkax

import "errors"

var (
	ErrProducerAlreadyInitialized = ErrProducerAlreadyRegistered          // ErrProducerAlreadyInitialized 表示默认 Kafka 生产者已经初始化。
	ErrConsumerAlreadyInitialized = ErrConsumerAlreadyRegistered          // ErrConsumerAlreadyInitialized 表示默认 Kafka 消费者已经初始化。
	ErrProducerNotInitialized     = ErrProducerNotRegistered              // ErrProducerNotInitialized 表示默认 Kafka 生产者未初始化。
	ErrConsumerNotInitialized     = ErrConsumerNotRegistered              // ErrConsumerNotInitialized 表示默认 Kafka 消费者未初始化。
	ErrNilProducer                = errors.New("kafkax: producer is nil") // ErrNilProducer 表示 Kafka 生产者为空。
	ErrNilConsumer                = errors.New("kafkax: consumer is nil") // ErrNilConsumer 表示 Kafka 消费者为空。
)

var defaultManager = NewManager() // defaultManager 默认 Kafka 客户端。

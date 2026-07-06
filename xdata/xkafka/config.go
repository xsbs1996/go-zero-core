package xkafka

import (
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	defaultDialTimeout  = 10 * time.Second // 默认建立连接超时时间。
	defaultReadTimeout  = 10 * time.Second // 默认读取超时时间。
	defaultWriteTimeout = 10 * time.Second // 默认写入超时时间。
	defaultBatchTimeout = 1000             // 默认批量等待时间，单位毫秒。
	defaultBatchSize    = 100              // 默认批量大小。
)

var (
	ErrMissingBrokers = errors.New("xkafka: missing kafka brokers")  // ErrMissingBrokers 表示 Kafka broker 地址为空。
	ErrMissingTopic   = errors.New("xkafka: missing kafka topic")    // ErrMissingTopic 表示 Kafka topic 为空。
	ErrMissingGroupID = errors.New("xkafka: missing kafka group id") // ErrMissingGroupID 表示 Kafka 消费组为空。
)

// Config Kafka 连接配置。
type Config struct {
	Brokers             []string           `json:"brokers" yaml:"brokers"`                                      // Brokers Kafka broker 地址列表。
	ClientID            string             `json:"clientId,optional" yaml:"clientId"`                           // ClientID Kafka 客户端 ID。
	DialTimeout         time.Duration      `json:"dialTimeout,default=10s" yaml:"dialTimeout"`                  // DialTimeout 建立连接超时时间，配置值示例：10s。
	ReadTimeout         time.Duration      `json:"readTimeout,default=10s" yaml:"readTimeout"`                  // ReadTimeout 读取超时时间，配置值示例：10s。
	WriteTimeout        time.Duration      `json:"writeTimeout,default=10s" yaml:"writeTimeout"`                // WriteTimeout 写入超时时间，配置值示例：10s。
	BatchSize           int                `json:"batchSize,default=100" yaml:"batchSize"`                      // BatchSize 生产者批量发送大小。
	BatchTimeout        int                `json:"batchTimeout,default=1000" yaml:"batchTimeout"`               // BatchTimeout 生产者批量发送等待时间，单位毫秒。
	ConsumeBatchSize    int                `json:"consumeBatchSize,default=100" yaml:"consumeBatchSize"`        // ConsumeBatchSize 批量消费上限，达到后立即触发 handler。
	ConsumeBatchTimeout int                `json:"consumeBatchTimeout,default=1000" yaml:"consumeBatchTimeout"` // ConsumeBatchTimeout 批量消费最长等待时间，单位毫秒。
	RequiredAcks        kafka.RequiredAcks `json:"requiredAcks,optional" yaml:"requiredAcks"`                   // RequiredAcks Kafka 写入确认级别，0 表示不等待确认，-1 表示等待所有副本确认。
	Async               bool               `json:"async,optional" yaml:"async"`                                 // Async 是否异步写入消息。
}

// WithDefault 返回补齐默认值后的配置。
func (c Config) WithDefault() Config {
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultDialTimeout
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.BatchSize == 0 {
		c.BatchSize = defaultBatchSize
	}
	if c.BatchTimeout == 0 {
		c.BatchTimeout = defaultBatchTimeout
	}
	if c.ConsumeBatchSize == 0 {
		c.ConsumeBatchSize = defaultBatchSize
	}
	if c.ConsumeBatchTimeout == 0 {
		c.ConsumeBatchTimeout = defaultBatchTimeout
	}
	return c
}

// ProducerBatchTimeout 返回生产者批量发送等待时间。
func (c Config) ProducerBatchTimeout() time.Duration {
	c = c.WithDefault()
	return time.Duration(c.BatchTimeout) * time.Millisecond
}

// ConsumerBatchTimeout 返回批量消费最长等待时间。
func (c Config) ConsumerBatchTimeout() time.Duration {
	c = c.WithDefault()
	return time.Duration(c.ConsumeBatchTimeout) * time.Millisecond
}

// ValidateProducer 校验 Kafka 生产者配置。
func (c Config) ValidateProducer() error {
	if len(c.Brokers) == 0 {
		return ErrMissingBrokers
	}
	return nil
}

// ValidateConsumer 校验 Kafka 消费者配置。
func (c Config) ValidateConsumer() error {
	if len(c.Brokers) == 0 {
		return ErrMissingBrokers
	}
	return nil
}

// balancer 返回默认 Kafka 分区均衡策略。
func balancer() kafka.Balancer {
	return &kafka.LeastBytes{}
}

package xkafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrProducerAlreadyRegistered = errors.New("xkafka: producer already registered") // ErrProducerAlreadyRegistered 表示生产者路由已经注册。
	ErrConsumerAlreadyRegistered = errors.New("xkafka: consumer already registered") // ErrConsumerAlreadyRegistered 表示消费者路由已经注册。
	ErrProducerNotRegistered     = errors.New("xkafka: producer not registered")     // ErrProducerNotRegistered 表示生产者路由未注册。
	ErrConsumerNotRegistered     = errors.New("xkafka: consumer not registered")     // ErrConsumerNotRegistered 表示消费者路由未注册。
	ErrNilManager                = errors.New("xkafka: manager is nil")              // ErrNilManager 表示管理器为空。
)

var _ Client = (*Manager)(nil) // 确保 Manager 实现 Client 接口。

// Manager 管理多个 Kafka 生产者与消费者路由。
type Manager struct {
	mu              sync.RWMutex             // mu 保护生产者和消费者缓存。
	producers       map[string]*kafka.Writer // producers 按 topic 缓存生产者。
	consumers       map[string]*kafka.Reader // consumers 按 topic+group 缓存消费者。
	consumerConfigs map[string]Config        // consumerConfigs 按 topic+group 缓存消费者配置。
}

// NewManager 创建一个 Kafka 路由管理器。
func NewManager() *Manager {
	return &Manager{
		producers:       make(map[string]*kafka.Writer),
		consumers:       make(map[string]*kafka.Reader),
		consumerConfigs: make(map[string]Config),
	}
}

// RegisterProducer 注册一个 Kafka 生产者路由。
func (m *Manager) RegisterProducer(topic string, conf Config, opts ...ProducerOption) error {
	if m == nil {
		return ErrNilManager
	}
	if err := conf.ValidateProducer(); err != nil {
		return err
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return err
	}

	conf = conf.WithDefault()
	options := newProducerOptions(conf, topic, opts...)
	if options.ping {
		if err := Ping(context.Background(), conf); err != nil {
			return err
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.producers[topic]; ok {
		return ErrProducerAlreadyRegistered
	}

	m.producers[topic] = options.writer
	return nil
}

// RegisterConsumer 注册一个 Kafka 消费者路由。
func (m *Manager) RegisterConsumer(topic, group string, conf Config, opts ...ConsumerOption) error {
	if m == nil {
		return ErrNilManager
	}
	if err := conf.ValidateConsumer(); err != nil {
		return err
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return err
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return err
	}

	key := consumerKey(topic, group)

	conf = conf.WithDefault()
	options := newConsumerOptions(conf, topic, group, opts...)
	if options.ping {
		if err := Ping(context.Background(), conf); err != nil {
			return err
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.consumers[key]; ok {
		return ErrConsumerAlreadyRegistered
	}

	m.consumers[key] = kafka.NewReader(*options.readerConfig)
	m.consumerConfigs[key] = conf
	return nil
}

// GetProducer 返回指定 topic 的 Kafka 生产者。
func (m *Manager) GetProducer(topic string) (*kafka.Writer, error) {
	if m == nil {
		return nil, ErrNilManager
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	writer, ok := m.producers[topic]
	if !ok || writer == nil {
		return nil, ErrProducerNotRegistered
	}
	return writer, nil
}

// GetConsumer 返回指定 topic 和 group 的 Kafka 消费者。
func (m *Manager) GetConsumer(topic, group string) (*kafka.Reader, error) {
	if m == nil {
		return nil, ErrNilManager
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return nil, err
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return nil, err
	}

	key := consumerKey(topic, group)

	m.mu.RLock()
	defer m.mu.RUnlock()

	reader, ok := m.consumers[key]
	if !ok || reader == nil {
		return nil, ErrConsumerNotRegistered
	}
	return reader, nil
}

// CloseProducer 关闭指定 topic 的 Kafka 生产者。
func (m *Manager) CloseProducer(topic string) error {
	if m == nil {
		return ErrNilManager
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	writer, ok := m.producers[topic]
	if !ok || writer == nil {
		return ErrProducerNotRegistered
	}
	if err := closeProducer(writer); err != nil {
		return err
	}

	delete(m.producers, topic)
	return nil
}

// CloseConsumer 关闭指定 topic 和 group 的 Kafka 消费者。
func (m *Manager) CloseConsumer(topic, group string) error {
	if m == nil {
		return ErrNilManager
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return err
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return err
	}

	key := consumerKey(topic, group)

	m.mu.Lock()
	defer m.mu.Unlock()

	reader, ok := m.consumers[key]
	if !ok || reader == nil {
		return ErrConsumerNotRegistered
	}
	if err := closeConsumer(reader); err != nil {
		return err
	}

	delete(m.consumers, key)
	delete(m.consumerConfigs, key)
	return nil
}

// Close 关闭管理器中的所有 Kafka 资源。
func (m *Manager) Close() error {
	if m == nil {
		return ErrNilManager
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for topic, writer := range m.producers {
		if writer == nil {
			continue
		}
		if err := writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka producer %q failed: %w", topic, err))
		}
		delete(m.producers, topic)
	}

	for key, reader := range m.consumers {
		if reader == nil {
			continue
		}
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close kafka consumer %q failed: %w", key, err))
		}
		delete(m.consumers, key)
		delete(m.consumerConfigs, key)
	}

	return errors.Join(errs...)
}

// Produce 通过指定 topic 发送单条消息。
func (m *Manager) Produce(ctx context.Context, topic string, msg kafka.Message) error {
	writer, err := m.GetProducer(topic)
	if err != nil {
		return err
	}
	return writeMessage(ctx, writer, msg)
}

// ProduceBatch 通过指定 topic 批量发送消息。
func (m *Manager) ProduceBatch(ctx context.Context, topic string, msgs ...kafka.Message) error {
	writer, err := m.GetProducer(topic)
	if err != nil {
		return err
	}
	return writeMessages(ctx, writer, msgs...)
}

// Consume 通过指定 topic 和 group 注册单条消息消费处理函数。
func (m *Manager) Consume(ctx context.Context, topic, group string, handler func(context.Context, kafka.Message) error) error {
	reader, err := m.GetConsumer(topic, group)
	if err != nil {
		return err
	}
	return readMessageLoop(ctx, reader, handler)
}

// ConsumeBatch 通过指定 topic 和 group 注册批量消息消费处理函数。
//
// batchSize 表示批量消费上限，传 0 时使用 Config.ConsumeBatchSize。
// batchTimeout 表示批量消费最长等待时间，传 0 时使用 Config.ConsumeBatchTimeout。
func (m *Manager) ConsumeBatch(ctx context.Context, topic, group string, batchSize int, batchTimeout time.Duration, handler func(context.Context, []kafka.Message) error) error {
	reader, err := m.GetConsumer(topic, group)
	if err != nil {
		return err
	}
	if batchSize <= 0 || batchTimeout <= 0 {
		conf, err := m.GetConsumerConfig(topic, group)
		if err != nil {
			return err
		}
		if batchSize <= 0 {
			batchSize = conf.ConsumeBatchSize
		}
		if batchTimeout <= 0 {
			batchTimeout = conf.ConsumerBatchTimeout()
		}
	}
	return readBatchLoop(ctx, reader, batchSize, batchTimeout, handler)
}

// GetConsumerConfig 返回指定 topic 和 group 的 Kafka 消费者配置。
func (m *Manager) GetConsumerConfig(topic, group string) (Config, error) {
	if m == nil {
		return Config{}, ErrNilManager
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return Config{}, err
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return Config{}, err
	}

	key := consumerKey(topic, group)

	m.mu.RLock()
	defer m.mu.RUnlock()

	conf, ok := m.consumerConfigs[key]
	if !ok {
		return Config{}, ErrConsumerNotRegistered
	}
	return conf, nil
}

// IsProducerRegistered 返回指定 topic 的 Kafka 生产者是否已注册。
func (m *Manager) IsProducerRegistered(topic string) bool {
	if m == nil {
		return false
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.producers[topic]
	return ok
}

// IsConsumerRegistered 返回指定 topic 和 group 的 Kafka 消费者是否已注册。
func (m *Manager) IsConsumerRegistered(topic, group string) bool {
	if m == nil {
		return false
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return false
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return false
	}

	key := consumerKey(topic, group)

	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.consumers[key]
	return ok
}

// normalizeTopic 标准化 topic 参数。
func normalizeTopic(topic string) (string, error) {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return "", ErrMissingTopic
	}
	return topic, nil
}

// normalizeGroup 标准化 group 参数。
func normalizeGroup(group string) (string, error) {
	group = strings.TrimSpace(group)
	if group == "" {
		return "", ErrMissingGroupID
	}
	return group, nil
}

// consumerKey 生成消费者缓存键。
func consumerKey(topic, group string) string {
	return topic + "\x00" + group
}

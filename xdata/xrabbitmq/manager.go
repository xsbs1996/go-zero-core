package xrabbitmq

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrMissingName               = errors.New("xrabbitmq: missing name")                 // ErrMissingName 表示注册名称为空。
	ErrNilConnection             = errors.New("xrabbitmq: connection is nil")            // ErrNilConnection 表示 RabbitMQ 连接为空。
	ErrNilChannel                = errors.New("xrabbitmq: channel is nil")               // ErrNilChannel 表示 RabbitMQ channel 为空。
	ErrNilHandler                = errors.New("xrabbitmq: handler is nil")               // ErrNilHandler 表示消费处理函数为空。
	ErrInvalidPublishBatch       = errors.New("xrabbitmq: empty publish batch")          // ErrInvalidPublishBatch 表示批量发布消息为空。
	ErrProducerAlreadyRegistered = errors.New("xrabbitmq: producer already registered")  // ErrProducerAlreadyRegistered 表示生产者已经注册。
	ErrConsumerAlreadyRegistered = errors.New("xrabbitmq: consumer already registered")  // ErrConsumerAlreadyRegistered 表示消费者已经注册。
	ErrProducerNotRegistered     = errors.New("xrabbitmq: producer not registered")      // ErrProducerNotRegistered 表示生产者未注册。
	ErrConsumerNotRegistered     = errors.New("xrabbitmq: consumer not registered")      // ErrConsumerNotRegistered 表示消费者未注册。
	ErrMissingQueue              = errors.New("xrabbitmq: missing rabbitmq queue")       // ErrMissingQueue 表示 RabbitMQ queue 为空。
	ErrMissingExchange           = errors.New("xrabbitmq: missing rabbitmq exchange")    // ErrMissingExchange 表示 RabbitMQ exchange 为空。
	ErrMissingRoutingKey         = errors.New("xrabbitmq: missing rabbitmq routing key") // ErrMissingRoutingKey 表示 RabbitMQ routing key 为空。
	ErrNilClient                 = errors.New("xrabbitmq: client is nil")                // ErrNilClient 表示 RabbitMQ 客户端为空。
)

var _ Client = (*Manager)(nil) // 确保 Manager 实现 Client 接口。

type producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	conf    ProducerConfig
}

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	conf    ConsumerConfig
}

// Manager 管理多个 RabbitMQ 生产者和消费者。
type Manager struct {
	mu        sync.RWMutex
	producers map[string]*producer
	consumers map[string]*consumer
}

// NewManager 创建一个 RabbitMQ 管理器。
func NewManager() *Manager {
	return &Manager{
		producers: make(map[string]*producer),
		consumers: make(map[string]*consumer),
	}
}

// RegisterProducer 注册生产者。
func (m *Manager) RegisterProducer(name string, conf ProducerConfig, opts ...ProducerOption) error {
	if m == nil {
		return ErrNilClient
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}
	if err := validateProducerConfig(conf); err != nil {
		return err
	}

	options := newProducerOptions(opts...)
	if options.ping {
		if err := Ping(conf.Connection); err != nil {
			return err
		}
	}

	channel := options.channel
	var conn *amqp.Connection
	if channel == nil {
		conn, err = Connect(conf.Connection)
		if err != nil {
			return err
		}
		channel, err = OpenChannel(conn)
		if err != nil {
			_ = conn.Close()
			return err
		}
	}

	if err := declareProducer(channel, conf); err != nil {
		if conn != nil {
			_ = conn.Close()
		}
		if channel != nil {
			_ = channel.Close()
		}
		return fmt.Errorf("declare rabbitmq producer failed: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.producers[name]; ok {
		if conn != nil {
			_ = conn.Close()
		}
		if channel != nil {
			_ = channel.Close()
		}
		return ErrProducerAlreadyRegistered
	}

	m.producers[name] = &producer{conn: conn, channel: channel, conf: conf}
	return nil
}

// RegisterConsumer 注册消费者。
func (m *Manager) RegisterConsumer(name string, conf ConsumerConfig, opts ...ConsumerOption) error {
	if m == nil {
		return ErrNilClient
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}
	if err := validateConsumerConfig(conf); err != nil {
		return err
	}

	options := newConsumerOptions(opts...)
	if options.ping {
		if err := Ping(conf.Connection); err != nil {
			return err
		}
	}

	channel := options.channel
	var conn *amqp.Connection
	if channel == nil {
		conn, err = Connect(conf.Connection)
		if err != nil {
			return err
		}
		channel, err = OpenChannel(conn)
		if err != nil {
			_ = conn.Close()
			return err
		}
	}

	if err := declareConsumer(channel, conf); err != nil {
		if conn != nil {
			_ = conn.Close()
		}
		if channel != nil {
			_ = channel.Close()
		}
		return fmt.Errorf("declare rabbitmq consumer failed: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.consumers[name]; ok {
		if conn != nil {
			_ = conn.Close()
		}
		if channel != nil {
			_ = channel.Close()
		}
		return ErrConsumerAlreadyRegistered
	}

	m.consumers[name] = &consumer{conn: conn, channel: channel, conf: conf}
	return nil
}

// CloseProducer 关闭指定生产者。
func (m *Manager) CloseProducer(name string) error {
	if m == nil {
		return ErrNilClient
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.producers[name]
	if !ok || p == nil {
		return ErrProducerNotRegistered
	}
	if err := closeProducer(p); err != nil {
		return err
	}
	delete(m.producers, name)
	return nil
}

// CloseConsumer 关闭指定消费者。
func (m *Manager) CloseConsumer(name string) error {
	if m == nil {
		return ErrNilClient
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.consumers[name]
	if !ok || c == nil {
		return ErrConsumerNotRegistered
	}
	if err := closeConsumer(c); err != nil {
		return err
	}
	delete(m.consumers, name)
	return nil
}

// Close 关闭全部生产者和消费者。
func (m *Manager) Close() error {
	if m == nil {
		return ErrNilClient
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, p := range m.producers {
		if err := closeProducer(p); err != nil {
			errs = append(errs, fmt.Errorf("close rabbitmq producer %q failed: %w", name, err))
		}
		delete(m.producers, name)
	}
	for name, c := range m.consumers {
		if err := closeConsumer(c); err != nil {
			errs = append(errs, fmt.Errorf("close rabbitmq consumer %q failed: %w", name, err))
		}
		delete(m.consumers, name)
	}
	return errors.Join(errs...)
}

// IsProducerRegistered 返回生产者是否已注册。
func (m *Manager) IsProducerRegistered(name string) bool {
	if m == nil {
		return false
	}
	name, err := normalizeName(name)
	if err != nil {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.producers[name]
	return ok
}

// IsConsumerRegistered 返回消费者是否已注册。
func (m *Manager) IsConsumerRegistered(name string) bool {
	if m == nil {
		return false
	}
	name, err := normalizeName(name)
	if err != nil {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.consumers[name]
	return ok
}

func normalizeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrMissingName
	}
	return name, nil
}

func validateProducerConfig(conf ProducerConfig) error {
	if err := conf.Connection.Validate(); err != nil {
		return err
	}
	if conf.Exchange.Kind == "fanout" {
		return nil
	}
	if conf.RoutingKey == "" && conf.Binding.RoutingKey == "" && conf.Queue.Name == "" {
		return ErrMissingRoutingKey
	}
	return nil
}

func validateConsumerConfig(conf ConsumerConfig) error {
	if err := conf.Connection.Validate(); err != nil {
		return err
	}
	if conf.Queue.Name == "" && conf.Binding.Queue == "" {
		return ErrMissingQueue
	}
	return nil
}

func closeProducer(p *producer) error {
	if p == nil {
		return ErrProducerNotRegistered
	}
	var errs []error
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func closeConsumer(c *consumer) error {
	if c == nil {
		return ErrConsumerNotRegistered
	}
	var errs []error
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

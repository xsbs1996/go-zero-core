package xrabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish 通过指定生产者发布单条消息。
func (m *Manager) Publish(ctx context.Context, name string, msg amqp.Publishing) error {
	if m == nil {
		return ErrNilClient
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}

	m.mu.RLock()
	p, ok := m.producers[name]
	m.mu.RUnlock()
	if !ok || p == nil {
		return ErrProducerNotRegistered
	}
	if p.channel == nil {
		return ErrNilChannel
	}

	routingKey := p.conf.RoutingKey
	if routingKey == "" {
		routingKey = p.conf.Binding.RoutingKey
	}

	if err := p.channel.PublishWithContext(ctx, p.conf.Exchange.Name, routingKey, p.conf.Mandatory, p.conf.Immediate, msg); err != nil {
		return fmt.Errorf("publish rabbitmq message failed: %w", err)
	}
	return nil
}

// PublishBatch 通过指定生产者批量发布消息。
func (m *Manager) PublishBatch(ctx context.Context, name string, msgs ...amqp.Publishing) error {
	if len(msgs) == 0 {
		return ErrInvalidPublishBatch
	}
	for _, msg := range msgs {
		if err := m.Publish(ctx, name, msg); err != nil {
			return err
		}
	}
	return nil
}

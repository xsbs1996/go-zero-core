package xrabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consume 通过指定消费者注册消费处理函数。
func (m *Manager) Consume(ctx context.Context, name string, handler Handler) error {
	if m == nil {
		return ErrNilClient
	}
	if handler == nil {
		return ErrNilHandler
	}

	name, err := normalizeName(name)
	if err != nil {
		return err
	}

	m.mu.RLock()
	c, ok := m.consumers[name]
	m.mu.RUnlock()
	if !ok || c == nil {
		return ErrConsumerNotRegistered
	}
	if c.channel == nil {
		return ErrNilChannel
	}

	queue := c.conf.Queue.Name
	if queue == "" {
		queue = c.conf.Binding.Queue
	}

	msgs, err := c.channel.ConsumeWithContext(ctx, queue, c.conf.Consumer, c.conf.AutoAck, c.conf.Exclusive, c.conf.NoLocal, c.conf.NoWait, nil)
	if err != nil {
		return fmt.Errorf("consume rabbitmq messages failed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			if err := handleDelivery(ctx, msg, c.conf, handler); err != nil {
				return err
			}
		}
	}
}

func handleDelivery(ctx context.Context, msg amqp.Delivery, conf ConsumerConfig, handler Handler) error {
	err := handler(ctx, msg)
	if conf.AutoAck {
		return err
	}
	if err != nil {
		if nackErr := msg.Nack(false, conf.NackRequeue); nackErr != nil {
			return fmt.Errorf("nack rabbitmq message failed: %w", nackErr)
		}
		return err
	}
	if ackErr := msg.Ack(false); ackErr != nil {
		return fmt.Errorf("ack rabbitmq message failed: %w", ackErr)
	}
	return nil
}

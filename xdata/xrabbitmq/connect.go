package xrabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connect 根据配置创建 RabbitMQ 连接。
func Connect(conf Config) (*amqp.Connection, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	conn, err := amqp.DialConfig(conf.DSN(), conf.AMQPConfig())
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq connection failed: %w", err)
	}
	return conn, nil
}

// MustConnect 根据配置创建 RabbitMQ 连接，失败时直接 panic。
func MustConnect(conf Config) *amqp.Connection {
	conn, err := Connect(conf)
	if err != nil {
		panic(err)
	}
	return conn
}

// OpenChannel 根据连接创建 RabbitMQ channel。
func OpenChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	if conn == nil {
		return nil, ErrNilConnection
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq channel failed: %w", err)
	}
	return channel, nil
}

// Ping 检查 RabbitMQ 连接是否可用。
func Ping(conf Config) error {
	conn, err := Connect(conf)
	if err != nil {
		return err
	}
	return conn.Close()
}

package xrabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connect 根据配置创建 RabbitMQ 连接。
//
// 参数：
//   - conf: RabbitMQ 连接配置。
//
// 返回值：
//   - *amqp.Connection: RabbitMQ 连接。
//   - error: 配置校验或连接失败时返回错误。
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
//
// 参数：
//   - conf: RabbitMQ 连接配置。
//
// 返回值：
//   - *amqp.Connection: RabbitMQ 连接。
func MustConnect(conf Config) *amqp.Connection {
	conn, err := Connect(conf)
	if err != nil {
		panic(err)
	}
	return conn
}

// OpenChannel 根据连接创建 RabbitMQ channel。
//
// 参数：
//   - conn: RabbitMQ 连接。
//
// 返回值：
//   - *amqp.Channel: RabbitMQ channel。
//   - error: 连接为空或打开 channel 失败时返回错误。
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
//
// 参数：
//   - conf: RabbitMQ 连接配置。
//
// 返回值：
//   - error: 配置校验、连接或关闭失败时返回错误。
func Ping(conf Config) error {
	conn, err := Connect(conf)
	if err != nil {
		return err
	}
	return conn.Close()
}

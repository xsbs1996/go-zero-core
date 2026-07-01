package rabbitmqx

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Handler RabbitMQ 消费处理函数。
type Handler func(ctx context.Context, msg amqp.Delivery) error

// ExchangeConfig RabbitMQ exchange 声明配置。
type ExchangeConfig struct {
	Name       string     `json:"name,optional" yaml:"name"`             // Name exchange 名称。
	Kind       string     `json:"kind,optional" yaml:"kind"`             // Kind exchange 类型，支持 direct、topic、fanout、headers。
	Durable    bool       `json:"durable,optional" yaml:"durable"`       // Durable 是否持久化。
	AutoDelete bool       `json:"autoDelete,optional" yaml:"autoDelete"` // AutoDelete 是否自动删除。
	Internal   bool       `json:"internal,optional" yaml:"internal"`     // Internal 是否为内部 exchange。
	NoWait     bool       `json:"noWait,optional" yaml:"noWait"`         // NoWait 是否不等待服务端确认。
	Args       amqp.Table `json:"args,optional" yaml:"args"`             // Args 声明参数。
}

// QueueConfig RabbitMQ queue 声明配置。
type QueueConfig struct {
	Name       string     `json:"name,optional" yaml:"name"`             // Name queue 名称。
	Durable    bool       `json:"durable,optional" yaml:"durable"`       // Durable 是否持久化。
	AutoDelete bool       `json:"autoDelete,optional" yaml:"autoDelete"` // AutoDelete 是否自动删除。
	Exclusive  bool       `json:"exclusive,optional" yaml:"exclusive"`   // Exclusive 是否排他。
	NoWait     bool       `json:"noWait,optional" yaml:"noWait"`         // NoWait 是否不等待服务端确认。
	Args       amqp.Table `json:"args,optional" yaml:"args"`             // Args 声明参数。
}

// BindingConfig RabbitMQ queue 绑定配置。
type BindingConfig struct {
	Queue      string     `json:"queue,optional" yaml:"queue"`           // Queue queue 名称。
	Exchange   string     `json:"exchange,optional" yaml:"exchange"`     // Exchange exchange 名称。
	RoutingKey string     `json:"routingKey,optional" yaml:"routingKey"` // RoutingKey 路由键。
	NoWait     bool       `json:"noWait,optional" yaml:"noWait"`         // NoWait 是否不等待服务端确认。
	Args       amqp.Table `json:"args,optional" yaml:"args"`             // Args 绑定参数。
}

// ProducerConfig RabbitMQ 生产者配置。
type ProducerConfig struct {
	Connection Config         `json:"connection" yaml:"connection"`          // Connection RabbitMQ 连接配置。
	Exchange   ExchangeConfig `json:"exchange,optional" yaml:"exchange"`     // Exchange exchange 声明配置。
	Queue      QueueConfig    `json:"queue,optional" yaml:"queue"`           // Queue queue 声明配置。
	Binding    BindingConfig  `json:"binding,optional" yaml:"binding"`       // Binding queue 绑定配置。
	RoutingKey string         `json:"routingKey,optional" yaml:"routingKey"` // RoutingKey 发布消息时使用的路由键。
	Mandatory  bool           `json:"mandatory,optional" yaml:"mandatory"`   // Mandatory 无法路由时是否返回消息。
	Immediate  bool           `json:"immediate,optional" yaml:"immediate"`   // Immediate 是否要求立即投递。
}

// ConsumerConfig RabbitMQ 消费者配置。
type ConsumerConfig struct {
	Connection    Config         `json:"connection" yaml:"connection"`                // Connection RabbitMQ 连接配置。
	Exchange      ExchangeConfig `json:"exchange,optional" yaml:"exchange"`           // Exchange exchange 声明配置。
	Queue         QueueConfig    `json:"queue,optional" yaml:"queue"`                 // Queue queue 声明配置。
	Binding       BindingConfig  `json:"binding,optional" yaml:"binding"`             // Binding queue 绑定配置。
	Consumer      string         `json:"consumer,optional" yaml:"consumer"`           // Consumer consumer 标识。
	AutoAck       bool           `json:"autoAck,optional" yaml:"autoAck"`             // AutoAck 是否自动 ack。
	Exclusive     bool           `json:"exclusive,optional" yaml:"exclusive"`         // Exclusive 是否排他消费。
	NoLocal       bool           `json:"noLocal,optional" yaml:"noLocal"`             // NoLocal 是否不接收当前连接发布的消息。
	NoWait        bool           `json:"noWait,optional" yaml:"noWait"`               // NoWait 是否不等待服务端确认。
	PrefetchCount int            `json:"prefetchCount,optional" yaml:"prefetchCount"` // PrefetchCount 消费者预取数量。
	PrefetchSize  int            `json:"prefetchSize,optional" yaml:"prefetchSize"`   // PrefetchSize 消费者预取大小。
	GlobalQos     bool           `json:"globalQos,optional" yaml:"globalQos"`         // GlobalQos 是否对整个 channel 生效。
	NackRequeue   bool           `json:"nackRequeue,optional" yaml:"nackRequeue"`     // NackRequeue 处理失败时是否重新入队。
}

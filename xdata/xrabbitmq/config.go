package xrabbitmq

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultVHost       = "/"              // 默认虚拟主机。
	defaultLocale      = "en_US"          // 默认 AMQP locale。
	defaultHeartbeat   = 10 * time.Second // 默认心跳间隔。
	defaultDialTimeout = 10 * time.Second // 默认建立连接超时时间。
)

var (
	ErrMissingHost     = errors.New("xrabbitmq: missing rabbitmq host")     // ErrMissingHost 表示 RabbitMQ 地址为空。
	ErrMissingUsername = errors.New("xrabbitmq: missing rabbitmq username") // ErrMissingUsername 表示 RabbitMQ 用户名为空。
)

// Config RabbitMQ 连接配置。
type Config struct {
	URL            string        `json:"url,optional" yaml:"url"`                       // URL RabbitMQ 完整连接地址，配置后优先使用。
	Host           string        `json:"host,optional" yaml:"host"`                     // Host RabbitMQ 地址。
	Port           int           `json:"port,optional" yaml:"port"`                     // Port RabbitMQ 端口。
	Username       string        `json:"username,optional" yaml:"username"`             // Username RabbitMQ 用户名。
	Password       string        `json:"password,optional" yaml:"password"`             // Password RabbitMQ 密码。
	VHost          string        `json:"vHost,default=/" yaml:"vHost"`                  // VHost RabbitMQ 虚拟主机。
	Heartbeat      time.Duration `json:"heartbeat,default=10s" yaml:"heartbeat"`        // Heartbeat 心跳间隔，配置值示例：10s。
	Locale         string        `json:"locale,default=en_US" yaml:"locale"`            // Locale AMQP locale 配置。
	ConnectionName string        `json:"connectionName,optional" yaml:"connectionName"` // ConnectionName RabbitMQ 连接名称。
	DialTimeout    time.Duration `json:"dialTimeout,default=10s" yaml:"dialTimeout"`    // DialTimeout 建立连接超时时间，配置值示例：10s。
	ChannelMax     int           `json:"channelMax,optional" yaml:"channelMax"`         // ChannelMax 最大 channel 数。
	FrameSize      int           `json:"frameSize,optional" yaml:"frameSize"`           // FrameSize AMQP frame 大小。
	TLS            bool          `json:"tls,optional" yaml:"tls"`                       // TLS 是否使用 amqps 协议。
}

// WithDefault 返回补齐默认值后的配置。
func (c Config) WithDefault() Config {
	if c.VHost == "" {
		c.VHost = defaultVHost
	}
	if c.Locale == "" {
		c.Locale = defaultLocale
	}
	if c.Heartbeat == 0 {
		c.Heartbeat = defaultHeartbeat
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultDialTimeout
	}
	return c
}

// Validate 校验 RabbitMQ 连接配置。
func (c Config) Validate() error {
	if c.URL != "" {
		return nil
	}
	if c.Host == "" {
		return ErrMissingHost
	}
	if c.Username == "" {
		return ErrMissingUsername
	}
	return nil
}

// DSN 根据配置生成 RabbitMQ 连接地址。
func (c Config) DSN() string {
	if c.URL != "" {
		return c.URL
	}

	c = c.WithDefault()
	scheme := "amqp"
	if c.TLS {
		scheme = "amqps"
	}

	dsn := url.URL{
		Scheme: scheme,
		User:   url.UserPassword(c.Username, c.Password),
		Host:   c.addr(),
		Path:   "/" + strings.TrimPrefix(c.VHost, "/"),
	}
	if c.VHost == defaultVHost {
		dsn.Path = "/"
	}

	return dsn.String()
}

// AMQPConfig 根据配置生成 RabbitMQ 原生连接配置。
func (c Config) AMQPConfig() amqp.Config {
	c = c.WithDefault()

	conf := amqp.Config{
		Heartbeat:  c.Heartbeat,
		Locale:     c.Locale,
		ChannelMax: uint16(c.ChannelMax),
		FrameSize:  c.FrameSize,
		Properties: amqp.Table{},
	}
	if c.ConnectionName != "" {
		conf.Properties.SetClientConnectionName(c.ConnectionName)
	}
	if c.DialTimeout > 0 {
		dialer := net.Dialer{Timeout: c.DialTimeout}
		conf.Dial = dialer.Dial
	}

	return conf
}

func (c Config) addr() string {
	if c.Port == 0 {
		return c.Host
	}
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

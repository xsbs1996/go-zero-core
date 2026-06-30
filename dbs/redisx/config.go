package redisx

import (
	"errors"
	"time"
)

const (
	defaultPoolSize     = 10
	defaultDialTimeout  = 5 * time.Second
	defaultReadTimeout  = 3 * time.Second
	defaultWriteTimeout = 3 * time.Second
)

var (
	ErrMissingAddr = errors.New("redisx: missing redis addr") // ErrMissingAddr 表示 Redis 地址为空。
)

// Config Redis 连接配置。
type Config struct {
	Addr         string        `json:"addr" yaml:"addr"`                            // Addr Redis 地址，例如 127.0.0.1:6379。
	Username     string        `json:"username,optional" yaml:"username"`           // Username Redis 用户名。
	Password     string        `json:"password,optional" yaml:"password"`           // Password Redis 密码。
	DB           int           `json:"db,optional" yaml:"db"`                       // DB Redis 数据库编号。
	PoolSize     int           `json:"poolSize,default=10" yaml:"poolSize"`         // PoolSize 连接池大小。
	MinIdleConns int           `json:"minIdleConns,optional" yaml:"minIdleConns"`   // MinIdleConns 最小空闲连接数。
	DialTimeout  time.Duration `json:"dialTimeout,default=5s" yaml:"dialTimeout"`   // DialTimeout 建立连接超时时间。
	ReadTimeout  time.Duration `json:"readTimeout,default=3s" yaml:"readTimeout"`   // ReadTimeout 读取超时时间。
	WriteTimeout time.Duration `json:"writeTimeout,default=3s" yaml:"writeTimeout"` // WriteTimeout 写入超时时间。
}

// WithDefault 返回补齐默认值后的配置。
func (c Config) WithDefault() Config {
	if c.PoolSize == 0 {
		c.PoolSize = defaultPoolSize
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultDialTimeout
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	return c
}

// Validate 校验 Redis 连接配置。
func (c Config) Validate() error {
	if c.Addr == "" {
		return ErrMissingAddr
	}
	return nil
}

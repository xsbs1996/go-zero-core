package xredis

import "github.com/redis/go-redis/v9"

type connectOptions struct {
	options *redis.Options
	ping    bool
}

// ConnectOption Redis 连接初始化的可选配置。
type ConnectOption func(*connectOptions)

// WithRedisOptions 使用 go-redis 原生配置。
func WithRedisOptions(options *redis.Options) ConnectOption {
	return func(o *connectOptions) {
		if options != nil {
			o.options = options
		}
	}
}

// WithoutPing 跳过连接后的 Ping 检查。
func WithoutPing() ConnectOption {
	return func(o *connectOptions) {
		o.ping = false
	}
}

func newOptions(conf Config, opts ...ConnectOption) connectOptions {
	options := connectOptions{
		options: redisOptions(conf),
		ping:    true,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return options
}

func redisOptions(conf Config) *redis.Options {
	conf = conf.WithDefault()
	return &redis.Options{
		Addr:         conf.Addr,
		Username:     conf.Username,
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
		DialTimeout:  conf.DialTimeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
}

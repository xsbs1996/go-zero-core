package xpostgres

import "gorm.io/gorm"

type connectOptions struct {
	dsn                    string
	gormOptions            []gorm.Option
	ping                   bool
	shardingPrimaryKeyFunc func(tableIdx int64) int64
}

// ConnectOption 是 PostgreSQL 连接初始化的可选配置。
type ConnectOption func(*connectOptions)

// WithDSN 使用自定义 DSN 连接 PostgreSQL。
func WithDSN(dsn string) ConnectOption {
	return func(o *connectOptions) {
		o.dsn = dsn
	}
}

// WithGormConfig 使用自定义 GORM 配置。
func WithGormConfig(conf *gorm.Config) ConnectOption {
	return func(o *connectOptions) {
		if conf != nil {
			o.gormOptions = []gorm.Option{conf}
		}
	}
}

// WithGormOptions 使用 GORM 原生 Option。
func WithGormOptions(opts ...gorm.Option) ConnectOption {
	return func(o *connectOptions) {
		if len(opts) > 0 {
			o.gormOptions = opts
		}
	}
}

// WithoutPing 跳过连接后的 Ping 检查。
func WithoutPing() ConnectOption {
	return func(o *connectOptions) {
		o.ping = false
	}
}

// WithShardingPrimaryKeyGenerator 使用自定义分表主键生成函数。
func WithShardingPrimaryKeyGenerator(fn func(tableIdx int64) int64) ConnectOption {
	return func(o *connectOptions) {
		o.shardingPrimaryKeyFunc = fn
	}
}

// newOptions 构建连接可选配置。
func newOptions(conf Config, opts ...ConnectOption) connectOptions {
	options := connectOptions{
		dsn:         conf.DSN(),
		gormOptions: []gorm.Option{defaultGormConfig(conf)},
		ping:        true,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return options
}

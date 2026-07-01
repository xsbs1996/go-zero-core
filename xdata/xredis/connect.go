package xredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect 根据配置创建 Redis 客户端。
func Connect(ctx context.Context, conf Config, opts ...ConnectOption) (*redis.Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	options := newOptions(conf, opts...)
	client := redis.NewClient(options.options)

	if options.ping {
		if err := client.Ping(ctx).Err(); err != nil {
			_ = client.Close()
			return nil, fmt.Errorf("ping redis failed: %w", err)
		}
	}

	return client, nil
}

// MustConnect 根据配置创建 Redis 客户端，失败时直接 panic。
func MustConnect(ctx context.Context, conf Config, opts ...ConnectOption) *redis.Client {
	client, err := Connect(ctx, conf, opts...)
	if err != nil {
		panic(err)
	}
	return client
}

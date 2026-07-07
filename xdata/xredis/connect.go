package xredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect 根据配置创建 Redis 客户端。
//
// 参数：
//   - ctx: 控制可选 Ping 检查的上下文。
//   - conf: Redis 连接配置。
//   - opts: 可选连接配置，例如 WithRedisOptions、WithoutPing。
//
// 返回值：
//   - *redis.Client: 创建成功的 Redis 客户端。
//   - error: 配置校验或 Ping 失败时返回错误。
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
//
// 参数：
//   - ctx: 控制可选 Ping 检查的上下文。
//   - conf: Redis 连接配置。
//   - opts: 可选连接配置。
//
// 返回值：
//   - *redis.Client: 创建成功的 Redis 客户端。
func MustConnect(ctx context.Context, conf Config, opts ...ConnectOption) *redis.Client {
	client, err := Connect(ctx, conf, opts...)
	if err != nil {
		panic(err)
	}
	return client
}

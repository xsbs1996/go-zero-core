package xredis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	ErrAlreadyInitialized = errors.New("xredis: client already initialized") // ErrAlreadyInitialized 表示全局 Redis 客户端已经初始化。
	ErrNotInitialized     = errors.New("xredis: client not initialized")     // ErrNotInitialized 表示全局 Redis 客户端未初始化。
)

var (
	redisClient *redis.Client
	redisMu     sync.RWMutex
)

// Init 根据配置初始化全局 Redis 客户端，重复初始化时返回 ErrAlreadyInitialized。
//
// 参数：
//   - ctx: 控制可选 Ping 检查的上下文。
//   - conf: Redis 连接配置。
//   - opts: 可选连接配置。
//
// 返回值：
//   - error: 初始化成功返回 nil；重复初始化或连接失败时返回错误。
func Init(ctx context.Context, conf Config, opts ...ConnectOption) error {
	redisMu.Lock()
	defer redisMu.Unlock()

	if redisClient != nil {
		return ErrAlreadyInitialized
	}

	client, err := Connect(ctx, conf, opts...)
	if err != nil {
		return err
	}

	redisClient = client
	return nil
}

// MustInit 根据配置初始化全局 Redis 客户端，失败时直接 panic。
//
// 参数：
//   - ctx: 控制可选 Ping 检查的上下文。
//   - conf: Redis 连接配置。
//   - opts: 可选连接配置。
func MustInit(ctx context.Context, conf Config, opts ...ConnectOption) {
	if err := Init(ctx, conf, opts...); err != nil {
		panic(err)
	}
}

// SetClient 设置全局 Redis 客户端，已初始化时不会覆盖已有实例。
//
// 参数：
//   - client: 外部创建的 Redis 客户端。
func SetClient(client *redis.Client) {
	if client == nil {
		return
	}

	redisMu.Lock()
	defer redisMu.Unlock()

	if redisClient != nil {
		return
	}

	redisClient = client
}

// GetClient 返回全局 Redis 客户端，未初始化时 panic。
//
// 返回值：
//   - *redis.Client: 全局 Redis 客户端。
func GetClient() *redis.Client {
	redisMu.RLock()
	defer redisMu.RUnlock()

	if redisClient == nil {
		panic("xredis: client not initialized, call xredis.Init or xredis.SetClient first")
	}
	return redisClient
}

// Close 关闭全局 Redis 客户端，并清空全局单例。
//
// 返回值：
//   - error: 关闭成功或未初始化返回 nil；底层客户端关闭失败时返回错误。
func Close() error {
	redisMu.Lock()
	defer redisMu.Unlock()

	if redisClient == nil {
		return ErrNotInitialized
	}

	if err := redisClient.Close(); err != nil {
		return fmt.Errorf("close redis client failed: %w", err)
	}

	redisClient = nil
	return nil
}

// IsInitialized 返回全局 Redis 客户端是否已初始化。
//
// 返回值：
//   - bool: true 表示全局 Redis 客户端已初始化。
func IsInitialized() bool {
	redisMu.RLock()
	defer redisMu.RUnlock()
	return redisClient != nil
}

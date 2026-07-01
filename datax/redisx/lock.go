package redisx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// unlockScript 只在锁 value 匹配时删除 key，避免误删其他实例持有的锁。
	unlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

	// renewScript 只在锁 value 匹配时刷新过期时间，避免续约其他实例持有的锁。
	renewScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0
`
)

var (
	ErrNilClient        = errors.New("redisx: redis client is nil")      // ErrNilClient 表示 Redis 客户端为空。
	ErrMissingLockKey   = errors.New("redisx: missing lock key")         // ErrMissingLockKey 表示锁 key 为空。
	ErrInvalidLockTTL   = errors.New("redisx: invalid lock ttl")         // ErrInvalidLockTTL 表示锁过期时间非法。
	ErrLockNotAcquired  = errors.New("redisx: lock not acquired")        // ErrLockNotAcquired 表示锁未获取成功。
	ErrLockNotOwned     = errors.New("redisx: lock not owned")           // ErrLockNotOwned 表示当前实例不持有该锁。
	ErrLockAlreadyEnded = errors.New("redisx: lock already ended")       // ErrLockAlreadyEnded 表示锁已经结束。
	ErrInvalidRenewal   = errors.New("redisx: invalid renewal interval") // ErrInvalidRenewal 表示续约间隔非法。
)

// Lock Redis 分布式锁。
type Lock struct {
	client        *redis.Client // client Redis 客户端。
	key           string        // key 锁 key。
	value         string        // value 锁唯一值，用于校验锁归属。
	ttl           time.Duration // ttl 锁过期时间。
	renewInterval time.Duration // renewInterval 自动续约间隔。
	stopRenew     chan struct{} // stopRenew 停止续约信号。
	renewDone     chan struct{} // renewDone 续约协程退出信号。
	renewOnce     sync.Once     // renewOnce 保证续约协程只停止一次。
	mu            sync.Mutex    // mu 保护 ended 和 renewErr。
	ended         bool          // ended 标记锁是否已经结束。
	renewErr      error         // renewErr 记录自动续约最后一次错误。
}

// AcquireLock 获取无续约分布式锁，锁到期后由 Redis 自动释放。
//
// ctx 控制本次获取锁请求的生命周期。
// client Redis 客户端。
// key 锁 key，建议使用业务前缀，例如 lock:order:1。
// ttl 锁过期时间，必须大于 0。
func AcquireLock(ctx context.Context, client *redis.Client, key string, ttl time.Duration) (*Lock, error) {
	return acquireLock(ctx, client, key, ttl, 0)
}

// AcquireGlobalLock 使用全局 Redis 客户端获取无续约分布式锁。
//
// ctx 控制本次获取锁请求的生命周期。
// key 锁 key，建议使用业务前缀，例如 lock:order:1。
// ttl 锁过期时间，必须大于 0。
func AcquireGlobalLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	return AcquireLock(ctx, GetClient(), key, ttl)
}

// AcquireRenewalLock 获取自动续约分布式锁，续约失败时可通过 RenewErr 获取错误。
//
// ctx 控制本次获取锁请求的生命周期，不控制后台续约协程。
// client Redis 客户端。
// key 锁 key，建议使用业务前缀，例如 lock:order:1。
// ttl 锁过期时间，必须大于 0。
// renewInterval 自动续约间隔，必须大于 0 且小于 ttl。
func AcquireRenewalLock(ctx context.Context, client *redis.Client, key string, ttl time.Duration, renewInterval time.Duration) (*Lock, error) {
	if ttl <= 0 {
		return nil, ErrInvalidLockTTL
	}
	if renewInterval <= 0 || renewInterval >= ttl {
		return nil, ErrInvalidRenewal
	}

	lock, err := acquireLock(ctx, client, key, ttl, renewInterval)
	if err != nil {
		return nil, err
	}
	lock.startRenewal()
	return lock, nil
}

// AcquireGlobalRenewalLock 使用全局 Redis 客户端获取自动续约分布式锁。
//
// ctx 控制本次获取锁请求的生命周期，不控制后台续约协程。
// key 锁 key，建议使用业务前缀，例如 lock:order:1。
// ttl 锁过期时间，必须大于 0。
// renewInterval 自动续约间隔，必须大于 0 且小于 ttl。
func AcquireGlobalRenewalLock(ctx context.Context, key string, ttl time.Duration, renewInterval time.Duration) (*Lock, error) {
	return AcquireRenewalLock(ctx, GetClient(), key, ttl, renewInterval)
}

// acquireLock 使用 SET NX PX 获取锁。
//
// renewInterval 为 0 时表示无续约锁，大于 0 时表示锁需要后台续约。
func acquireLock(ctx context.Context, client *redis.Client, key string, ttl time.Duration, renewInterval time.Duration) (*Lock, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	if key == "" {
		return nil, ErrMissingLockKey
	}
	if ttl <= 0 {
		return nil, ErrInvalidLockTTL
	}

	value, err := randomLockValue()
	if err != nil {
		return nil, err
	}

	ok, err := client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("acquire redis lock failed: %w", err)
	}
	if !ok {
		return nil, ErrLockNotAcquired
	}

	return &Lock{
		client:        client,
		key:           key,
		value:         value,
		ttl:           ttl,
		renewInterval: renewInterval,
		stopRenew:     make(chan struct{}),
		renewDone:     make(chan struct{}),
	}, nil
}

// Key 返回锁 key。
func (l *Lock) Key() string {
	return l.key
}

// Value 返回锁唯一值，主要用于排查问题，不建议业务依赖。
func (l *Lock) Value() string {
	return l.value
}

// RenewErr 返回自动续约协程最后一次错误；无续约锁始终返回 nil。
func (l *Lock) RenewErr() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.renewErr
}

// Unlock 释放分布式锁，只有锁 value 匹配时才会删除 Redis key。
//
// ctx 控制本次释放锁请求的生命周期。
// 自动续约锁会先停止续约协程，再执行释放逻辑。
func (l *Lock) Unlock(ctx context.Context) error {
	if l == nil || l.client == nil {
		return ErrNilClient
	}

	l.stopRenewal()

	l.mu.Lock()
	if l.ended {
		l.mu.Unlock()
		return ErrLockAlreadyEnded
	}
	l.mu.Unlock()

	result, err := l.client.Eval(ctx, unlockScript, []string{l.key}, l.value).Int64()
	if err != nil {
		return fmt.Errorf("unlock redis lock failed: %w", err)
	}
	if result == 0 {
		l.markEnded(nil)
		return ErrLockNotOwned
	}

	l.markEnded(nil)
	return nil
}

// startRenewal 启动后台续约协程。
//
// 续约协程使用 renewInterval 定时刷新 ttl。
// 一旦续约失败，协程退出并把错误记录到 renewErr。
func (l *Lock) startRenewal() {
	go func() {
		defer close(l.renewDone)

		ticker := time.NewTicker(l.renewInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := l.renew(context.Background()); err != nil {
					l.setRenewErr(err)
					return
				}
			case <-l.stopRenew:
				return
			}
		}
	}()
}

// stopRenewal 停止后台续约协程。
//
// 无续约锁不会启动续约协程，因此这里会直接返回。
func (l *Lock) stopRenewal() {
	if l.renewInterval <= 0 {
		return
	}

	l.renewOnce.Do(func() {
		close(l.stopRenew)
		<-l.renewDone
	})
}

// renew 刷新锁过期时间。
//
// ctx 控制本次续约请求的生命周期。
// 只有锁 value 匹配时才会刷新过期时间。
func (l *Lock) renew(ctx context.Context) error {
	timeout := l.ttl / 2
	if timeout <= 0 || timeout > 3*time.Second {
		timeout = 3 * time.Second
	}

	renewCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := l.client.Eval(renewCtx, renewScript, []string{l.key}, l.value, int64(l.ttl/time.Millisecond)).Int64()
	if err != nil {
		return fmt.Errorf("renew redis lock failed: %w", err)
	}
	if result == 0 {
		return ErrLockNotOwned
	}
	return nil
}

// markEnded 标记锁已经结束。
func (l *Lock) markEnded(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ended = true
	if err != nil {
		l.renewErr = err
	}
}

// setRenewErr 记录自动续约错误。
func (l *Lock) setRenewErr(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.renewErr = err
}

// randomLockValue 生成锁唯一值。
func randomLockValue() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate redis lock value failed: %w", err)
	}
	return hex.EncodeToString(b), nil
}

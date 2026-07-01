# go-zero-core

`go-zero-core` 是面向 go-zero 服务的 Go 基础能力库。它不替代 go-zero、GORM、go-redis、kafka-go 或 RabbitMQ 客户端，而是在这些库之上沉淀微服务项目中常见的初始化、连接管理、中间件、统一响应、日志、加密和类型转换能力。

当前模块名为 `go-zero-core`，要求 Go `1.25+`。

## 适用范围

适合使用本库的场景：

- 多个 go-zero REST 服务需要统一鉴权、CORS、限流、客户端 IP 提取和响应结构。
- 服务需要以一致方式接入 MySQL、PostgreSQL、Redis、Kafka、RabbitMQ。
- 项目需要复用 GORM 日志适配、trace 上下文、结构化日志输出。
- 业务代码中存在大量类型转换、JSON 序列化、时间戳转换、指针取值等重复逻辑。
- 服务需要 JWT、BCrypt、AES、RSA、HMAC、MD5、SHA、UUID、随机值等常用工具。
- 需要 WebSocket 会话管理、连接复用、广播和重连能力，并在业务端按需扩展鉴权、分组、跨实例广播等能力。
- 需要基于 cron 的任务注册、调度和生命周期管理，并在业务端按需扩展持久化、失败重试、分布式调度等能力。

## 版本兼容

| go-zero-core | Go | go-zero |
| --- | --- | --- |
| v1.0.x | 1.25+ | 1.10.1 |

## 安装

```bash
go get go-zero-core
```

如果该仓库作为内部模块维护，请将业务项目中的 import path 替换为实际仓库地址。

## 模块总览

```text
.
├── xcast/             # 类型转换、JSON、时间、map/struct、指针工具
├── xcrypto/           # 加密、摘要、签名、编码、JWT、UUID、随机值
├── xdata/
│   ├── xmysql/        # MySQL GORM 连接、全局实例、sharding
│   ├── xpostgres/     # PostgreSQL GORM 连接、全局实例、sharding
│   ├── xredis/        # Redis 客户端、全局实例、分布式锁
│   ├── xkafka/        # Kafka producer、consumer、manager
│   └── xrabbitmq/     # RabbitMQ 连接、声明、生产、消费、manager
├── xlog/              # go-zero logx 配置、结构化日志、GORM logger
├── xmid/              # go-zero REST 中间件
├── xreply/            # 统一 API 响应和业务错误码
├── xtask/             # cron 定时任务管理
└── xws/               # WebSocket 会话管理
```

## 快速接入

### MySQL

`xdata/xmysql` 使用 GORM 作为底层 ORM。`Connect` 返回错误，适合初始化流程中显式处理失败；`MustConnect` 在失败时 panic，适合服务启动期快速失败。

```go
package main

import "go-zero-core/xdata/xmysql"

func main() {
	db := xmysql.MustConnect(xmysql.Config{
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "password",
		DBName:       "app",
		ParseTime:    true,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogLevel:     "warn",
	})

	_ = db
}
```

核心配置：

- `Host`、`Port`、`User`、`Password`、`DBName`: 基础连接信息。
- `Charset`: 默认 `utf8mb4`。
- `ParseTime`: 是否将 MySQL 时间类型解析为 `time.Time`。
- `Loc`: 默认 `Local`，可设置为 `Asia/Shanghai` 等时区。
- `Timeout`、`ReadTimeout`、`WriteTimeout`: 连接、读、写超时。
- `MaxIdleConns`、`MaxOpenConns`、`ConnMaxLifetime`、`ConnMaxIdleTime`: 连接池配置。
- `SkipDefaultTransaction`、`PrepareStmt`、`LogLevel`: GORM 行为配置。
- `Sharding`: GORM sharding 分表配置。

全局实例适合在项目基础设施初始化阶段统一设置：

```go
if err := xmysql.Init(xmysql.Config{Host: "127.0.0.1", User: "root", DBName: "app"}); err != nil {
	panic(err)
}
db := xmysql.GetDB()
defer xmysql.Close()
```

### PostgreSQL

`xdata/xpostgres` 的使用方式与 `xmysql` 保持一致，底层使用 `gorm.io/driver/postgres`。

```go
db := xpostgres.MustConnect(xpostgres.Config{
	Host:     "127.0.0.1",
	Port:     5432,
	User:     "postgres",
	Password: "password",
	DBName:   "app",
	SSLMode:  "disable",
	TimeZone: "Asia/Shanghai",
	LogLevel: "warn",
})
```

### Redis

`xdata/xredis` 基于 `github.com/redis/go-redis/v9`，支持独立客户端和全局客户端两种方式。

```go
ctx := context.Background()

client := xredis.MustConnect(ctx, xredis.Config{
	Addr:         "127.0.0.1:6379",
	Password:     "",
	DB:           0,
	PoolSize:     20,
	MinIdleConns: 5,
})
defer client.Close()
```

分布式锁使用 Redis `SET NX PX` 获取锁，并通过 Lua 脚本校验锁 value 后释放，避免误删其他实例持有的锁。

```go
lock, err := xredis.AcquireLock(ctx, client, "lock:order:1001", 30*time.Second)
if err != nil {
	return err
}
defer func() {
	_ = lock.Unlock(ctx)
}()

// 执行需要互斥的业务逻辑
```

长任务可使用自动续约锁。`renewInterval` 必须大于 `0` 且小于 `ttl`。

```go
lock, err := xredis.AcquireRenewalLock(ctx, client, "lock:job:sync", time.Minute, 20*time.Second)
if err != nil {
	return err
}
defer lock.Unlock(ctx)

if err := lock.RenewErr(); err != nil {
	return err
}
```

### Kafka

`xdata/xkafka` 基于 `github.com/segmentio/kafka-go`，推荐通过 `Manager` 统一注册和管理多个 topic 的生产者、消费者。

```go
ctx := context.Background()

conf := xkafka.Config{
	Brokers:             []string{"127.0.0.1:9092"},
	ClientID:            "order-service",
	BatchSize:           100,
	BatchTimeout:        1000,
	ConsumeBatchSize:    100,
	ConsumeBatchTimeout: 1000,
}

manager := xkafka.NewManager()
defer manager.Close()

if err := manager.RegisterProducer("order.created", conf); err != nil {
	return err
}
if err := manager.RegisterConsumer("order.created", "order-worker", conf); err != nil {
	return err
}
```

发送单条消息：

```go
err := manager.Produce(ctx, "order.created", kafka.Message{
	Key:   []byte("1001"),
	Value: []byte(`{"orderId":1001}`),
})
if err != nil {
	return err
}
```

批量发送消息：

```go
err := manager.ProduceBatch(ctx, "order.created",
	kafka.Message{Key: []byte("1001"), Value: []byte(`{"orderId":1001}`)},
	kafka.Message{Key: []byte("1002"), Value: []byte(`{"orderId":1002}`)},
)
if err != nil {
	return err
}
```

单条消费会在 handler 成功返回后提交 offset；handler 返回错误时会停止本轮消费并向调用方返回错误。

```go
err := manager.Consume(ctx, "order.created", "order-worker", func(ctx context.Context, msg kafka.Message) error {
	// 处理 msg.Key 和 msg.Value
	return nil
})
if err != nil {
	return err
}
```

批量消费会在达到 `batchSize` 或 `batchTimeout` 后触发 handler。`batchSize` 或 `batchTimeout` 传 `0` 时使用注册消费者时的 `Config` 默认值。

```go
err := manager.ConsumeBatch(ctx, "order.created", "order-worker", 0, 0, func(ctx context.Context, msgs []kafka.Message) error {
	// 批量处理消息
	return nil
})
if err != nil {
	return err
}
```

核心配置：

- `Brokers`: Kafka broker 地址列表，生产者和消费者必填。
- `ClientID`: Kafka client id，用于标识客户端。
- `DialTimeout`、`ReadTimeout`、`WriteTimeout`: 连接、读取和写入超时。
- `BatchSize`、`BatchTimeout`: 生产者批量写入配置，`BatchTimeout` 单位为毫秒。
- `ConsumeBatchSize`、`ConsumeBatchTimeout`: 批量消费配置，`ConsumeBatchTimeout` 单位为毫秒。
- `RequiredAcks`: 写入确认级别，透传 `kafka.RequiredAcks`。
- `Async`: 是否异步写入消息。

### RabbitMQ

`xdata/xrabbitmq` 基于 `github.com/rabbitmq/amqp091-go`，通过命名 producer/consumer 管理连接、channel、exchange、queue、binding、发布和消费。

```go
ctx := context.Background()

connConf := xrabbitmq.Config{
	Host:           "127.0.0.1",
	Port:           5672,
	Username:       "guest",
	Password:       "guest",
	VHost:          "/",
	ConnectionName: "order-service",
}

manager := xrabbitmq.NewManager()
defer manager.Close()
```

注册生产者时可同时声明 exchange、queue 和 binding。`RoutingKey` 为空时，会使用 `Binding.RoutingKey` 作为发布路由键。

```go
producerConf := xrabbitmq.ProducerConfig{
	Connection: connConf,
	Exchange: xrabbitmq.ExchangeConfig{
		Name:    "order.exchange",
		Kind:    "direct",
		Durable: true,
	},
	Queue: xrabbitmq.QueueConfig{
		Name:    "order.created.queue",
		Durable: true,
	},
	Binding: xrabbitmq.BindingConfig{
		Queue:      "order.created.queue",
		Exchange:   "order.exchange",
		RoutingKey: "order.created",
	},
	RoutingKey: "order.created",
}

if err := manager.RegisterProducer("order-producer", producerConf); err != nil {
	return err
}
```

发布消息：

```go
err := manager.Publish(ctx, "order-producer", amqp.Publishing{
	ContentType:  "application/json",
	DeliveryMode: amqp.Persistent,
	Body:         []byte(`{"orderId":1001}`),
})
if err != nil {
	return err
}
```

注册消费者时同样可声明 exchange、queue 和 binding，并可配置 QoS、ack 行为和失败是否重新入队。

```go
consumerConf := xrabbitmq.ConsumerConfig{
	Connection: connConf,
	Exchange: xrabbitmq.ExchangeConfig{
		Name:    "order.exchange",
		Kind:    "direct",
		Durable: true,
	},
	Queue: xrabbitmq.QueueConfig{
		Name:    "order.created.queue",
		Durable: true,
	},
	Binding: xrabbitmq.BindingConfig{
		Queue:      "order.created.queue",
		Exchange:   "order.exchange",
		RoutingKey: "order.created",
	},
	Consumer:      "order-worker",
	PrefetchCount: 10,
	NackRequeue:   true,
}

if err := manager.RegisterConsumer("order-consumer", consumerConf); err != nil {
	return err
}
```

消费消息时，`AutoAck=false` 会在 handler 成功后自动 `Ack`；handler 返回错误时执行 `Nack`，是否重新入队由 `NackRequeue` 控制。

```go
err := manager.Consume(ctx, "order-consumer", func(ctx context.Context, msg amqp.Delivery) error {
	// 处理 msg.Body
	return nil
})
if err != nil {
	return err
}
```

核心配置：

- `Config.URL`: 完整 AMQP 连接地址，配置后优先使用。
- `Config.Host`、`Port`、`Username`、`Password`、`VHost`: 连接信息。
- `ExchangeConfig`: exchange 名称、类型、持久化、自动删除、声明参数。
- `QueueConfig`: queue 名称、持久化、排他、自动删除、声明参数。
- `BindingConfig`: queue、exchange、routing key 绑定关系。
- `ProducerConfig.Mandatory`、`Immediate`: 发布行为参数。
- `ConsumerConfig.AutoAck`、`PrefetchCount`、`NackRequeue`: 消费确认、预取和失败处理策略。

### JWT

`xcrypto/xjwt` 用于生成、解析和刷新 JWT。业务 claims 放在 `Data` 中。

```go
conf := xjwt.Config{
	Secret: "your-secret",
	Expire: time.Hour,
}

token, err := xjwt.Generate(conf, map[string]any{
	"userId": 1001,
})
if err != nil {
	return err
}

claims, err := xjwt.Parse(conf, token)
if err != nil {
	return err
}

fmt.Println(claims.Data)
```

### REST 中间件

`xmid` 提供 go-zero REST 中间件。中间件遵循 `rest.Middleware` 类型，可直接通过 `server.Use` 注册。

```go
server := rest.MustNewServer(rest.RestConf{})
defer server.Stop()

server.Use(xmid.Auth(xmid.AuthConfig{
	SkipPaths: []string{"/healthz", "/login"},
	Verify: func(r *http.Request, token string) (any, error) {
		claims, err := xjwt.Parse(jwtConf, token)
		if err != nil {
			return nil, err
		}
		return claims.Data, nil
	},
}))
```

在 handler 中读取认证信息：

```go
authInfo, ok := xmid.AuthInfo(r.Context())
if !ok {
	xreply.FailStatus(w, http.StatusUnauthorized, xreply.CodeUnauthorized)
	return
}
_ = authInfo
```

可用中间件：

- `Auth`: Bearer token 鉴权，支持自定义 header、prefix、跳过路径、校验函数和未授权响应。
- `Cors`: CORS 处理，支持来源、方法、请求头、暴露响应头、凭证和预检缓存时间。
- `RateLimit`: 进程内窗口限流，适合单实例或轻量保护场景。
- `IP`: 将客户端 IP 写入请求上下文。
- `ClientIP`: 独立 IP 提取函数，可在非中间件场景使用。

### 统一响应

`xreply` 固定响应结构为：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

基本用法：

```go
xreply.Success(w, user)
xreply.SuccessPage(w, users, total, page, pageSize)
xreply.Fail(w, xreply.CodeInvalidParam, xreply.Vars{"field": "name"})
xreply.FailStatus(w, http.StatusUnauthorized, xreply.CodeUnauthorized)
```

错误码约定：

- `0`: 成功。
- `1-99`: 公共错误码，由 `xreply` 保留。
- `100+`: 业务项目自定义错误码，可通过 `RegisterCodes` 注册。

```go
const CodeOrderNotFound = 10001

xreply.RegisterCodes(map[int]string{
	CodeOrderNotFound: "order not found",
})
```

### WebSocket

`xws` 封装 WebSocket 升级、会话存储、读写通道、广播和同 code 热重连。`code` 是业务传入的连接标识，库本身不绑定用户、设备或权限模型。

```go
manager := xws.NewManager(xws.Config{
	MaxConnTotal:  10000,
	ReadChanSize:  1024,
	WriteChanSize: 1024,
})

func wsHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	session, isNew, err := manager.Create(w, r, code)
	if err != nil {
		return
	}

	if isNew {
		go func() {
			for msg := range session.Read() {
				_ = session.Write(msg)
			}
		}()
	}
}
```

常用方法：

- `Create(w, r, code)`: 创建或复用指定 code 的会话。
- `Get(code)`: 获取会话。
- `CloseConn(code)`: 关闭指定会话。
- `Count()`: 当前在线会话数。
- `Range(fn)`: 遍历会话快照。
- `Broadcast(msg)`: 向所有在线会话写入消息。

### 定时任务

`xtask` 基于 `github.com/robfig/cron/v3`，用于进程内定时任务。

```go
manager := xtask.NewManager(xtask.WithSeconds())

err := manager.AddFuncNow("refresh-cache", "*/30 * * * * *", func(ctx context.Context) {
	// 执行业务逻辑
})
if err != nil {
	panic(err)
}

manager.Start()
defer manager.Stop()
```

说明：

- `WithSeconds()` 启用秒级 cron 表达式。
- `WithLocation(location)` 指定调度时区。
- 同名任务重复注册时，新任务会替换旧任务。
- `Stop()` 会停止调度器，并等待已经运行中的任务退出。

## 包级 API 速查

### `xcast`

- 字符串、整数、无符号整数、浮点数、布尔值互转。
- `any` 转字符串、整数、`int64`、`float64`、布尔值，并提供默认值版本。
- `time.Time` 与 Unix 秒、毫秒时间戳互转。
- JSON marshal、unmarshal、格式化输出。
- struct 与 map 转换。
- 泛型指针工具：创建指针、取值、默认值。

### `xcrypto`

- `xaes`: AES-GCM、AES-CBC 加解密。
- `xrsa`: RSA 密钥生成、PEM 编解码、OAEP 加解密、PSS 签名和验签。
- `xjwt`: JWT 生成、解析、刷新。
- `xbcrypt`: 密码哈希和校验。
- `xhmac`: HMAC-SHA1、HMAC-SHA256、HMAC-SHA512 签名。
- `xmd5`: 字符串、字节切片、文件 MD5。
- `xsha`: SHA1、SHA256、SHA384、SHA512。
- `xbase64`: Base64、Base64URL、Base62、Base58。
- `xuuid`: UUID 生成、解析、校验、去横线 UUID。
- `xrand`: 随机字节、十六进制字符串、数字字符串、字母字符串、混合字符串。

### `xdata`

- `xmysql`: `Connect`、`MustConnect`、`Init`、`MustInit`、`GetDB`、`SetDB`、`Close`、sharding。
- `xpostgres`: PostgreSQL 版本的 GORM 连接与全局实例管理。
- `xredis`: `Connect`、`MustConnect`、`Init`、`MustInit`、`GetClient`、`SetClient`、`Close`、分布式锁。
- `xkafka`: producer、consumer、topic/group manager、单条和批量生产消费。
- `xrabbitmq`: 连接、channel、exchange/queue/binding 声明、producer、consumer、manager。

### `xlog`

- 复用 go-zero `logx`，支持 console、file、volume 输出模式。
- 支持 json、plain 编码，debug、info、error、severe 级别。
- 提供 `Debug`、`Info`、`Warn`、`Error`、`ErrorStack`、`Severe`。
- 支持 context trace 注入、禁用日志输出。
- 提供 GORM logger 适配，统一 SQL 日志结构。

### `xmid`

- `Auth`: 请求鉴权。
- `AuthInfo`: 从 context 获取鉴权信息。
- `Cors`: 跨域处理。
- `RateLimit`: 进程内限流。
- `IP`、`ClientIP`: 客户端 IP 提取。

### `xreply`

- `Success`: 成功响应。
- `SuccessPage`: 分页成功响应。
- `SuccessMsg`: 自定义成功消息。
- `Fail`: 业务失败响应，HTTP 状态码为 200。
- `FailStatus`: 指定 HTTP 状态码的失败响应。
- `RegisterCodes`: 注册业务错误码。

### `xws`

- `Manager`: 会话创建、获取、关闭、遍历、广播。
- `Session`: `Code`、`Read`、`Write`、`Close`、`IsAlive`。
- 支持相同 code 重连，旧连接退出不会关闭新会话。

### `xtask`

- `NewManager`: 创建任务管理器。
- `Add`、`AddFunc`、`AddFuncNow`: 注册任务。
- `Start`、`Stop`: 启动和停止调度器。
- `Remove`、`Count`、`Entries`: 任务管理和快照查询。

## 接入建议

- 服务启动期优先初始化数据库、Redis、消息队列和日志，失败时直接终止进程，避免服务以不完整状态运行。
- 业务代码中优先传递显式客户端或 `*gorm.DB`；全局实例更适合基础设施较简单的服务。
- Redis 分布式锁必须始终设置合理 TTL，并在业务结束后调用 `Unlock`。长任务使用自动续约锁时，应定期检查 `RenewErr`。
- `xmid.RateLimit` 是进程内限流，不适合多实例全局限流。
- WebSocket 的 `code` 应由业务层完成鉴权和唯一性设计，本库只负责连接生命周期。
- `xreply` 的 `100+` 错误码建议由业务模块集中注册，避免跨模块冲突。

## 开发

```bash
go mod tidy
go test ./...
```

## 主要依赖

- `github.com/zeromicro/go-zero`
- `gorm.io/gorm`
- `gorm.io/driver/mysql`
- `gorm.io/driver/postgres`
- `gorm.io/sharding`
- `github.com/redis/go-redis/v9`
- `github.com/segmentio/kafka-go`
- `github.com/rabbitmq/amqp091-go`
- `github.com/gorilla/websocket`
- `github.com/robfig/cron/v3`
- `github.com/golang-jwt/jwt/v5`

## 许可证

请查看 [LICENSE](./LICENSE)。

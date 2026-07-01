# go-zero-core

`go-zero-core` 是一个面向 go-zero 项目的 Go 基础能力库，主要用于沉淀微服务项目中反复使用的通用组件。项目围绕 go-zero REST 服务、GORM、Redis、Kafka、RabbitMQ、JWT、日志和基础类型转换做轻量封装，目标是在业务项目中减少重复样板代码，并统一常见基础设施的接入方式。

当前模块名为 `go-zero-core`，Go 版本要求为 `1.25+`。

## 版本适配

| go-zero-core 版本 | go-zero 版本 |
| --- | --- |
| v1.0 系列 | v1.10.1 |

## 项目定位

这个仓库更适合作为业务服务的内部基础库使用，而不是完整的应用框架。它不接管 go-zero 的工程结构，也不替代 go-zero、GORM、go-redis、kafka-go 或 RabbitMQ 客户端本身，而是在这些库之上提供更贴近业务项目的初始化、配置、封装和辅助函数。

适用场景：

- go-zero REST 服务需要统一鉴权、CORS、限流、IP 提取等中间件。
- 多个服务需要复用 MySQL、PostgreSQL、Redis、Kafka、RabbitMQ 的连接管理方式。
- 项目中需要统一日志内容格式、trace 上下文和 GORM 日志输出。
- 业务代码中频繁出现类型转换、JSON 序列化、时间转换、指针取值等基础操作。
- 需要常用加密、签名、摘要、JWT、UUID、随机字符串等工具函数。

## 目录结构

```text
.
├── xcast/                   # 类型转换和序列化辅助工具
├── xdata/                   # 数据库、缓存和消息队列封装
│   ├── xkafka/              # Kafka 生产者、消费者和客户端管理
│   ├── xmysql/              # MySQL GORM 连接、全局实例和分表支持
│   ├── xpostgres/           # PostgreSQL GORM 连接、全局实例和分表支持
│   ├── xrabbitmq/           # RabbitMQ 连接、声明、生产者、消费者和客户端管理
│   └── xredis/              # Redis 连接、全局客户端和分布式锁
├── xcrypto/                 # 加密、摘要、编码、JWT、UUID 和随机值工具
│   ├── xaes/                # AES-GCM、AES-CBC 加解密
│   ├── xbase64/             # Base64、Base64URL、Base62、Base58 编解码
│   ├── xbcrypt/             # BCrypt 密码哈希和校验
│   ├── xhmac/               # HMAC-SHA 签名和校验
│   ├── xjwt/                # JWT 生成、解析和刷新
│   ├── xmd5/                # MD5 字符串、字节和文件摘要
│   ├── xrand/               # 随机字节、十六进制、数字和字符串
│   ├── xrsa/                # RSA 密钥、OAEP 加解密和 PSS 签名
│   ├── xsha/                # SHA1、SHA256、SHA384、SHA512 摘要
│   └── xuuid/               # UUID 生成、解析和校验
├── xlog/                    # go-zero logx 配置、结构化日志和 GORM 日志适配
├── xmid/                    # go-zero REST 中间件
├── xreply/                  # 统一 API 响应结构和响应辅助能力
├── xtask/                   # 定时任务管理
├── xws/                     # WebSocket 会话和连接管理
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖版本锁定
├── LICENSE                  # 许可证
└── README.md                # 项目说明文档
```

## 功能介绍

### xcast

`xcast` 包提供基础类型转换，适合替代业务代码中分散的 `strconv`、`json.Marshal`、`time.Format` 等重复逻辑。

主要能力：

- 字符串转整数、无符号整数、浮点数、布尔值和时间。
- 整数、无符号整数、浮点数、布尔值转字符串。
- `any` 类型转字符串、整数、`int64`、`float64`、布尔值，并提供默认值版本。
- `time.Time` 和 Unix 秒、毫秒时间戳互转。
- 标准 JSON 和格式化 JSON 序列化、反序列化。
- struct 和 map 之间转换。
- 泛型指针工具：创建指针、取指针值、取默认值。

### xdata/xmysql

`xdata/xmysql` 基于 GORM 和 `gorm.io/driver/mysql` 封装 MySQL 连接。

主要能力：

- 使用结构化 `Config` 生成 MySQL DSN。
- 支持连接超时、读写超时、连接池、连接生命周期、GORM 日志级别等配置。
- 支持 `Connect`、`MustConnect` 创建独立连接。
- 支持 `Init`、`MustInit`、`GetDB`、`SetDB`、`Close` 管理全局 GORM 实例。
- 支持通过 `ConnectOption` 传入自定义 DSN、GORM 配置、GORM options，以及跳过 ping。
- 支持 GORM sharding 分表配置和自定义分表主键生成器。

### xdata/xpostgres

`xdata/xpostgres` 基于 GORM 和 `gorm.io/driver/postgres` 封装 PostgreSQL 连接，整体使用方式和 `xmysql` 保持一致。

主要能力：

- 使用结构化 `Config` 生成 PostgreSQL DSN。
- 支持 host、port、user、password、database、sslmode、timezone 等连接参数。
- 支持连接池、连接生命周期、GORM 日志级别配置。
- 支持独立连接创建和全局 GORM 实例管理。
- 支持 GORM sharding 分表配置。

### xdata/xredis

`xdata/xredis` 基于 `github.com/redis/go-redis/v9` 封装 Redis 客户端。

主要能力：

- 支持 Redis 地址、用户名、密码、DB、连接池、超时配置。
- 支持 `Connect`、`MustConnect` 创建客户端。
- 支持 `Init`、`MustInit`、`GetClient`、`SetClient`、`Close` 管理全局 Redis 客户端。
- 支持自定义 `redis.Options` 和跳过 ping。
- 提供基于 Redis 的分布式锁。
- 分布式锁支持普通锁和自动续期锁，并通过 Lua 脚本保证释放锁时校验锁值。

### xdata/xkafka

`xdata/xkafka` 基于 `github.com/segmentio/kafka-go` 封装 Kafka 生产和消费。

主要能力：

- 支持 broker、client id、读写超时、批量大小、批量等待时间、RequiredAcks、异步写入等配置。
- 支持创建生产者和消费者。
- 支持按 topic 注册和管理生产者。
- 支持按 topic + group 注册和管理消费者。
- 支持单条生产、批量生产。
- 支持单条消费和批量消费。
- 提供默认客户端和自定义 Manager，便于在服务中集中管理多个 topic。

### xdata/xrabbitmq

`xdata/xrabbitmq` 基于 `github.com/rabbitmq/amqp091-go` 封装 RabbitMQ 连接和消息收发。

主要能力：

- 支持 RabbitMQ host、port、username、password、vhost、heartbeat、locale、TLS 等连接配置。
- 支持连接、打开 channel 和 ping。
- 支持 exchange、queue、binding 的声明配置。
- 支持注册生产者和消费者。
- 支持按名称发布消息和批量发布消息。
- 支持消费消息，并根据 handler 返回值执行 ack、nack 或 reject。
- 提供默认客户端和自定义 Manager。

### xcrypto

`xcrypto` 目录按算法拆分为多个子包，覆盖业务项目中常见的加密、签名、摘要、编码和令牌场景。

主要能力：

- `xaes`: AES-GCM、AES-CBC 加解密。
- `xrsa`: RSA 密钥生成、PEM 编解码、OAEP 加解密、PSS 签名和验签。
- `xjwt`: JWT token 生成、解析、刷新，支持 issuer、subject、audience、expire 等配置。
- `xbcrypt`: 密码哈希和密码校验。
- `xhmac`: HMAC-SHA1、HMAC-SHA256、HMAC-SHA512 签名，支持十六进制和 Base64 输出。
- `xmd5`: 字符串、字节切片和文件 MD5 摘要。
- `xsha`: SHA1、SHA256、SHA384、SHA512 摘要。
- `xbase64`: Base64、Base64URL、Base62、Base58 编解码。
- `xuuid`: UUID 生成、去横线 UUID、解析和格式校验。
- `xrand`: 随机字节、十六进制字符串、数字字符串、字母字符串、混合字符串。

### xlog

`xlog` 包基于 go-zero `logx` 做日志配置和结构化内容封装，同时提供 GORM 日志适配。

主要能力：

- 复用 go-zero `logx.LogConf`、`logx.LogField`、`logx.Writer`。
- 提供更直接的 `Config`，支持 console、file、volume 输出模式。
- 支持 json、plain 编码格式。
- 支持 debug、info、error、severe 日志级别。
- 支持日志路径、保留天数、压缩、切割策略、最大文件大小、最大备份数量等配置。
- 支持统一日志 body：`msg`、`content`、`err`。
- 提供 `Debug`、`Info`、`Warn`、`Error`、`ErrorStack`、`Severe` 等方法。
- 支持在 context 中禁用日志输出。
- 支持为 context 注入 trace 信息。
- 提供 GORM logger 适配，统一 SQL 日志输出格式。

### xmid

`xmid` 包提供 go-zero REST 中间件。

主要能力：

- `Auth`: Bearer token 鉴权中间件，支持自定义 header、prefix、跳过路径、校验函数和未授权响应。
- `AuthInfo`: 从请求 context 中读取鉴权结果。
- `Cors`: CORS 中间件，支持允许来源、方法、请求头、暴露响应头、凭证和预检缓存时间。
- `RateLimit`: 基于内存的简单限流中间件，支持窗口时间、请求次数、key 函数和超限响应。
- `IP`: 客户端 IP 提取中间件，支持从指定 header 和 `RemoteAddr` 获取 IP。
- `ClientIP`: 独立 IP 提取函数，可在非中间件场景复用。

### xreply

`xreply` 包用于统一 API 响应结构、公共错误码和错误响应文案，并基于 go-zero `httpx.OkJson` 直接向客户端输出 JSON。响应字段固定为 `code`、`msg`、`data`。

主要能力：

- `0-99`: 公共保留错误码区间，其中 `0` 表示成功，`1-99` 表示常用错误。
- `0`: 成功码，同样维护在内置 code map 中。
- `RegisterCodes`: 注册业务项目自定义错误码，自定义 code 必须从 `100` 开始。
- `Vars`: 失败 msg 变量，支持 `{name}` 形式的占位符替换。
- 内置 `CodeInvalidParam` 默认文案包含 `{field}`，例如传入 `xreply.Vars{"field": "name"}` 后输出 `invalid param: name`。
- `Success`: 输出成功响应，固定使用 `code=0`，支持可选 msg 变量。
- `Fail`: 根据 code 输出失败响应，msg 从错误码表自动填充。
- `FailStatus`: 根据 HTTP status 和 code 输出失败响应，msg 从错误码表自动填充。
- `SuccessPage`: 输出分页成功响应。
- `SuccessMsg`: 输出指定 code 的成功响应，并使用自定义 msg 覆盖 code map 中的默认 msg。

### xws

`xws` 包用于 WebSocket 会话和连接管理，适合在 go-zero REST handler 中接入升级后的长连接。

主要能力：

- `Manager`: 管理 WebSocket 会话创建、复用、关闭和遍历。
- `Session`: 封装单个连接编码对应的 WebSocket 连接、读写通道和生命周期。
- `Config`: 支持最大连接数、读写缓冲区、读写通道长度、超时时间、消息类型和来源校验配置。
- 创建会话时显式传入连接编码 `code`，不绑定具体用户、设备或鉴权概念。
- 支持同一连接编码热重连，替换连接时不会因为旧连接退出而关闭新会话。
- 支持 `Get`、`CloseConn`、`Count`、`Range`、`Broadcast` 等常用管理能力。
- 日志统一使用本库 `xlog` 包输出。

最小用法：

```go
manager := xws.NewManager()
session, isNew, err := manager.Create(w, r, code)
```

### xtask

`xtask` 包用于管理基于 cron 表达式的定时任务。

主要能力：

- `Manager`: 管理定时任务注册、启动、停止和移除。
- `Job`: 描述任务名称、cron 表达式、是否立即执行和执行函数。
- 支持秒级 cron 表达式。
- 支持指定调度时区。
- 支持同名任务覆盖注册。

最小用法：

```go
manager := xtask.NewManager(xtask.WithSeconds())
_ = manager.AddFunc("refresh-cache", "*/30 * * * * *", func(ctx context.Context) {
	// 执行业务逻辑
})
manager.Start()
defer manager.Stop()
```

## 安装

```bash
go get go-zero-core
```

如果作为内部模块使用，请根据实际仓库地址调整 `go.mod` 中的 module path 和业务项目的 import path。

## 使用示例

### JWT

```go
package main

import (
	"fmt"
	"time"

	"go-zero-core/xcrypto/xjwt"
)

func main() {
	conf := xjwt.Config{
		Secret: "your-secret",
		Expire: time.Hour,
	}

	token, err := xjwt.Generate(conf, map[string]any{
		"userId": 1001,
	})
	if err != nil {
		panic(err)
	}

	claims, err := xjwt.Parse(conf, token)
	if err != nil {
		panic(err)
	}

	fmt.Println(claims.Data)
}
```

### MySQL

```go
package main

import "go-zero-core/xdata/xmysql"

func main() {
	db := xmysql.MustConnect(xmysql.Config{
		Host:      "127.0.0.1",
		Port:      3306,
		User:      "root",
		Password:  "password",
		DBName:    "app",
		ParseTime: true,
	})

	_ = db
}
```

### Redis 分布式锁

```go
package main

import (
	"context"
	"time"

	"go-zero-core/xdata/xredis"
)

func main() {
	ctx := context.Background()
	client := xredis.MustConnect(ctx, xredis.Config{
		Addr: "127.0.0.1:6379",
	})
	defer client.Close()

	lock, err := xredis.AcquireLock(ctx, client, "lock:job", 30*time.Second)
	if err != nil {
		panic(err)
	}
	defer lock.Unlock(ctx)

	// 执行业务逻辑
}
```

### 鉴权中间件

```go
package main

import (
	"net/http"

	"go-zero-core/xmid"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	server := rest.MustNewServer(rest.RestConf{})
	defer server.Stop()

	server.Use(xmid.Auth(xmid.AuthConfig{
		Verify: func(r *http.Request, token string) (any, error) {
			return token, nil
		},
	}))

	server.Start()
}
```

## 开发命令

```bash
go mod tidy
go test ./...
```

## 依赖

主要依赖：

- `github.com/zeromicro/go-zero`
- `gorm.io/gorm`
- `gorm.io/driver/mysql`
- `gorm.io/driver/postgres`
- `gorm.io/sharding`
- `github.com/redis/go-redis/v9`
- `github.com/segmentio/kafka-go`
- `github.com/rabbitmq/amqp091-go`
- `github.com/golang-jwt/jwt/v5`

## 许可证

请查看 [LICENSE](./LICENSE)。

package xws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096
	defaultMaxConnTotal    = 10240
	defaultChannelSize     = 1024
	defaultReadDeadline    = 120 * time.Minute
	defaultWriteDeadline   = 10 * time.Second
)

// Config 表示 WebSocket 会话管理配置。
type Config struct {
	MaxConnTotal    int32                      // MaxConnTotal 表示最大在线会话数量
	ReadBufferSize  int                        // ReadBufferSize 表示 WebSocket 读缓冲区大小
	WriteBufferSize int                        // WriteBufferSize 表示 WebSocket 写缓冲区大小
	ReadChanSize    int                        // ReadChanSize 表示上行消息通道长度
	WriteChanSize   int                        // WriteChanSize 表示下行消息通道长度
	ReadDeadline    time.Duration              // ReadDeadline 表示单次读取超时时间
	WriteDeadline   time.Duration              // WriteDeadline 表示单次写入超时时间
	MessageType     int                        // MessageType 表示下行消息类型
	CheckOrigin     func(r *http.Request) bool // CheckOrigin 表示 WebSocket 来源校验函数
}

// DefaultConfig 返回默认配置。
//
// 返回值：
//   - Config: 已填充默认缓冲区、通道长度、读写超时、消息类型和来源校验函数的配置。
func DefaultConfig() Config {
	return Config{
		MaxConnTotal:    defaultMaxConnTotal,
		ReadBufferSize:  defaultReadBufferSize,
		WriteBufferSize: defaultWriteBufferSize,
		ReadChanSize:    defaultChannelSize,
		WriteChanSize:   defaultChannelSize,
		ReadDeadline:    defaultReadDeadline,
		WriteDeadline:   defaultWriteDeadline,
		MessageType:     websocket.BinaryMessage,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func normalizeConfig(config Config) Config {
	defaultConfig := DefaultConfig()

	if config.MaxConnTotal <= 0 {
		config.MaxConnTotal = defaultConfig.MaxConnTotal
	}
	if config.ReadBufferSize <= 0 {
		config.ReadBufferSize = defaultConfig.ReadBufferSize
	}
	if config.WriteBufferSize <= 0 {
		config.WriteBufferSize = defaultConfig.WriteBufferSize
	}
	if config.ReadChanSize <= 0 {
		config.ReadChanSize = defaultConfig.ReadChanSize
	}
	if config.WriteChanSize <= 0 {
		config.WriteChanSize = defaultConfig.WriteChanSize
	}
	if config.ReadDeadline <= 0 {
		config.ReadDeadline = defaultConfig.ReadDeadline
	}
	if config.WriteDeadline <= 0 {
		config.WriteDeadline = defaultConfig.WriteDeadline
	}
	if config.MessageType == 0 {
		config.MessageType = defaultConfig.MessageType
	}
	if config.CheckOrigin == nil {
		config.CheckOrigin = defaultConfig.CheckOrigin
	}

	return config
}

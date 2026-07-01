package xws

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"go-zero-core/xlog"

	"github.com/gorilla/websocket"
)

// Session 表示一个 WebSocket 连接会话
type Session struct {
	ctx    context.Context    // ctx 控制会话读写协程退出
	cancel context.CancelFunc // cancel 主动取消会话上下文
	code   string             // code 表示连接编码

	connMu  sync.RWMutex    // connMu 保护 conn 热替换过程
	conn    *websocket.Conn // conn 表示当前生效的 WebSocket 连接
	connSeq atomic.Uint64   // connSeq 表示连接版本，用于区分新旧连接

	readCh  chan []byte // readCh 接收客户端上行消息
	writeCh chan []byte // writeCh 发送服务端下行消息

	closeOnce sync.Once    // closeOnce 确保会话只关闭一次
	state     atomic.Int32 // state 表示会话状态，1 为运行中，0 为已关闭

	readDeadline  time.Duration          // readDeadline 表示单次读取超时时间
	writeDeadline time.Duration          // writeDeadline 表示单次写入超时时间
	messageType   int                    // messageType 表示下行消息类型
	onClose       func(session *Session) // onClose 表示会话关闭后的回调
}

func createSession(conn *websocket.Conn, code string, config Config, onClose func(session *Session)) *Session {
	ctx, cancel := context.WithCancel(xlog.ContextWithTrace(context.Background()))

	session := &Session{
		ctx:           ctx,
		cancel:        cancel,
		code:          code,
		conn:          conn,
		readCh:        make(chan []byte, config.ReadChanSize),
		writeCh:       make(chan []byte, config.WriteChanSize),
		readDeadline:  config.ReadDeadline,
		writeDeadline: config.WriteDeadline,
		messageType:   config.MessageType,
		onClose:       onClose,
	}

	session.state.Store(1)
	session.connSeq.Store(1)

	go session.readLoop(conn, 1)
	go session.writeLoop()

	return session
}

// Code 返回会话连接编码
func (s *Session) Code() string {
	if s == nil {
		return ""
	}
	return s.code
}

// Close 关闭会话
func (s *Session) Close() {
	if s == nil {
		return
	}

	s.closeOnce.Do(func() {
		s.state.Store(0)
		s.cancel()

		s.connMu.Lock()
		if s.conn != nil {
			_ = s.conn.Close()
			s.conn = nil
		}
		s.connMu.Unlock()

		if s.onClose != nil {
			s.onClose(s)
		}

		xlog.Info(s.ctx, "websocket session closed", map[string]any{"code": s.code})
	})
}

// replaceConn 热替换当前连接，旧连接退出不会关闭新会话
func (s *Session) replaceConn(newConn *websocket.Conn) {
	s.connMu.Lock()
	oldConn := s.conn
	s.conn = newConn
	seq := s.connSeq.Add(1)
	s.connMu.Unlock()

	go s.readLoop(newConn, seq)

	if oldConn != nil {
		_ = oldConn.Close()
	}

	xlog.Info(s.ctx, "websocket connection replaced", map[string]any{"code": s.code})
}

// readLoop 持续读取客户端上行消息
func (s *Session) readLoop(conn *websocket.Conn, seq uint64) {
	defer func() {
		if value := recover(); value != nil {
			xlog.Error(s.ctx, "websocket readLoop panic", map[string]any{
				"code":  s.code,
				"panic": value,
				"stack": string(debug.Stack()),
			}, nil)
			if s.isCurrentConn(conn, seq) {
				s.Close()
			}
		}
		xlog.Info(s.ctx, "websocket readLoop exited", map[string]any{"code": s.code, "seq": seq})
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(s.readDeadline))

		_, data, err := conn.ReadMessage()
		if err != nil {
			if s.isCurrentConn(conn, seq) {
				xlog.Error(s.ctx, "websocket read message failed", map[string]any{"code": s.code}, err)
				s.Close()
			}
			return
		}

		select {
		case s.readCh <- data:
		default:
			xlog.Warn(s.ctx, "websocket read channel full", map[string]any{"code": s.code, "len": len(data)}, nil)
		}
	}
}

// writeLoop 持续发送服务端下行消息
func (s *Session) writeLoop() {
	defer func() {
		if value := recover(); value != nil {
			xlog.Error(s.ctx, "websocket writeLoop panic", map[string]any{
				"code":  s.code,
				"panic": value,
				"stack": string(debug.Stack()),
			}, nil)
			s.Close()
		}
		xlog.Info(s.ctx, "websocket writeLoop exited", map[string]any{"code": s.code})
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.writeCh:
			conn, seq := s.currentConn()
			if conn == nil {
				continue
			}

			_ = conn.SetWriteDeadline(time.Now().Add(s.writeDeadline))
			if err := conn.WriteMessage(s.messageType, msg); err != nil {
				xlog.Error(s.ctx, "websocket write message failed", map[string]any{"code": s.code, "len": len(msg)}, err)
				if s.isCurrentConn(conn, seq) {
					s.Close()
				}
			}
		}
	}
}

// Write 写入下行消息
func (s *Session) Write(msg []byte) bool {
	if s == nil || !s.IsAlive() {
		return false
	}

	select {
	case <-s.ctx.Done():
		return false
	case s.writeCh <- msg:
		return true
	default:
		select {
		case <-s.writeCh:
		default:
		}

		select {
		case s.writeCh <- msg:
			return true
		default:
			return false
		}
	}
}

// IsAlive 判断会话是否存活
func (s *Session) IsAlive() bool {
	return s != nil && s.state.Load() == 1
}

// Read 返回上行消息读取通道
func (s *Session) Read() <-chan []byte {
	if s == nil {
		return nil
	}
	return s.readCh
}

// currentConn 返回当前连接和连接版本
func (s *Session) currentConn() (*websocket.Conn, uint64) {
	s.connMu.RLock()
	conn := s.conn
	seq := s.connSeq.Load()
	s.connMu.RUnlock()
	return conn, seq
}

// isCurrentConn 判断连接和版本是否仍是当前生效连接
func (s *Session) isCurrentConn(conn *websocket.Conn, seq uint64) bool {
	s.connMu.RLock()
	current := s.conn == conn && s.connSeq.Load() == seq
	s.connMu.RUnlock()
	return current
}

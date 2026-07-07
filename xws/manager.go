package xws

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/xsbs1996/go-zero-core/xlog"

	"github.com/gorilla/websocket"
)

// Manager 管理 WebSocket 会话生命周期。
type Manager struct {
	mu       sync.RWMutex        // mu 保护 sessions 的并发读写
	sessions map[string]*Session // sessions 保存连接编码到会话的映射
	total    atomic.Int32        // total 表示当前在线会话数量
	config   Config              // config 表示管理器配置
	upgrader websocket.Upgrader  // upgrader 负责 HTTP 到 WebSocket 的协议升级
}

// NewManager 创建 WebSocket 会话管理器。
//
// 参数：
//   - config: 可选配置；不传时使用 DefaultConfig，传入时会通过 normalizeConfig 补齐默认值。
//
// 返回值：
//   - *Manager: 可直接用于创建、获取、关闭和广播会话的管理器实例。
func NewManager(config ...Config) *Manager {
	conf := DefaultConfig()
	if len(config) > 0 {
		conf = normalizeConfig(config[0])
	}

	return &Manager{
		sessions: make(map[string]*Session),
		config:   conf,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  conf.ReadBufferSize,
			WriteBufferSize: conf.WriteBufferSize,
			CheckOrigin:     conf.CheckOrigin,
		},
	}
}

// Create 创建或复用指定编码的 WebSocket 会话。
//
// 参数：
//   - w: HTTP 响应写入器，用于完成 WebSocket 协议升级。
//   - r: HTTP 请求对象，用于读取升级请求和上下文。
//   - code: 连接编码，作为会话唯一键；同一个 code 再次连接会热替换旧连接。
//
// 返回值：
//   - *Session: 创建或复用后的会话。
//   - bool: true 表示新建会话，false 表示复用已有会话并替换连接。
//   - error: code 为空、连接数超限或协议升级失败时返回错误。
func (m *Manager) Create(w http.ResponseWriter, r *http.Request, code string) (*Session, bool, error) {
	reqCtx := xlog.ContextWithTrace(r.Context())
	code = strings.TrimSpace(code)

	if code == "" {
		xlog.Error(reqCtx, "invalid websocket code", map[string]any{"addr": r.RemoteAddr}, ErrInvalidCode)
		return nil, false, ErrInvalidCode
	}

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		xlog.Error(reqCtx, "websocket upgrade failed", map[string]any{"addr": r.RemoteAddr, "code": code}, err)
		return nil, false, err
	}

	m.mu.Lock()
	if oldSession, exists := m.sessions[code]; exists {
		m.mu.Unlock()
		oldSession.replaceConn(conn)
		xlog.Info(reqCtx, "websocket session reconnected", map[string]any{"code": code})
		return oldSession, false, nil
	}

	if m.total.Load() >= m.config.MaxConnTotal {
		m.mu.Unlock()
		_ = conn.Close()
		xlog.Warn(reqCtx, "websocket max connection exceeded", map[string]any{"code": code, "total": m.total.Load()}, ErrMaxConnExceeded)
		return nil, false, ErrMaxConnExceeded
	}

	session := createSession(conn, code, m.config, m.remove)
	m.sessions[code] = session
	total := m.total.Add(1)
	m.mu.Unlock()

	xlog.Info(reqCtx, "websocket session created", map[string]any{"code": code, "total": total})
	return session, true, nil
}

// Get 获取指定编码的会话。
//
// 参数：
//   - code: 连接编码。
//
// 返回值：
//   - *Session: 命中的会话；不存在时为 nil。
//   - bool: true 表示会话存在，false 表示会话不存在。
func (m *Manager) Get(code string) (*Session, bool) {
	m.mu.RLock()
	session, ok := m.sessions[code]
	m.mu.RUnlock()
	return session, ok
}

// CloseConn 根据连接编码关闭会话。
//
// 参数：
//   - code: 连接编码。
//
// 返回值：
//   - error: 会话不存在时返回 ErrSessionNotFound；关闭动作本身不返回底层连接错误。
func (m *Manager) CloseConn(code string) error {
	session, ok := m.Get(code)
	if !ok {
		return ErrSessionNotFound
	}

	session.Close()
	return nil
}

// Count 返回当前会话数量。
//
// 返回值：
//   - int32: 当前仍由管理器持有的在线会话数量。
func (m *Manager) Count() int32 {
	return m.total.Load()
}

// Range 遍历当前会话。
//
// 参数：
//   - fn: 遍历回调；返回 false 时提前停止遍历。
//
// Range 会先复制会话快照再执行回调，避免回调内操作管理器造成锁重入。
func (m *Manager) Range(fn func(session *Session) bool) {
	m.mu.RLock()
	sessions := make([]*Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	m.mu.RUnlock()

	for _, session := range sessions {
		if !fn(session) {
			return
		}
	}
}

// Broadcast 向全部在线会话发送消息。
//
// 参数：
//   - msg: 下行消息内容；会写入每个会话的发送通道。
func (m *Manager) Broadcast(msg []byte) {
	m.Range(func(session *Session) bool {
		session.Write(msg)
		return true
	})
}

// remove 从管理器中移除会话
func (m *Manager) remove(session *Session) {
	m.mu.Lock()
	current := m.sessions[session.Code()]
	if current == session {
		delete(m.sessions, session.Code())
		m.total.Add(-1)
	}
	m.mu.Unlock()
}

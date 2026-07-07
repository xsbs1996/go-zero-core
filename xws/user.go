package xws

import (
	"net/http"
	"strconv"
	"strings"
)

const userCodePrefix = "user:"

// UserCode 将用户 ID 转换为 WebSocket 会话编码。
//
// 参数：
//   - userID: 用户 ID，必须大于 0。
//
// 返回值：
//   - string: user:<id> 格式的连接编码；userID 小于等于 0 时返回空字符串。
//
// 用户会话编码使用 user:<id> 格式，避免和设备、房间等其它业务 code 冲突。
// userID 小于等于 0 时返回空字符串，调用 CreateUser 会得到 ErrInvalidCode。
func UserCode(userID int64) string {
	if userID <= 0 {
		return ""
	}
	return userCodePrefix + strconv.FormatInt(userID, 10)
}

// UserID 从 WebSocket 会话中解析用户 ID。
//
// 参数：
//   - session: WebSocket 会话。
//
// 返回值：
//   - uint64: 解析成功的用户 ID；会话为空、code 为空或格式非法时返回 0。
//
// 优先解析 UserCode 生成的 user:<id> 格式；为了兼容业务直接使用数字字符串作为 code 的场景，
// 当 code 不带 user: 前缀时，也会尝试按无符号整数解析。
func UserID(session *Session) uint64 {
	if session == nil {
		return 0
	}

	code := strings.TrimSpace(session.Code())
	if code == "" {
		return 0
	}
	code = strings.TrimPrefix(code, userCodePrefix)

	userID, err := strconv.ParseUint(code, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}

// CreateUser 根据用户 ID 创建或复用 WebSocket 会话。
//
// 参数：
//   - w: HTTP 响应写入器，用于完成 WebSocket 协议升级。
//   - r: HTTP 请求对象，用于读取升级请求和上下文。
//   - userID: 用户 ID，必须大于 0。
//
// 返回值：
//   - *Session: 创建或复用后的会话。
//   - bool: true 表示新建会话，false 表示复用已有用户会话并替换连接。
//   - error: userID 非法、连接数超限或协议升级失败时返回错误。
//
// 这是 Create 的用户 ID 便捷封装，适合一人一连接或按用户热重连的场景。
func (m *Manager) CreateUser(w http.ResponseWriter, r *http.Request, userID int64) (*Session, bool, error) {
	return m.Create(w, r, UserCode(userID))
}

// GetUser 根据用户 ID 获取 WebSocket 会话。
//
// 参数：
//   - userID: 用户 ID。
//
// 返回值：
//   - *Session: 命中的用户会话；不存在时为 nil。
//   - bool: true 表示会话存在，false 表示会话不存在。
func (m *Manager) GetUser(userID int64) (*Session, bool) {
	return m.Get(UserCode(userID))
}

// CloseUserConn 根据用户 ID 关闭 WebSocket 会话。
//
// 参数：
//   - userID: 用户 ID。
//
// 返回值：
//   - error: 会话不存在时返回 ErrSessionNotFound。
func (m *Manager) CloseUserConn(userID int64) error {
	return m.CloseConn(UserCode(userID))
}

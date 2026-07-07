package xticket

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const ticketParts = 2

// Config 表示签名票据配置。
type Config struct {
	Secret []byte        `json:"secret" yaml:"secret"`          // Secret 表示 HMAC-SHA256 签名密钥。
	Issuer string        `json:"issuer,optional" yaml:"issuer"` // Issuer 表示签发者，可为空。
	TTL    time.Duration `json:"ttl,optional" yaml:"ttl"`       // TTL 表示票据有效期，配置值示例：10m；小于等于 0 时不写入过期时间。
}

// Claims 表示票据载荷。
//
// Claims 会被 JSON 序列化后参与签名。业务层如果需要实现一次性票据、
// 防重放或撤销能力，可以保存 Nonce 或完整票据字符串，并在 Verify 后自行比对和删除。
type Claims[T any] struct {
	Payload   T      `json:"payload"`          // Payload 表示业务载荷。
	Issuer    string `json:"issuer,omitempty"` // Issuer 表示签发者。
	IssuedAt  int64  `json:"iat"`              // IssuedAt 表示签发 Unix 毫秒时间。
	ExpiresAt int64  `json:"exp,omitempty"`    // ExpiresAt 表示过期 Unix 毫秒时间。
	Nonce     string `json:"nonce"`            // Nonce 表示随机值，业务层可用于防重放存储。
}

// Generate 生成签名票据。
//
// 参数：
//   - conf: 票据签名配置，必须包含 Secret。
//   - payload: 业务载荷，会写入 Claims.Payload。
//
// 返回值：
//   - string: 签名后的票据字符串。
//   - error: 密钥为空或 JSON 序列化失败时返回错误。
//
// 票据格式为 base64url(JSON claims) + "." + base64url(HMAC-SHA256 signature)。
// Generate 只负责生成随机 Nonce、写入签发/过期时间并完成签名，不负责服务端存储。
func Generate[T any](conf Config, payload T) (string, error) {
	if len(conf.Secret) == 0 {
		return "", ErrMissingSecret
	}

	now := time.Now()
	claims := Claims[T]{
		Payload:  payload,
		Issuer:   conf.Issuer,
		IssuedAt: now.UnixMilli(),
		Nonce:    randomNonce(),
	}
	if conf.TTL > 0 {
		claims.ExpiresAt = now.Add(conf.TTL).UnixMilli()
	}

	body, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("xticket: marshal claims: %w", err)
	}

	bodyText := base64.RawURLEncoding.EncodeToString(body)
	signText := sign(conf.Secret, bodyText)
	return bodyText + "." + signText, nil
}

// Parse 解析票据并校验签名，不校验过期时间。
//
// 参数：
//   - conf: 票据签名配置，必须包含 Secret。
//   - value: 待解析的票据字符串。
//
// 返回值：
//   - *Claims[T]: 解析成功后的票据载荷。
//   - error: 密钥为空、票据格式非法、签名非法或 JSON 解码失败时返回错误。
//
// Parse 适合需要读取过期票据内容的场景；如果需要同时校验过期时间，应使用 Verify。
func Parse[T any](conf Config, value string) (*Claims[T], error) {
	if len(conf.Secret) == 0 {
		return nil, ErrMissingSecret
	}

	bodyText, signText, err := splitTicket(value)
	if err != nil {
		return nil, err
	}
	if !verifySign(conf.Secret, bodyText, signText) {
		return nil, ErrInvalidTicket
	}

	body, err := base64.RawURLEncoding.DecodeString(bodyText)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTicket, err)
	}

	var claims Claims[T]
	if err := json.Unmarshal(body, &claims); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTicket, err)
	}
	if claims.Nonce == "" || claims.IssuedAt == 0 {
		return nil, ErrInvalidTicket
	}
	return &claims, nil
}

// Verify 解析票据、校验签名并校验过期时间。
//
// 参数：
//   - conf: 票据签名配置，必须包含 Secret。
//   - value: 待校验的票据字符串。
//
// 返回值：
//   - *Claims[T]: 校验成功后的票据载荷。
//   - error: 解析失败、签名非法或票据过期时返回错误。
//
// Verify 不处理服务端存储、防重放或删除逻辑。业务层可在 Verify 成功后根据
// Claims.Nonce 或原始票据字符串自行完成 Redis/DB 校验和删除。
func Verify[T any](conf Config, value string) (*Claims[T], error) {
	claims, err := Parse[T](conf, value)
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt > 0 && time.Now().UnixMilli() > claims.ExpiresAt {
		return nil, ErrTicketExpired
	}
	return claims, nil
}

// splitTicket 拆分票据正文和签名。
func splitTicket(value string) (string, string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != ticketParts || parts[0] == "" || parts[1] == "" {
		return "", "", ErrInvalidTicket
	}
	return parts[0], parts[1], nil
}

// sign 使用 HMAC-SHA256 对票据正文签名。
func sign(secret []byte, body string) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(body))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// verifySign 使用常量时间比较校验签名，避免普通字符串比较带来的时序风险。
func verifySign(secret []byte, body string, signText string) bool {
	expected := sign(secret, body)
	return hmac.Equal([]byte(expected), []byte(signText))
}

// randomNonce 生成票据随机值。
//
// 正常情况下使用 crypto/rand；如果系统随机源失败，则退化为纳秒时间戳，保证仍能生成票据。
func randomNonce() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(b[:])
}

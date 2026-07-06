package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config JWT 配置。
type Config struct {
	Secret        string        `json:"secret" yaml:"secret"`                        // Secret JWT 签名密钥。
	Issuer        string        `json:"issuer,optional" yaml:"issuer"`               // Issuer 签发者。
	Subject       string        `json:"subject,optional" yaml:"subject"`             // Subject 主题。
	Audience      []string      `json:"audience,optional" yaml:"audience"`           // Audience 受众。
	Expire        time.Duration `json:"expire,optional" yaml:"expire"`               // Expire token 有效期，配置值示例：10m。
	RefreshExpire time.Duration `json:"refreshExpire,optional" yaml:"refreshExpire"` // RefreshExpire 刷新有效期，配置值示例：24h。
}

// Claims JWT 载荷。
type Claims struct {
	Data map[string]any `json:"data"` // Data 业务数据。
	jwt.RegisteredClaims
}

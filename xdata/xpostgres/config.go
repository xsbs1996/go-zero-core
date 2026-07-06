package xpostgres

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm/logger"
)

const (
	defaultSSLMode      = "disable"        // 默认 SSL 模式。
	defaultDialTimeout  = 10 * time.Second // 默认建立连接超时时间。
	defaultReadTimeout  = 30 * time.Second // 默认读取超时时间。
	defaultWriteTimeout = 30 * time.Second // 默认写入超时时间。
)

var (
	ErrMissingHost   = errors.New("xpostgres: missing postgres host")     // ErrMissingHost 表示 PostgreSQL 地址为空。
	ErrMissingUser   = errors.New("xpostgres: missing postgres user")     // ErrMissingUser 表示 PostgreSQL 用户名为空。
	ErrMissingDBName = errors.New("xpostgres: missing postgres database") // ErrMissingDBName 表示 PostgreSQL 数据库名为空。
)

// Config PostgreSQL 的 GORM 连接配置。
type Config struct {
	Host                   string        `json:"host" yaml:"host"`                                              // Host PostgreSQL 地址。
	Port                   int           `json:"port,optional" yaml:"port"`                                     // Port PostgreSQL 端口。
	User                   string        `json:"user,optional" yaml:"user"`                                     // User PostgreSQL 用户名。
	Username               string        `json:"username,optional" yaml:"username"`                             // Username PostgreSQL 用户名，兼容旧字段，优先使用 User。
	Password               string        `json:"password,optional" yaml:"password"`                             // Password PostgreSQL 密码。
	DBName                 string        `json:"dbName,optional" yaml:"dbName"`                                 // DBName PostgreSQL 数据库名。
	Database               string        `json:"database,optional" yaml:"database"`                             // Database PostgreSQL 数据库名，兼容旧字段，优先使用 DBName。
	SSLMode                string        `json:"sslMode,default=disable" yaml:"sslMode"`                        // SSLMode SSL 模式，例如 disable、require、verify-full。
	TimeZone               string        `json:"timeZone,optional" yaml:"timeZone"`                             // TimeZone 连接时区，例如 Asia/Shanghai。
	ConnectTimeout         time.Duration `json:"connectTimeout,default=10s" yaml:"connectTimeout"`              // ConnectTimeout 建立连接超时时间，配置值示例：10s。
	ReadTimeout            time.Duration `json:"readTimeout,default=30s" yaml:"readTimeout"`                    // ReadTimeout 读取超时时间，配置值示例：30s。
	WriteTimeout           time.Duration `json:"writeTimeout,default=30s" yaml:"writeTimeout"`                  // WriteTimeout 写入超时时间，配置值示例：30s。
	MaxIdleConns           int           `json:"maxIdleConns,optional" yaml:"maxIdleConns"`                     // MaxIdleConns 最大空闲连接数。
	MaxOpenConns           int           `json:"maxOpenConns,optional" yaml:"maxOpenConns"`                     // MaxOpenConns 最大打开连接数。
	ConnMaxLifetime        int           `json:"connMaxLifetime,optional" yaml:"connMaxLifetime"`               // ConnMaxLifetime 连接最大生命周期，单位分钟。
	ConnMaxIdleTime        int           `json:"connMaxIdleTime,optional" yaml:"connMaxIdleTime"`               // ConnMaxIdleTime 连接最大空闲时间，单位分钟。
	SkipDefaultTransaction bool          `json:"skipDefaultTransaction,optional" yaml:"skipDefaultTransaction"` // SkipDefaultTransaction 是否跳过 GORM 默认事务。
	PrepareStmt            bool          `json:"prepareStmt,optional" yaml:"prepareStmt"`                       // PrepareStmt 是否启用 GORM 预编译语句缓存。
	LogLevel               string        `json:"logLevel,optional" yaml:"logLevel"`                             // LogLevel GORM 日志级别，支持 silent、error、warn、info。
	Sharding               Sharding      `json:"sharding,optional" yaml:"sharding"`                             // Sharding 自动分表配置。
}

// WithDefault 返回补齐默认值后的配置。
func (c Config) WithDefault() Config {
	if c.SSLMode == "" {
		c.SSLMode = defaultSSLMode
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = defaultDialTimeout
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	return c
}

// Validate 校验 PostgreSQL 连接配置。
func (c Config) Validate() error {
	if c.Host == "" {
		return ErrMissingHost
	}
	if c.user() == "" {
		return ErrMissingUser
	}
	if c.dbName() == "" {
		return ErrMissingDBName
	}
	if err := c.Sharding.Validate(); err != nil {
		return err
	}
	return nil
}

// DSN 根据配置生成 PostgreSQL 连接字符串。
func (c Config) DSN() string {
	c = c.WithDefault()

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.user(), c.Password),
		Host:   c.addr(),
		Path:   "/" + c.dbName(),
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	query.Set("connect_timeout", strconv.Itoa(int(c.ConnectTimeout.Seconds())))
	if c.TimeZone != "" {
		query.Set("TimeZone", c.TimeZone)
	}
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

// GormLogLevel 将字符串日志级别转换为 GORM 日志级别。
func (c Config) GormLogLevel() logger.LogLevel {
	switch c.LogLevel {
	case "silent", "Silent", "SILENT":
		return logger.Silent
	case "error", "Error", "ERROR":
		return logger.Error
	case "warn", "Warn", "WARN", "warning", "Warning", "WARNING":
		return logger.Warn
	case "info", "Info", "INFO":
		return logger.Info
	default:
		return logger.Warn
	}
}

// GormLogConfig 将配置转换为 GORM 日志配置。
func (c Config) GormLogConfig() logger.Config {
	return logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  c.GormLogLevel(),
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
	}
}

func (c Config) addr() string {
	if c.Port == 0 {
		return c.Host
	}
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c Config) user() string {
	if c.User != "" {
		return c.User
	}
	return c.Username
}

func (c Config) dbName() string {
	if c.DBName != "" {
		return c.DBName
	}
	return c.Database
}

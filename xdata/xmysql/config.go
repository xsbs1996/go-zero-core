package xmysql

import (
	"errors"
	"strconv"
	"time"

	drivermysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm/logger"
)

const (
	defaultCharset      = "utf8mb4"
	defaultLoc          = "Local"
	defaultTimeout      = 10 * time.Second
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

var (
	ErrMissingHost   = errors.New("xmysql: missing mysql host")     // ErrMissingHost 表示 MySQL 地址为空。
	ErrMissingUser   = errors.New("xmysql: missing mysql user")     // ErrMissingUser 表示 MySQL 用户名为空。
	ErrMissingDBName = errors.New("xmysql: missing mysql database") // ErrMissingDBName 表示 MySQL 数据库名为空。
)

// Config MySQL 的 GORM 连接配置。
type Config struct {
	Host                   string        `json:"host" yaml:"host"`                                              // Host MySQL 地址。
	Port                   int           `json:"port,optional" yaml:"port"`                                     // Port MySQL 端口。
	User                   string        `json:"user,optional" yaml:"user"`                                     // User MySQL 用户名。
	Username               string        `json:"username,optional" yaml:"username"`                             // Username MySQL 用户名，兼容旧字段，优先使用 User。
	Password               string        `json:"password,optional" yaml:"password"`                             // Password MySQL 密码。
	DBName                 string        `json:"dbName,optional" yaml:"dbName"`                                 // DBName MySQL 数据库名。
	Database               string        `json:"database,optional" yaml:"database"`                             // Database MySQL 数据库名，兼容旧字段，优先使用 DBName。
	Charset                string        `json:"charset,default=utf8mb4" yaml:"charset"`                        // Charset 连接字符集。
	ParseTime              bool          `json:"parseTime,optional" yaml:"parseTime"`                           // ParseTime 是否将 DATE/DATETIME 自动解析为 time.Time。
	Loc                    string        `json:"loc,default=Local" yaml:"loc"`                                  // Loc 时间时区，例如 Local、Asia/Shanghai。
	Timeout                time.Duration `json:"timeout,default=10s" yaml:"timeout"`                            // Timeout 建立连接超时时间。
	ReadTimeout            time.Duration `json:"readTimeout,default=30s" yaml:"readTimeout"`                    // ReadTimeout 读取超时时间。
	WriteTimeout           time.Duration `json:"writeTimeout,default=30s" yaml:"writeTimeout"`                  // WriteTimeout 写入超时时间。
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
	if c.Charset == "" {
		c.Charset = defaultCharset
	}
	if c.Loc == "" {
		c.Loc = defaultLoc
	}
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	return c
}

// Validate 校验 MySQL 连接配置。
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

// DSN 根据配置生成 MySQL 连接字符串。
func (c Config) DSN() string {
	c = c.WithDefault()

	dsn := drivermysql.NewConfig()
	dsn.User = c.user()
	dsn.Passwd = c.Password
	dsn.Net = "tcp"
	dsn.Addr = c.addr()
	dsn.DBName = c.dbName()
	dsn.ParseTime = c.ParseTime
	dsn.Loc = c.location()
	dsn.Timeout = c.Timeout
	dsn.ReadTimeout = c.ReadTimeout
	dsn.WriteTimeout = c.WriteTimeout
	dsn.Params = map[string]string{
		"charset": c.Charset,
	}

	return dsn.FormatDSN()
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
	c = c.WithDefault()
	if c.Port == 0 {
		return c.Host
	}
	return c.Host + ":" + strconv.Itoa(c.Port)
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

func (c Config) location() *time.Location {
	loc, err := time.LoadLocation(c.WithDefault().Loc)
	if err != nil {
		return time.Local
	}
	return loc
}

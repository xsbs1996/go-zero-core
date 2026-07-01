package xlog

import (
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	// LogConf 日志配置，直接复用 go-zero logx.LogConf。
	LogConf = logx.LogConf
	// LogField 日志字段，直接复用 go-zero logx.LogField。
	LogField = logx.LogField
	// Writer 日志写入器，直接复用 go-zero logx.Writer。
	Writer = logx.Writer
)

const (
	// ModeConsole 控制台输出模式。
	ModeConsole = "console"
	// ModeFile 文件输出模式。
	ModeFile = "file"
	// ModeVolume 容器挂载卷输出模式。
	ModeVolume = "volume"
	// EncodingJSON JSON 编码格式。
	EncodingJSON = "json"
	// EncodingPlain 普通文本编码格式。
	EncodingPlain = "plain"
	// LevelDebug debug 日志级别。
	LevelDebug = "debug"
	// LevelInfo info 日志级别。
	LevelInfo = "info"
	// LevelError error 日志级别。
	LevelError = "error"
	// LevelSevere severe 日志级别。
	LevelSevere = "severe"
	// RotationDaily 按天切割日志。
	RotationDaily = "daily"
	// RotationSize 按大小切割日志。
	RotationSize = "size"
)

// Config 日志配置，字段语义与 go-zero logx.LogConf 保持一致。
type Config struct {
	ServiceName         string    `json:"serviceName,optional" yaml:"serviceName"`                           // ServiceName 服务名称。
	Mode                string    `json:"mode,default=console,options=[console,file,volume]" yaml:"mode"`    // Mode 日志输出模式，支持 console、file、volume。
	Encoding            string    `json:"encoding,default=json,options=[json,plain]" yaml:"encoding"`        // Encoding 日志编码格式，支持 json、plain。
	TimeFormat          string    `json:"timeFormat,optional" yaml:"timeFormat"`                             // TimeFormat 日志时间格式。
	Path                string    `json:"path,default=logs" yaml:"path"`                                     // Path 日志文件目录。
	Level               string    `json:"level,default=info,options=[debug,info,error,severe]" yaml:"level"` // Level 日志级别，支持 debug、info、error、severe。
	MaxContentLength    uint32    `json:"maxContentLength,optional" yaml:"maxContentLength"`                 // MaxContentLength 单条日志最大内容长度，0 表示不限制。
	Compress            bool      `json:"compress,optional" yaml:"compress"`                                 // Compress 是否压缩日志文件。
	Stat                *bool     `json:"stat,default=true" yaml:"stat"`                                     // Stat 是否输出统计日志，不配置时默认开启。
	KeepDays            int       `json:"keepDays,optional" yaml:"keepDays"`                                 // KeepDays 日志保留天数，0 表示不按天数清理。
	StackCooldownMillis int       `json:"stackCooldownMillis,default=100" yaml:"stackCooldownMillis"`        // StackCooldownMillis 堆栈日志冷却时间，单位毫秒。
	MaxBackups          int       `json:"maxBackups,default=0" yaml:"maxBackups"`                            // MaxBackups 最大备份文件数量，0 表示不限制。
	MaxSize             int       `json:"maxSize,default=0" yaml:"maxSize"`                                  // MaxSize 单个日志文件最大体积，单位 MB，0 表示不限制。
	Rotation            string    `json:"rotation,default=daily,options=[daily,size]" yaml:"rotation"`       // Rotation 日志切割规则，支持 daily、size。
	FileTimeFormat      string    `json:"fileTimeFormat,optional" yaml:"fileTimeFormat"`                     // FileTimeFormat 日志文件时间格式。
	FieldKeys           FieldKeys `json:"fieldKeys,optional" yaml:"fieldKeys"`                               // FieldKeys 日志字段名称配置。
}

// FieldKeys 日志字段名称配置，字段语义与 go-zero logx.FieldKeys 保持一致。
type FieldKeys struct {
	CallerKey    string `json:"callerKey,default=caller" yaml:"callerKey"`           // CallerKey 调用位置字段名。
	ContentKey   string `json:"contentKey,default=content" yaml:"contentKey"`        // ContentKey 日志内容字段名。
	DurationKey  string `json:"durationKey,default=duration" yaml:"durationKey"`     // DurationKey 耗时字段名。
	LevelKey     string `json:"levelKey,default=level" yaml:"levelKey"`              // LevelKey 日志级别字段名。
	SpanKey      string `json:"spanKey,default=span" yaml:"spanKey"`                 // SpanKey 链路跨度字段名。
	TimestampKey string `json:"timestampKey,default=@timestamp" yaml:"timestampKey"` // TimestampKey 时间字段名。
	TraceKey     string `json:"traceKey,default=trace" yaml:"traceKey"`              // TraceKey 链路追踪字段名。
	TruncatedKey string `json:"truncatedKey,default=truncated" yaml:"truncatedKey"`  // TruncatedKey 内容截断标记字段名。
}

type logConfAdapter struct {
	ServiceName         string
	Mode                string
	Encoding            string
	TimeFormat          string
	Path                string
	Level               string
	MaxContentLength    uint32
	Compress            bool
	Stat                bool
	KeepDays            int
	StackCooldownMillis int
	MaxBackups          int
	MaxSize             int
	Rotation            string
	FileTimeFormat      string
	FieldKeys           FieldKeys
}

// WithDefault 返回填充默认值后的日志配置。
func (c Config) WithDefault() Config {
	if c.Mode == "" {
		c.Mode = ModeConsole
	}
	if c.Encoding == "" {
		c.Encoding = EncodingJSON
	}
	if c.Path == "" {
		c.Path = "logs"
	}
	if c.Level == "" {
		c.Level = LevelInfo
	}
	if c.Stat == nil {
		enabled := true
		c.Stat = &enabled
	}
	if c.StackCooldownMillis <= 0 {
		c.StackCooldownMillis = 100
	}
	if c.Rotation == "" {
		c.Rotation = RotationDaily
	}

	return c
}

// ToLogConf 转换为 go-zero logx.LogConf。
func (c Config) ToLogConf() (LogConf, error) {
	c = c.WithDefault()

	adapter := logConfAdapter{
		ServiceName:         c.ServiceName,
		Mode:                c.Mode,
		Encoding:            c.Encoding,
		TimeFormat:          c.TimeFormat,
		Path:                c.Path,
		Level:               c.Level,
		MaxContentLength:    c.MaxContentLength,
		Compress:            c.Compress,
		Stat:                *c.Stat,
		KeepDays:            c.KeepDays,
		StackCooldownMillis: c.StackCooldownMillis,
		MaxBackups:          c.MaxBackups,
		MaxSize:             c.MaxSize,
		Rotation:            c.Rotation,
		FileTimeFormat:      c.FileTimeFormat,
		FieldKeys:           c.FieldKeys,
	}

	data, err := json.Marshal(adapter)
	if err != nil {
		return LogConf{}, err
	}

	var conf LogConf
	if err := json.Unmarshal(data, &conf); err != nil {
		return LogConf{}, err
	}

	return conf, nil
}

// SetupConfig 使用日志配置初始化 logx。
func SetupConfig(c Config) error {
	conf, err := c.ToLogConf()
	if err != nil {
		return err
	}

	return Setup(conf)
}

// MustSetupConfig 使用日志配置初始化 logx，失败时按 logx 规则退出或抛出异常。
func MustSetupConfig(c Config) {
	conf, err := c.ToLogConf()
	logx.Must(err)
	MustSetup(conf)
}

// Setup 初始化 logx 日志配置。
func Setup(c LogConf) error {
	return logx.SetUp(c)
}

// MustSetup 初始化 logx 日志配置，失败时按 logx 规则退出或抛出异常。
func MustSetup(c LogConf) {
	logx.MustSetup(c)
}

// Close 关闭 logx 日志写入器。
func Close() error {
	return logx.Close()
}

// SetWriter 设置日志写入器。
func SetWriter(w Writer) {
	logx.SetWriter(w)
}

// AddWriter 增加日志写入器。
func AddWriter(w Writer) {
	logx.AddWriter(w)
}

// Field 创建日志字段。
func Field(key string, value any) LogField {
	return logx.Field(key, value)
}

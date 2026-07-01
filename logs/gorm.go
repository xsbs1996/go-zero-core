package logs

import (
	"context"
	"errors"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

var _ gormlogger.Interface = (*GormLogger)(nil)

// GormLogger GORM 日志适配器。
type GormLogger struct {
	Config gormlogger.Config // Config GORM 日志配置。
}

// NewGormLogger 创建 GORM 日志适配器。
func NewGormLogger(config ...gormlogger.Config) gormlogger.Interface {
	conf := gormlogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	}
	if len(config) > 0 {
		conf = config[0]
	}

	return &GormLogger{Config: conf}
}

// LogMode 设置 GORM 日志级别。
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.Config.LogLevel = level
	return &newLogger
}

// Info 输出 GORM info 日志。
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if IsDisabled(ctx) {
		return
	}
	if l.Config.LogLevel >= gormlogger.Info {
		Info(ctx, formatMessage(msg, data...), nil)
	}
}

// Warn 输出 GORM warn 日志。
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if IsDisabled(ctx) {
		return
	}
	if l.Config.LogLevel >= gormlogger.Warn {
		Warn(ctx, formatMessage(msg, data...), nil, nil)
	}
}

// Error 输出 GORM error 日志。
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if IsDisabled(ctx) {
		return
	}
	if l.Config.LogLevel >= gormlogger.Error {
		Error(ctx, formatMessage(msg, data...), nil, nil)
	}
}

// Trace 输出 GORM SQL 日志。
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if IsDisabled(ctx) {
		return
	}
	if l.Config.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		Error(ctx, "gorm sql error", gormBody(sql, rows, elapsed), err)
	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold > 0 && l.Config.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		Warn(ctx, fmt.Sprintf("slow sql >= %v", l.Config.SlowThreshold), gormBody(sql, rows, elapsed), nil)
	case l.Config.LogLevel == gormlogger.Info:
		sql, rows := fc()
		Info(ctx, "gorm sql", gormBody(sql, rows, elapsed))
	}
}

// ParamsFilter 过滤 GORM SQL 参数。
func (l *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}

	return sql, params
}

func formatMessage(msg string, data ...interface{}) string {
	if len(data) == 0 {
		return msg
	}

	return fmt.Sprintf(msg, data...)
}

func gormBody(sql string, rows int64, elapsed time.Duration) map[string]any {
	return map[string]any{
		"sql":        sql,
		"rows":       rows,
		"elapsed_ms": float64(elapsed.Nanoseconds()) / 1e6,
	}
}

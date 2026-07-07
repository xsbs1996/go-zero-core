package xtask

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Option 表示定时任务管理器配置函数。
type Option func(options *options)

type options struct {
	withSeconds bool
	location    *time.Location
}

// WithSeconds 开启秒级 cron 表达式。
//
// 返回值：
//   - Option: 允许使用 6 段秒级 cron 表达式的配置函数。
func WithSeconds() Option {
	return func(options *options) {
		options.withSeconds = true
	}
}

// WithLocation 设置 cron 调度时区。
//
// 参数：
//   - location: cron 调度时区；传 nil 时该配置不生效。
//
// 返回值：
//   - Option: 设置 cron 调度时区的配置函数。
func WithLocation(location *time.Location) Option {
	return func(options *options) {
		if location != nil {
			options.location = location
		}
	}
}

func buildCronOptions(options options) []cron.Option {
	cronOptions := make([]cron.Option, 0, 2)

	if options.withSeconds {
		cronOptions = append(cronOptions, cron.WithSeconds())
	}
	if options.location != nil {
		cronOptions = append(cronOptions, cron.WithLocation(options.location))
	}

	return cronOptions
}

package taskx

import "errors"

var (
	ErrNilManager = errors.New("taskx: nil manager") // ErrNilManager 表示定时任务管理器为空
	ErrEmptyName  = errors.New("taskx: empty name")  // ErrEmptyName 表示任务名称为空
	ErrEmptySpec  = errors.New("taskx: empty spec")  // ErrEmptySpec 表示 cron 表达式为空
	ErrNilFunc    = errors.New("taskx: nil func")    // ErrNilFunc 表示任务执行函数为空
)

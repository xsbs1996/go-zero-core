package xtask

import "errors"

var (
	ErrNilManager = errors.New("xtask: nil manager") // ErrNilManager 表示定时任务管理器为空
	ErrEmptyName  = errors.New("xtask: empty name")  // ErrEmptyName 表示任务名称为空
	ErrEmptySpec  = errors.New("xtask: empty spec")  // ErrEmptySpec 表示 cron 表达式为空
	ErrNilFunc    = errors.New("xtask: nil func")    // ErrNilFunc 表示任务执行函数为空
)

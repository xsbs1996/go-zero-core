package xtask

import (
	"context"
	"sync"

	"go-zero-core/xlog"

	"github.com/robfig/cron/v3"
)

// Job 表示一个定时任务
type Job struct {
	Name   string                    // Name 表示任务名称
	Spec   string                    // Spec 表示 cron 表达式
	RunNow bool                      // RunNow 表示注册后是否立即执行一次
	Run    func(ctx context.Context) // Run 表示任务执行函数
}

// Manager 管理定时任务生命周期
type Manager struct {
	mu      sync.RWMutex            // mu 保护 entries 的并发读写
	cron    *cron.Cron              // cron 表示底层调度器
	entries map[string]cron.EntryID // entries 保存任务名称到任务 ID 的映射
	ctx     context.Context         // ctx 表示任务管理器上下文
	cancel  context.CancelFunc      // cancel 用于停止任务上下文
}

// NewManager 创建定时任务管理器
func NewManager(opts ...Option) *Manager {
	optionValues := options{}
	for _, opt := range opts {
		if opt != nil {
			opt(&optionValues)
		}
	}

	ctx, cancel := context.WithCancel(xlog.ContextWithTrace(context.Background()))

	return &Manager{
		cron:    cron.New(buildCronOptions(optionValues)...),
		entries: make(map[string]cron.EntryID),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动定时任务调度器
func (m *Manager) Start() {
	if m == nil {
		return
	}
	m.cron.Start()
}

// Stop 停止定时任务调度器并等待运行中的任务退出
func (m *Manager) Stop() {
	if m == nil {
		return
	}

	m.cancel()
	ctx := m.cron.Stop()
	<-ctx.Done()
}

// AddFunc 注册定时任务函数
func (m *Manager) AddFunc(name string, spec string, fn func(ctx context.Context)) error {
	return m.Add(Job{
		Name: name,
		Spec: spec,
		Run:  fn,
	})
}

// AddFuncNow 注册定时任务函数并立即执行一次
func (m *Manager) AddFuncNow(name string, spec string, fn func(ctx context.Context)) error {
	return m.Add(Job{
		Name:   name,
		Spec:   spec,
		RunNow: true,
		Run:    fn,
	})
}

// Add 注册定时任务
func (m *Manager) Add(job Job) error {
	if m == nil {
		return ErrNilManager
	}
	if job.Name == "" {
		return ErrEmptyName
	}
	if job.Spec == "" {
		return ErrEmptySpec
	}
	if job.Run == nil {
		return ErrNilFunc
	}

	id, err := m.cron.AddFunc(job.Spec, func() {
		job.Run(m.ctx)
	})
	if err != nil {
		return err
	}

	m.mu.Lock()
	if oldID, exists := m.entries[job.Name]; exists {
		m.cron.Remove(oldID)
	}
	m.entries[job.Name] = id
	m.mu.Unlock()

	if job.RunNow {
		go job.Run(m.ctx)
	}

	return nil
}

// Remove 移除定时任务
func (m *Manager) Remove(name string) bool {
	if m == nil {
		return false
	}

	m.mu.Lock()
	id, ok := m.entries[name]
	if ok {
		delete(m.entries, name)
		m.cron.Remove(id)
	}
	m.mu.Unlock()

	return ok
}

// Count 返回已注册任务数量
func (m *Manager) Count() int {
	if m == nil {
		return 0
	}

	m.mu.RLock()
	count := len(m.entries)
	m.mu.RUnlock()
	return count
}

// Entries 返回底层 cron 任务快照
func (m *Manager) Entries() []cron.Entry {
	if m == nil {
		return nil
	}
	return m.cron.Entries()
}

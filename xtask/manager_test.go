package xtask

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestManagerAddRemoveAndCount 验证任务添加、删除和数量统计。
func TestManagerAddRemoveAndCount(t *testing.T) {
	manager := NewManager(WithSeconds(), WithLocation(time.UTC))
	if err := manager.AddFunc("job", "*/5 * * * * *", func(context.Context) {}); err != nil {
		t.Fatalf("AddFunc() error = %v", err)
	}
	if got := manager.Count(); got != 1 {
		t.Fatalf("Count() = %d", got)
	}
	if !manager.Remove("job") {
		t.Fatal("Remove(job) = false")
	}
	if got := manager.Count(); got != 0 {
		t.Fatalf("Count() after remove = %d", got)
	}
}

// TestManagerAddValidation 验证任务注册时的必填参数校验。
func TestManagerAddValidation(t *testing.T) {
	manager := NewManager()
	if err := manager.Add(Job{}); !errors.Is(err, ErrEmptyName) {
		t.Fatalf("Add(empty) error = %v", err)
	}
	if err := manager.Add(Job{Name: "job"}); !errors.Is(err, ErrEmptySpec) {
		t.Fatalf("Add(missing spec) error = %v", err)
	}
	if err := manager.Add(Job{Name: "job", Spec: "* * * * *"}); !errors.Is(err, ErrNilFunc) {
		t.Fatalf("Add(missing handler) error = %v", err)
	}
}

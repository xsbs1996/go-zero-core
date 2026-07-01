package xpostgres

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var (
	ErrAlreadyInitialized = errors.New("xpostgres: database already initialized") // ErrAlreadyInitialized 表示全局 GORM 单例已经初始化。
	ErrNotInitialized     = errors.New("xpostgres: database not initialized")     // ErrNotInitialized 表示全局 GORM 单例未初始化。
)

var (
	gormDB *gorm.DB
	gormMu sync.RWMutex
)

// Init 根据配置初始化全局 GORM 单例，重复初始化时返回 ErrAlreadyInitialized。
func Init(conf Config, opts ...ConnectOption) error {
	gormMu.Lock()
	defer gormMu.Unlock()

	if gormDB != nil {
		return ErrAlreadyInitialized
	}

	db, err := Connect(conf, opts...)
	if err != nil {
		return err
	}

	gormDB = db
	return nil
}

// MustInit 根据配置初始化全局 GORM 单例，失败时直接 panic。
func MustInit(conf Config, opts ...ConnectOption) {
	if err := Init(conf, opts...); err != nil {
		panic(err)
	}
}

// SetDB 设置全局 GORM 单例，已初始化时不会覆盖已有实例。
func SetDB(db *gorm.DB) {
	if db == nil {
		return
	}

	gormMu.Lock()
	defer gormMu.Unlock()

	if gormDB != nil {
		return
	}

	gormDB = db
}

// GetDB 返回全局 GORM 单例，未初始化时 panic。
func GetDB() *gorm.DB {
	gormMu.RLock()
	defer gormMu.RUnlock()

	if gormDB == nil {
		panic("xpostgres: database not initialized, call xpostgres.Init or xpostgres.SetDB first")
	}
	return gormDB
}

// Close 关闭全局 GORM 单例的底层连接，并清空全局单例。
func Close() error {
	gormMu.Lock()
	defer gormMu.Unlock()

	if gormDB == nil {
		return ErrNotInitialized
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get postgres sql db failed: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close postgres connection failed: %w", err)
	}

	gormDB = nil
	return nil
}

// IsInitialized 返回全局 GORM 单例是否已初始化。
func IsInitialized() bool {
	gormMu.RLock()
	defer gormMu.RUnlock()
	return gormDB != nil
}

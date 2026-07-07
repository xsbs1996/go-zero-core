package xmysql

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var (
	ErrAlreadyInitialized = errors.New("xmysql: database already initialized") // ErrAlreadyInitialized 表示全局 GORM 单例已经初始化。
	ErrNotInitialized     = errors.New("xmysql: database not initialized")     // ErrNotInitialized 表示全局 GORM 单例未初始化。
)

var (
	gormDB *gorm.DB
	gormMu sync.RWMutex
)

// Init 根据配置初始化全局 GORM 单例，重复初始化时返回 ErrAlreadyInitialized。
//
// 参数：
//   - conf: MySQL 连接配置。
//   - opts: 可选连接配置。
//
// 返回值：
//   - error: 初始化成功返回 nil；重复初始化或连接失败时返回错误。
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
//
// 参数：
//   - conf: MySQL 连接配置。
//   - opts: 可选连接配置。
func MustInit(conf Config, opts ...ConnectOption) {
	if err := Init(conf, opts...); err != nil {
		panic(err)
	}
}

// SetDB 设置全局 GORM 单例，已初始化时不会覆盖已有实例。
//
// 参数：
//   - db: 外部创建的 GORM DB。
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
//
// 返回值：
//   - *gorm.DB: 全局 GORM DB。
func GetDB() *gorm.DB {
	gormMu.RLock()
	defer gormMu.RUnlock()

	if gormDB == nil {
		panic("xmysql: database not initialized, call xmysql.Init or xmysql.SetDB first")
	}
	return gormDB
}

// Close 关闭全局 GORM 单例的底层连接，并清空全局单例。
//
// 返回值：
//   - error: 关闭成功或未初始化返回 nil；底层连接关闭失败时返回错误。
func Close() error {
	gormMu.Lock()
	defer gormMu.Unlock()

	if gormDB == nil {
		return ErrNotInitialized
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get mysql sql db failed: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close mysql connection failed: %w", err)
	}

	gormDB = nil
	return nil
}

// IsInitialized 返回全局 GORM 单例是否已初始化。
//
// 返回值：
//   - bool: true 表示全局 GORM DB 已初始化。
func IsInitialized() bool {
	gormMu.RLock()
	defer gormMu.RUnlock()
	return gormDB != nil
}

package xpostgres

import (
	"fmt"
	"time"

	"go-zero-core/xlog"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect 根据配置创建 GORM PostgreSQL 连接。
func Connect(conf Config, opts ...ConnectOption) (*gorm.DB, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	if conf.Sharding.Enabled && conf.PrepareStmt {
		return nil, ErrShardingPrepareStmtUnsupported
	}

	conf = conf.WithDefault()
	options := newOptions(conf, opts...)

	db, err := gorm.Open(gormpostgres.Open(options.dsn), options.gormOptions...)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection failed: %w", err)
	}

	if err := useSharding(db, conf.Sharding, options); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get postgres sql db failed: %w", err)
	}

	if conf.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Minute)
	}
	if conf.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTime) * time.Minute)
	}

	if options.ping {
		if err := sqlDB.Ping(); err != nil {
			return nil, fmt.Errorf("ping postgres failed: %w", err)
		}
	}

	return db, nil
}

// MustConnect 根据配置创建 GORM PostgreSQL 连接，失败时直接 panic。
func MustConnect(conf Config, opts ...ConnectOption) *gorm.DB {
	db, err := Connect(conf, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

func defaultGormConfig(conf Config) *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction: conf.SkipDefaultTransaction,
		PrepareStmt:            conf.PrepareStmt,
		Logger:                 xlog.NewGormLogger(conf.GormLogConfig()),
	}
}

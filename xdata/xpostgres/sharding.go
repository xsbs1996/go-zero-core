package xpostgres

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	gormsharding "gorm.io/sharding"
)

var (
	ErrInvalidPrimaryKeyGenerator     = errors.New("xpostgres: invalid sharding primary key generator") // ErrInvalidPrimaryKeyGenerator 表示分表主键生成策略非法，当前支持 snowflake、custom、none。
	ErrShardingPrepareStmtUnsupported = errors.New("xpostgres: sharding does not support prepare stmt") // ErrShardingPrepareStmtUnsupported 表示分表插件不支持 GORM PrepareStmt。
	ErrMissingShards                  = errors.New("xpostgres: missing sharding shards")                // ErrMissingShards 表示分表数量为空。
	ErrMissingTables                  = errors.New("xpostgres: missing sharding tables")                // ErrMissingTables 表示分表表名映射为空。
	ErrMissingKey                     = errors.New("xpostgres: missing sharding key")                   // ErrMissingKey 表示分表键为空。
)

// Sharding GORM 自动分表配置。
type Sharding struct {
	Enabled bool                     `json:"enabled,optional" yaml:"enabled"` // Enabled 是否启用自动分表，关闭时不注册 sharding 插件。
	Tables  map[string]ShardingTable `json:"tables,optional" yaml:"tables"`   // Tables 表名到分表配置的映射，key 是逻辑表名，例如 orders。
}

// ShardingTable 单表自动分表配置。
type ShardingTable struct {
	ShardingKey         string `json:"shardingKey" yaml:"shardingKey"`                          // ShardingKey 分表键字段名，查询、修改、删除分表数据时需要带等值条件。
	NumberOfShards      uint   `json:"numberOfShards" yaml:"numberOfShards"`                    // NumberOfShards 分表数量，例如 64 会生成 orders_00 到 orders_63。
	PrimaryKeyGenerator string `json:"primaryKeyGenerator,optional" yaml:"primaryKeyGenerator"` // PrimaryKeyGenerator 主键生成策略，snowflake 使用插件雪花算法，custom/none 使用自定义函数或不自动生成。
	DoubleWrite         bool   `json:"doubleWrite,optional" yaml:"doubleWrite"`                 // DoubleWrite 是否同时写入逻辑表和分表，通常只在老表迁移期启用。
}

// Validate 校验自动分表配置。
func (s Sharding) Validate() error {
	if !s.Enabled {
		return nil
	}
	if len(s.Tables) == 0 {
		return ErrMissingTables
	}
	for _, table := range s.Tables {
		if table.ShardingKey == "" {
			return ErrMissingKey
		}
		if table.NumberOfShards == 0 {
			return ErrMissingShards
		}
	}
	return nil
}

// useSharding 根据配置为每张逻辑表注册 GORM sharding 插件。
func useSharding(db *gorm.DB, conf Sharding, opts connectOptions) error {
	if !conf.Enabled {
		return nil
	}

	for tableName, tableConf := range conf.Tables {
		pkGenerator, err := primaryKeyGenerator(tableConf.PrimaryKeyGenerator)
		if err != nil {
			return err
		}

		// 使用 custom/none 且未注入自定义函数时，返回 0 表示不自动生成主键。
		pkGeneratorFn := opts.shardingPrimaryKeyFunc
		if pkGenerator == gormsharding.PKCustom && pkGeneratorFn == nil {
			pkGeneratorFn = func(tableIdx int64) int64 {
				return 0
			}
		}

		if err := db.Use(gormsharding.Register(gormsharding.Config{
			DoubleWrite:           tableConf.DoubleWrite,
			ShardingKey:           tableConf.ShardingKey,
			NumberOfShards:        tableConf.NumberOfShards,
			PrimaryKeyGenerator:   pkGenerator,
			PrimaryKeyGeneratorFn: pkGeneratorFn,
		}, tableName)); err != nil {
			return fmt.Errorf("register postgres sharding failed: %w", err)
		}
	}

	return nil
}

// primaryKeyGenerator 将配置字符串转换为 sharding 插件的主键生成策略。
func primaryKeyGenerator(name string) (int, error) {
	switch name {
	case "", "snowflake", "Snowflake", "SNOWFLAKE":
		return gormsharding.PKSnowflake, nil
	case "none", "None", "NONE", "custom", "Custom", "CUSTOM":
		return gormsharding.PKCustom, nil
	default:
		return 0, ErrInvalidPrimaryKeyGenerator
	}
}

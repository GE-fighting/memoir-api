package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository 基础仓库接口
type Repository interface {
	DB() *gorm.DB
}

// BaseRepository 基础仓库实现
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓库
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB 获取数据库连接
func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}

// WithTx 使用事务执行操作
func (r *BaseRepository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

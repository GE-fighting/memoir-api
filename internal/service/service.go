package service

import (
	"context"

	"memoir-api/internal/repository"

	"gorm.io/gorm"
)

// Service 基础服务接口
type Service interface {
	Repository() repository.Repository
}

// BaseService 基础服务实现
type BaseService struct {
	repo repository.Repository
}

// NewBaseService 创建基础服务
func NewBaseService(repo repository.Repository) *BaseService {
	return &BaseService{
		repo: repo,
	}
}

// Repository 获取仓库实例
func (s *BaseService) Repository() repository.Repository {
	return s.repo
}

// WithTx 在事务中执行操作
func (s *BaseService) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	baseRepo, ok := s.repo.(*repository.BaseRepository)
	if !ok {
		// 如果不是BaseRepository，可能是一个特定仓库的实现
		// 尝试调用其DB()方法获取数据库连接
		return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 创建一个带有事务上下文的新上下文
			txCtx := context.WithValue(ctx, "tx_db", tx)
			// 使用新上下文执行函数
			return fn(txCtx)
		})
	}

	// 使用BaseRepository的WithTx方法
	return baseRepo.WithTx(ctx, func(tx *gorm.DB) error {
		// 创建一个带有事务上下文的新上下文
		txCtx := context.WithValue(ctx, "tx_db", tx)
		// 使用新上下文执行函数
		return fn(txCtx)
	})
}

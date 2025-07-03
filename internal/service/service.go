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

// EmailService 邮件服务接口
type EmailService interface {
	// 发送验证邮件
	SendVerificationEmail(ctx context.Context, toAddress, username, verificationCode string) error

	// 发送密码重置邮件
	SendPasswordResetEmail(ctx context.Context, toAddress, resetToken string) error

	// 发送通知邮件
	SendNotificationEmail(ctx context.Context, toAddress, username, message string) error

	// 发送欢迎邮件
	SendWelcomeEmail(ctx context.Context, toAddress, username string) error

	// 发送纪念日邮件
	SendAnniversaryEmail(ctx context.Context, toAddress, username, partnerName string, days int, date string) error

	// 发送节日邮件
	SendFestivalEmail(ctx context.Context, toAddress, username, partnerName, festivalName string) error

	// 处理邮件队列
	ProcessEmailQueue(ctx context.Context)

	// 存储验证码到Redis
	StoreVerificationCode(ctx context.Context, email, code string) error

	// 验证验证码
	VerifyCode(ctx context.Context, email, code string) (bool, error)

	// 存储密码重置令牌
	StorePasswordResetToken(ctx context.Context, email, token string) error

	// 验证密码重置令牌
	VerifyPasswordResetToken(ctx context.Context, email, token string) (bool, error)
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

// Repository 获取仓库
func (s *BaseService) Repository() repository.Repository {
	return s.repo
}

// WithTx 事务包装
func (s *BaseService) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// 使用数据库连接的事务功能
	return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建一个包含事务的上下文
		txCtx := context.WithValue(ctx, "tx", tx)
		// 执行业务逻辑
		return fn(txCtx)
	})
}

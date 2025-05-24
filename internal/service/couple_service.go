package service

import (
	"context"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

// CoupleService 情侣关系服务接口
type CoupleService interface {
	Service
	CreateCouple(ctx context.Context, couple *models.Couple) (*models.Couple, error)
	GetCoupleByID(ctx context.Context, id int64) (*models.Couple, error)
	UpdateCouple(ctx context.Context, couple *models.Couple) error
	DeleteCouple(ctx context.Context, id int64) error
	ListCouples(ctx context.Context, offset, limit int) ([]*models.Couple, int64, error)
	GetCoupleUsers(ctx context.Context, coupleID int64) ([]*models.User, error)
}

// coupleService 情侣关系服务实现
type coupleService struct {
	*BaseService
	coupleRepo repository.CoupleRepository
	userRepo   repository.UserRepository
}

// NewCoupleService 创建情侣关系服务
func NewCoupleService(
	coupleRepo repository.CoupleRepository,
	userRepo repository.UserRepository,
) CoupleService {
	return &coupleService{
		BaseService: NewBaseService(coupleRepo),
		coupleRepo:  coupleRepo,
		userRepo:    userRepo,
	}
}

// CreateCouple 创建情侣关系
func (s *coupleService) CreateCouple(ctx context.Context, couple *models.Couple) (*models.Couple, error) {

	// 默认开启自动生成视频和提醒
	couple.AutoGenerateVideo = true
	couple.ReminderNotifications = true

	// 创建情侣关系
	if err := s.coupleRepo.Create(ctx, couple); err != nil {
		return nil, err
	}

	return couple, nil
}

// GetCoupleByID 通过ID获取情侣关系
func (s *coupleService) GetCoupleByID(ctx context.Context, id int64) (*models.Couple, error) {
	return s.coupleRepo.GetByID(ctx, id)
}

// UpdateCouple 更新情侣关系
func (s *coupleService) UpdateCouple(ctx context.Context, couple *models.Couple) error {
	// 检查情侣关系是否存在
	_, err := s.coupleRepo.GetByID(ctx, couple.ID)
	if err != nil {
		return err
	}

	return s.coupleRepo.Update(ctx, couple)
}

// DeleteCouple 删除情侣关系
func (s *coupleService) DeleteCouple(ctx context.Context, id int64) error {
	// 获取该情侣关系下的用户
	users, err := s.userRepo.ListByCoupleID(ctx, id)
	if err != nil {
		return err
	}

	// 开启事务
	return s.WithTx(ctx, func(ctx context.Context) error {
		// 删除所有用户
		for _, user := range users {
			if err := s.userRepo.Delete(ctx, user.ID); err != nil {
				return err
			}
		}

		// 删除情侣关系
		return s.coupleRepo.Delete(ctx, id)
	})
}

// ListCouples 获取情侣关系列表
func (s *coupleService) ListCouples(ctx context.Context, offset, limit int) ([]*models.Couple, int64, error) {
	return s.coupleRepo.List(ctx, offset, limit)
}

// GetCoupleUsers 获取情侣关系下的用户
func (s *coupleService) GetCoupleUsers(ctx context.Context, coupleID int64) ([]*models.User, error) {
	// 检查情侣关系是否存在
	_, err := s.coupleRepo.GetByID(ctx, coupleID)
	if err != nil {
		return nil, err
	}

	return s.userRepo.ListByCoupleID(ctx, coupleID)
}

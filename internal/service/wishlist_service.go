package service

import (
	"context"
	"errors"
	"fmt"
	"memoir-api/internal/api/dto"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

var (
	ErrWishlistNotFound = errors.New("心愿不存在")
)

// WishlistService 心愿清单服务接口
type WishlistService interface {
	Service
	CreateWishlist(ctx context.Context, wishlistDTO *dto.CreateWishlistRequest) (*models.Wishlist, error)
	GetWishlistByID(ctx context.Context, id int64) (*models.Wishlist, error)
	ListWishlistsByCoupleID(ctx context.Context, coupleID int64) ([]dto.WishlistDTO, error)
	ListWishlistsByStatus(ctx context.Context, coupleID int64, status string) ([]*models.Wishlist, error)
	ListWishlistsByPriority(ctx context.Context, coupleID int64, priority int) ([]*models.Wishlist, error)
	ListUpcomingReminders(ctx context.Context, daysAhead int) ([]*models.Wishlist, error)
	UpdateWishlist(ctx context.Context, wishlist *models.Wishlist) error
	UpdateWishlistStatus(ctx context.Context, id int64, status string) error
	DeleteWishlist(ctx context.Context, id int64) error
	UpdateWishlistByRequest(ctx context.Context, req *dto.UpdateWishlistRequest) (*models.Wishlist, error)
}

// wishlistService 心愿清单服务实现
type wishlistService struct {
	*BaseService
	wishlistRepo repository.WishlistRepository
}

// NewWishlistService 创建心愿清单服务
func NewWishlistService(wishlistRepo repository.WishlistRepository) WishlistService {
	return &wishlistService{
		BaseService:  NewBaseService(wishlistRepo),
		wishlistRepo: wishlistRepo,
	}
}

// CreateWishlist 创建心愿
func (s *wishlistService) CreateWishlist(ctx context.Context, wishlistDTO *dto.CreateWishlistRequest) (*models.Wishlist, error) {
	wishlist, err := wishlistDTO.ToModel()
	if err != nil {
		return nil, err
	}
	if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("创建心愿失败: %w", err)
	}
	return wishlist, nil
}

// GetWishlistByID 通过ID获取心愿
func (s *wishlistService) GetWishlistByID(ctx context.Context, id int64) (*models.Wishlist, error) {
	wishlist, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWishlistNotFound) {
			return nil, ErrWishlistNotFound
		}
		return nil, fmt.Errorf("获取心愿失败: %w", err)
	}
	return wishlist, nil
}

// ListWishlistsByCoupleID 获取情侣关系下的所有心愿
func (s *wishlistService) ListWishlistsByCoupleID(ctx context.Context, coupleID int64) ([]dto.WishlistDTO, error) {
	entities, err := s.wishlistRepo.ListByCoupleID(ctx, coupleID)
	if err != nil {
		return nil, err
	}
	return dto.WishlistsFromModels(entities), nil
}

// ListWishlistsByStatus 按状态获取心愿
func (s *wishlistService) ListWishlistsByStatus(ctx context.Context, coupleID int64, status string) ([]*models.Wishlist, error) {
	return s.wishlistRepo.ListByStatus(ctx, coupleID, status)
}

// ListWishlistsByPriority 按优先级获取心愿
func (s *wishlistService) ListWishlistsByPriority(ctx context.Context, coupleID int64, priority int) ([]*models.Wishlist, error) {
	return s.wishlistRepo.ListByPriority(ctx, coupleID, priority)
}

// ListUpcomingReminders 获取即将到期的提醒
func (s *wishlistService) ListUpcomingReminders(ctx context.Context, daysAhead int) ([]*models.Wishlist, error) {
	return s.wishlistRepo.ListUpcomingReminders(ctx, daysAhead)
}

// UpdateWishlist 更新心愿
func (s *wishlistService) UpdateWishlist(ctx context.Context, wishlist *models.Wishlist) error {
	// 检查心愿是否存在
	_, err := s.wishlistRepo.GetByID(ctx, wishlist.ID)
	if err != nil {
		if errors.Is(err, repository.ErrWishlistNotFound) {
			return ErrWishlistNotFound
		}
		return fmt.Errorf("查询心愿失败: %w", err)
	}

	// 更新心愿
	if err := s.wishlistRepo.Update(ctx, wishlist); err != nil {
		return fmt.Errorf("更新心愿失败: %w", err)
	}
	return nil
}

// UpdateWishlistByRequest 根据请求更新心愿
func (s *wishlistService) UpdateWishlistByRequest(ctx context.Context, req *dto.UpdateWishlistRequest) (*models.Wishlist, error) {
	// 获取现有的心愿项
	existingWishlist, err := s.GetWishlistByID(ctx, req.ID)
	if err != nil {
		return nil, err // GetWishlistByID 已经处理了错误包装
	}

	// 将请求中的更新应用到现有心愿项
	if err := req.ApplyToModel(existingWishlist); err != nil {
		return nil, fmt.Errorf("更新心愿参数无效: %w", err)
	}

	// 调用更新方法
	if err := s.UpdateWishlist(ctx, existingWishlist); err != nil {
		return nil, err // UpdateWishlist 已经处理了错误包装
	}

	return existingWishlist, nil
}

// UpdateWishlistStatus 更新心愿状态
func (s *wishlistService) UpdateWishlistStatus(ctx context.Context, id int64, status string) error {
	// 检查心愿是否存在
	_, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWishlistNotFound) {
			return ErrWishlistNotFound
		}
		return fmt.Errorf("查询心愿失败: %w", err)
	}

	// 更新心愿状态
	if err := s.wishlistRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("更新心愿状态失败: %w", err)
	}
	return nil
}

// DeleteWishlist 删除心愿
func (s *wishlistService) DeleteWishlist(ctx context.Context, id int64) error {
	// 检查心愿是否存在
	_, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWishlistNotFound) {
			return ErrWishlistNotFound
		}
		return fmt.Errorf("查询心愿失败: %w", err)
	}

	// 删除心愿
	if err := s.wishlistRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除心愿失败: %w", err)
	}
	return nil
}

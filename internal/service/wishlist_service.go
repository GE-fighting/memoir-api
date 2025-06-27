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

	// 附件关联管理
	AssociateAttachments(ctx context.Context, wishlistID int64, attachmentIDs []int64) error
	RemoveAttachment(ctx context.Context, wishlistID int64, attachmentID int64) error
	GetAttachments(ctx context.Context, wishlistID int64) ([]models.Attachment, error)
}

// wishlistService 心愿清单服务实现
type wishlistService struct {
	*BaseService
	wishlistRepo           repository.WishlistRepository
	wishlistAttachmentRepo repository.WishlistAttachmentRepository
	attachmentRepo         repository.AttachmentRepository
}

// NewWishlistService 创建心愿清单服务
func NewWishlistService(
	wishlistRepo repository.WishlistRepository,
	wishlistAttachmentRepo repository.WishlistAttachmentRepository,
	attachmentRepo repository.AttachmentRepository,
) WishlistService {
	return &wishlistService{
		BaseService:            NewBaseService(wishlistRepo),
		wishlistRepo:           wishlistRepo,
		wishlistAttachmentRepo: wishlistAttachmentRepo,
		attachmentRepo:         attachmentRepo,
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

	// 处理附件关联
	if wishlistDTO.AttachmentIDs != nil && len(wishlistDTO.AttachmentIDs) > 0 {
		if err := s.AssociateAttachments(ctx, wishlist.ID, wishlistDTO.AttachmentIDs); err != nil {
			return nil, fmt.Errorf("关联附件失败: %w", err)
		}
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

	// 加载附件
	attachments, err := s.GetAttachments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("加载附件失败: %w", err)
	}
	wishlist.Attachments = attachments

	return wishlist, nil
}

// ListWishlistsByCoupleID 获取情侣关系下的所有心愿
func (s *wishlistService) ListWishlistsByCoupleID(ctx context.Context, coupleID int64) ([]dto.WishlistDTO, error) {
	entities, err := s.wishlistRepo.ListByCoupleID(ctx, coupleID)
	if err != nil {
		return nil, err
	}

	// 为每个心愿加载附件
	result := make([]dto.WishlistDTO, len(entities))
	for i, entity := range entities {
		// 加载附件
		attachments, err := s.GetAttachments(ctx, entity.ID)
		if err == nil && len(attachments) > 0 {
			entity.Attachments = attachments
		}

		// 转换为DTO
		result[i] = dto.WishlistFromModel(entity)
	}

	return result, nil
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

	// 处理附件关联
	if req.AttachmentIDs != nil {
		if err := s.AssociateAttachments(ctx, existingWishlist.ID, req.AttachmentIDs); err != nil {
			return nil, fmt.Errorf("关联附件失败: %w", err)
		}
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

// AssociateAttachments 将附件关联到心愿清单
func (s *wishlistService) AssociateAttachments(ctx context.Context, wishlistID int64, attachmentIDs []int64) error {
	// 首先验证心愿是否存在
	_, err := s.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		return err // GetWishlistByID 已经包装了错误
	}

	// 删除现有关联
	if err := s.wishlistAttachmentRepo.DeleteByWishlistID(ctx, wishlistID); err != nil {
		return fmt.Errorf("删除现有附件关联失败: %w", err)
	}

	// 创建新的关联
	for _, attachmentID := range attachmentIDs {
		// 验证附件是否存在
		_, err := s.attachmentRepo.GetByID(ctx, attachmentID)
		if err != nil {
			return fmt.Errorf("附件 %d 不存在: %w", attachmentID, err)
		}

		// 创建关联
		wishlistAttachment := &models.WishlistAttachment{
			WishlistID:   wishlistID,
			AttachmentID: attachmentID,
		}
		if err := s.wishlistAttachmentRepo.Create(ctx, wishlistAttachment); err != nil {
			return fmt.Errorf("创建附件关联失败: %w", err)
		}
	}

	return nil
}

// RemoveAttachment 从心愿清单中移除单个附件
func (s *wishlistService) RemoveAttachment(ctx context.Context, wishlistID int64, attachmentID int64) error {
	// 首先验证心愿是否存在
	_, err := s.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		return err // GetWishlistByID 已经包装了错误
	}

	// 删除关联
	if err := s.wishlistAttachmentRepo.DeleteByWishlistAndAttachmentID(ctx, wishlistID, attachmentID); err != nil {
		return fmt.Errorf("删除附件关联失败: %w", err)
	}

	return nil
}

// GetAttachments 获取心愿清单关联的所有附件
func (s *wishlistService) GetAttachments(ctx context.Context, wishlistID int64) ([]models.Attachment, error) {
	// 首先验证心愿是否存在
	_, err := s.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		return nil, err // GetWishlistByID 已经包装了错误
	}

	// 获取关联
	wishlistAttachments, err := s.wishlistAttachmentRepo.FindByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("获取附件关联失败: %w", err)
	}

	// 如果没有关联的附件，返回空数组
	if len(wishlistAttachments) == 0 {
		return []models.Attachment{}, nil
	}

	// 收集附件ID
	var attachmentIDs []int64
	for _, wa := range wishlistAttachments {
		attachmentIDs = append(attachmentIDs, wa.AttachmentID)
	}

	// 查询附件详情
	var attachments []models.Attachment
	for _, attachmentID := range attachmentIDs {
		attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
		if err != nil {
			if errors.Is(err, repository.ErrAttachmentNotFound) {
				// 如果附件不存在，跳过
				continue
			}
			return nil, fmt.Errorf("获取附件详情失败: %w", err)
		}
		attachments = append(attachments, *attachment)
	}

	return attachments, nil
}

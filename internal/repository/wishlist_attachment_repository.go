package repository

import (
	"context"
	"errors"
	"memoir-api/internal/models"

	"gorm.io/gorm"
)

// ErrWishlistAttachmentNotFound 心愿与附件关联不存在错误
var ErrWishlistAttachmentNotFound = errors.New("wishlist attachment not found")

// WishlistAttachmentRepository 心愿与附件关联仓库接口
type WishlistAttachmentRepository interface {
	Repository
	Create(ctx context.Context, wishlistAttachment *models.WishlistAttachment) error
	FindByWishlistID(ctx context.Context, wishlistID int64) ([]models.WishlistAttachment, error)
	FindByAttachmentID(ctx context.Context, attachmentID int64) ([]models.WishlistAttachment, error)
	DeleteByWishlistID(ctx context.Context, wishlistID int64) error
	DeleteByAttachmentID(ctx context.Context, attachmentID int64) error
	DeleteByWishlistAndAttachmentID(ctx context.Context, wishlistID, attachmentID int64) error
}

// wishlistAttachmentRepository 心愿与附件关联仓库实现
type wishlistAttachmentRepository struct {
	*BaseRepository
}

// NewWishlistAttachmentRepository 创建心愿与附件关联仓库
func NewWishlistAttachmentRepository(db *gorm.DB) WishlistAttachmentRepository {
	return &wishlistAttachmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建心愿与附件关联
func (r *wishlistAttachmentRepository) Create(ctx context.Context, wishlistAttachment *models.WishlistAttachment) error {
	return r.DB().WithContext(ctx).Create(wishlistAttachment).Error
}

// FindByWishlistID 根据心愿ID查询关联
func (r *wishlistAttachmentRepository) FindByWishlistID(ctx context.Context, wishlistID int64) ([]models.WishlistAttachment, error) {
	var wishlistAttachments []models.WishlistAttachment
	err := r.DB().WithContext(ctx).Where("wishlist_id = ?", wishlistID).Find(&wishlistAttachments).Error
	if err != nil {
		return nil, err
	}
	return wishlistAttachments, nil
}

// FindByAttachmentID 根据附件ID查询关联
func (r *wishlistAttachmentRepository) FindByAttachmentID(ctx context.Context, attachmentID int64) ([]models.WishlistAttachment, error) {
	var wishlistAttachments []models.WishlistAttachment
	err := r.DB().WithContext(ctx).Where("attachment_id = ?", attachmentID).Find(&wishlistAttachments).Error
	if err != nil {
		return nil, err
	}
	return wishlistAttachments, nil
}

// DeleteByWishlistID 根据心愿ID删除关联
func (r *wishlistAttachmentRepository) DeleteByWishlistID(ctx context.Context, wishlistID int64) error {
	return r.DB().WithContext(ctx).Where("wishlist_id = ?", wishlistID).Delete(&models.WishlistAttachment{}).Error
}

// DeleteByAttachmentID 根据附件ID删除关联
func (r *wishlistAttachmentRepository) DeleteByAttachmentID(ctx context.Context, attachmentID int64) error {
	return r.DB().WithContext(ctx).Where("attachment_id = ?", attachmentID).Delete(&models.WishlistAttachment{}).Error
}

// DeleteByWishlistAndAttachmentID 根据心愿ID和附件ID删除单个关联
func (r *wishlistAttachmentRepository) DeleteByWishlistAndAttachmentID(ctx context.Context, wishlistID, attachmentID int64) error {
	return r.DB().WithContext(ctx).Where("wishlist_id = ? AND attachment_id = ?", wishlistID, attachmentID).Delete(&models.WishlistAttachment{}).Error
}

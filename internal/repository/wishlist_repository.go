package repository

import (
	"context"
	"errors"
	"time"

	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrWishlistNotFound = errors.New("心愿不存在")
)

// WishlistRepository 心愿清单仓库接口
type WishlistRepository interface {
	Repository
	Create(ctx context.Context, wishlist *models.Wishlist) error
	GetByID(ctx context.Context, id int64) (*models.Wishlist, error)
	ListByCoupleID(ctx context.Context, coupleID int64) ([]*models.Wishlist, error)
	ListByStatus(ctx context.Context, coupleID int64, status string) ([]*models.Wishlist, error)
	ListByPriority(ctx context.Context, coupleID int64, priority int) ([]*models.Wishlist, error)
	ListUpcomingReminders(ctx context.Context, daysAhead int) ([]*models.Wishlist, error)
	Update(ctx context.Context, wishlist *models.Wishlist) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
}

// wishlistRepository 心愿清单仓库实现
type wishlistRepository struct {
	*BaseRepository
}

// NewWishlistRepository 创建心愿清单仓库
func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建心愿
func (r *wishlistRepository) Create(ctx context.Context, wishlist *models.Wishlist) error {
	return r.DB().WithContext(ctx).Create(wishlist).Error
}

// GetByID 通过ID获取心愿
func (r *wishlistRepository) GetByID(ctx context.Context, id int64) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.DB().WithContext(ctx).First(&wishlist, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}
	return &wishlist, nil
}

// ListByCoupleID 获取情侣关系下的所有心愿，按优先级和创建时间排序
func (r *wishlistRepository) ListByCoupleID(ctx context.Context, coupleID int64) ([]*models.Wishlist, error) {
	var wishlists []*models.Wishlist
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.Wishlist{}).Where("couple_id = ?", coupleID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取列表，优先级从高到低（数字从小到大），同优先级的按创建时间从新到旧
	query := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID).Order("priority ASC, created_at DESC")

	if err := query.Find(&wishlists).Error; err != nil {
		return nil, err
	}

	return wishlists, nil
}

// ListByStatus 按状态获取心愿
func (r *wishlistRepository) ListByStatus(ctx context.Context, coupleID int64, status string) ([]*models.Wishlist, error) {
	var wishlists []*models.Wishlist
	err := r.DB().WithContext(ctx).
		Where("couple_id = ? AND status = ?", coupleID, status).
		Order("priority ASC, created_at DESC").
		Find(&wishlists).Error

	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

// ListByPriority 按优先级获取心愿
func (r *wishlistRepository) ListByPriority(ctx context.Context, coupleID int64, priority int) ([]*models.Wishlist, error) {
	var wishlists []*models.Wishlist
	err := r.DB().WithContext(ctx).
		Where("couple_id = ? AND priority = ?", coupleID, priority).
		Order("created_at DESC").
		Find(&wishlists).Error

	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

// ListUpcomingReminders 获取即将到期的提醒
func (r *wishlistRepository) ListUpcomingReminders(ctx context.Context, daysAhead int) ([]*models.Wishlist, error) {
	var wishlists []*models.Wishlist
	today := time.Now().Truncate(24 * time.Hour)
	future := today.AddDate(0, 0, daysAhead)

	err := r.DB().WithContext(ctx).
		Where("reminder_date IS NOT NULL AND reminder_date BETWEEN ? AND ? AND status != ?", today, future, "completed").
		Order("reminder_date ASC").
		Find(&wishlists).Error

	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

// Update 更新心愿
func (r *wishlistRepository) Update(ctx context.Context, wishlist *models.Wishlist) error {
	result := r.DB().WithContext(ctx).Save(wishlist)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

// UpdateStatus 更新心愿状态
func (r *wishlistRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	result := r.DB().WithContext(ctx).Model(&models.Wishlist{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

// Delete 删除心愿
func (r *wishlistRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB().WithContext(ctx).Delete(&models.Wishlist{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

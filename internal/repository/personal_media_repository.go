package repository

import (
	"context"

	"gorm.io/gorm"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"
)

// PersonalMediaRepository 个人媒体仓库接口
type PersonalMediaRepository interface {
	// 创建个人媒体
	Create(ctx context.Context, media *models.PersonalMedia) error
	// 获取用户所有个人媒体
	FindByUserID(ctx context.Context, userID int64, category string) ([]models.PersonalMedia, error)
	// 获取单个个人媒体
	FindByID(ctx context.Context, id int64) (*models.PersonalMedia, error)
	// 更新个人媒体
	Update(ctx context.Context, media *models.PersonalMedia) error
	// 删除个人媒体
	Delete(ctx context.Context, id int64) error
	// 查询个人媒体（支持分页和筛选）
	Query(ctx context.Context, pageRequest dto.QueryPersonalMediaRequest) ([]models.PersonalMedia, int64, error)
}

// GormPersonalMediaRepository 个人媒体仓库的GORM实现
type GormPersonalMediaRepository struct {
	db *gorm.DB
}

// NewGormPersonalMediaRepository 创建个人媒体仓库实例
func NewGormPersonalMediaRepository(db *gorm.DB) PersonalMediaRepository {
	return &GormPersonalMediaRepository{
		db: db,
	}
}

// Create 创建个人媒体记录
func (r *GormPersonalMediaRepository) Create(ctx context.Context, media *models.PersonalMedia) error {
	return r.db.WithContext(ctx).Create(media).Error
}

// FindByUserID 获取用户所有个人媒体
func (r *GormPersonalMediaRepository) FindByUserID(ctx context.Context, userID int64, category string) ([]models.PersonalMedia, error) {
	var media []models.PersonalMedia
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("created_at DESC").Find(&media).Error
	return media, err
}

// FindByID 获取单个个人媒体
func (r *GormPersonalMediaRepository) FindByID(ctx context.Context, id int64) (*models.PersonalMedia, error) {
	var media models.PersonalMedia
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

// Update 更新个人媒体
func (r *GormPersonalMediaRepository) Update(ctx context.Context, media *models.PersonalMedia) error {
	return r.db.WithContext(ctx).Save(media).Error
}

// Delete 删除个人媒体
func (r *GormPersonalMediaRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.PersonalMedia{}, id).Error
}

// Query 查询个人媒体（支持分页和筛选）
func (r *GormPersonalMediaRepository) Query(ctx context.Context, pageRequest dto.QueryPersonalMediaRequest) ([]models.PersonalMedia, int64, error) {
	var media []models.PersonalMedia
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PersonalMedia{}).Where("user_id = ?", pageRequest.UserID)

	if pageRequest.Category != "" {
		query = query.Where("category = ?", pageRequest.Category)
	}

	if pageRequest.MediaType != "" {
		query = query.Where("media_type = ?", pageRequest.MediaType)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := pageRequest.Offset()
	limit := pageRequest.Limit()
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error
	if err != nil {
		return nil, 0, err
	}

	return media, total, nil
}

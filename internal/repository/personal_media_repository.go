package repository

import (
	"context"
	"net/url"
	"strings"

	"gorm.io/gorm"

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
	Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error)
	// 迁移现有媒体数据，添加路径信息
	MigrateExistingMedia(ctx context.Context) error
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
func (r *GormPersonalMediaRepository) Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error) {
	var media []models.PersonalMedia
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PersonalMedia{}).Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if mediaType != "" {
		query = query.Where("media_type = ?", mediaType)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&media).Error
	if err != nil {
		return nil, 0, err
	}

	return media, total, nil
}

// MigrateExistingMedia 为现有媒体记录添加路径信息
func (r *GormPersonalMediaRepository) MigrateExistingMedia(ctx context.Context) error {
	// 查找所有没有path字段的记录
	var mediaList []models.PersonalMedia
	err := r.db.WithContext(ctx).Where("path = '' OR path IS NULL").Find(&mediaList).Error
	if err != nil {
		return err
	}

	// 处理每条记录
	for _, media := range mediaList {
		// 从URL中提取路径
		parsedURL, err := url.Parse(media.MediaURL)
		if err != nil {
			continue // 跳过无法解析的URL
		}

		// 获取路径部分（去除开头的斜杠）
		path := strings.TrimPrefix(parsedURL.Path, "/")
		if path == "" {
			continue // 跳过无法提取路径的记录
		}

		// 更新记录
		err = r.db.WithContext(ctx).Model(&models.PersonalMedia{}).Where("id = ?", media.ID).Update("path", path).Error
		if err != nil {
			return err
		}
	}

	return nil
}

package repository

import (
	"context"
	"errors"
	"memoir-api/internal/api/dto"

	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrPhotoVideoNotFound = errors.New("照片/视频不存在")
)

// PhotoVideoRepository 照片和视频仓库接口
type PhotoVideoRepository interface {
	Repository
	Create(ctx context.Context, photoVideo *models.PhotoVideo) error
	GetByID(ctx context.Context, id int64) (*models.PhotoVideo, error)
	Query(ctx context.Context, params *dto.PhotoVideoQueryParams) ([]*models.PhotoVideo, int64, error)
	Update(ctx context.Context, photoVideo *models.PhotoVideo) error
	Delete(ctx context.Context, id int64) error
}

// photoVideoRepository 照片和视频仓库实现
type photoVideoRepository struct {
	*BaseRepository
}

func (r *photoVideoRepository) Query(ctx context.Context, params *dto.PhotoVideoQueryParams) ([]*models.PhotoVideo, int64, error) {
	db := r.DB().WithContext(ctx).Model(&models.PhotoVideo{}).Order("created_at DESC")
	// 构建查询条件
	if params.CoupleID != 0 {
		db = db.Where("couple_id = ?", params.CoupleID)
	}
	if params.AlbumID != 0 {
		db = db.Where("album_id = ?", params.AlbumID)
	}
	if params.MediaType != "" {
		db = db.Where("media_type = ?", params.MediaType)
	}
	if params.EventID != 0 {
		db = db.Where("event_id = ?", params.EventID)
	}
	if params.LocationID != 0 {
		db = db.Where("location_id = ?", params.LocationID)
	}
	var total int64

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	limit, offset := params.Limit(), params.Limit()
	if limit >= 0 && offset > 0 {
		db = db.Offset(offset).Limit(limit)
	}
	var results []*models.PhotoVideo
	if err := db.Find(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// NewPhotoVideoRepository 创建照片和视频仓库
func NewPhotoVideoRepository(db *gorm.DB) PhotoVideoRepository {
	return &photoVideoRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建照片/视频
func (r *photoVideoRepository) Create(ctx context.Context, photoVideo *models.PhotoVideo) error {
	return r.DB().WithContext(ctx).Create(photoVideo).Error
}

// GetByID 通过ID获取照片/视频
func (r *photoVideoRepository) GetByID(ctx context.Context, id int64) (*models.PhotoVideo, error) {
	var photoVideo models.PhotoVideo
	err := r.DB().WithContext(ctx).First(&photoVideo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPhotoVideoNotFound
		}
		return nil, err
	}
	return &photoVideo, nil
}

// ListByCoupleID 获取情侣关系下的所有照片/视频，按上传时间倒序排列
func (r *photoVideoRepository) ListByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	var photoVideos []*models.PhotoVideo
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.PhotoVideo{}).Where("couple_id = ?", coupleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID).Order("created_at DESC")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&photoVideos).Error; err != nil {
		return nil, 0, err
	}

	return photoVideos, total, nil
}

// ListByCategory 按分类获取照片/视频
func (r *photoVideoRepository) ListByCategory(ctx context.Context, coupleID int64, category string, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	var photoVideos []*models.PhotoVideo
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.PhotoVideo{}).Where("couple_id = ? AND category = ?", coupleID, category).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx).Where("couple_id = ? AND category = ?", coupleID, category).Order("created_at DESC")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&photoVideos).Error; err != nil {
		return nil, 0, err
	}

	return photoVideos, total, nil
}

// ListByMediaType 按媒体类型获取照片/视频
func (r *photoVideoRepository) ListByMediaType(ctx context.Context, coupleID int64, mediaType string, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	var photoVideos []*models.PhotoVideo
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.PhotoVideo{}).Where("couple_id = ? AND media_type = ?", coupleID, mediaType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx).Where("couple_id = ? AND media_type = ?", coupleID, mediaType).Order("created_at DESC")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&photoVideos).Error; err != nil {
		return nil, 0, err
	}

	return photoVideos, total, nil
}

// ListByEventID 按事件ID获取照片/视频
func (r *photoVideoRepository) ListByEventID(ctx context.Context, eventID int64) ([]*models.PhotoVideo, error) {
	var photoVideos []*models.PhotoVideo
	err := r.DB().WithContext(ctx).Where("event_id = ?", eventID).Order("created_at DESC").Find(&photoVideos).Error
	if err != nil {
		return nil, err
	}
	return photoVideos, nil
}

// ListByLocationID 按地点ID获取照片/视频
func (r *photoVideoRepository) ListByLocationID(ctx context.Context, locationID int64) ([]*models.PhotoVideo, error) {
	var photoVideos []*models.PhotoVideo
	err := r.DB().WithContext(ctx).Where("location_id = ?", locationID).Order("created_at DESC").Find(&photoVideos).Error
	if err != nil {
		return nil, err
	}
	return photoVideos, nil
}

// Update 更新照片/视频
func (r *photoVideoRepository) Update(ctx context.Context, photoVideo *models.PhotoVideo) error {
	result := r.DB().WithContext(ctx).Save(photoVideo)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPhotoVideoNotFound
	}
	return nil
}

// Delete 删除照片/视频
func (r *photoVideoRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB().WithContext(ctx).Delete(&models.PhotoVideo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPhotoVideoNotFound
	}
	return nil
}

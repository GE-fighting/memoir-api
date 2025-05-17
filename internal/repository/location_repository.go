package repository

import (
	"context"
	"errors"

	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrLocationNotFound = errors.New("地点不存在")
)

// LocationRepository 地点仓库接口
type LocationRepository interface {
	Repository
	Create(ctx context.Context, location *models.Location) error
	GetByID(ctx context.Context, id int64) (*models.Location, error)
	ListByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.Location, int64, error)
	Update(ctx context.Context, location *models.Location) error
	Delete(ctx context.Context, id int64) error
	// 根据地理位置查询附近地点
	FindNearby(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*models.Location, error)
}

// locationRepository 地点仓库实现
type locationRepository struct {
	*BaseRepository
}

// NewLocationRepository 创建地点仓库
func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建地点
func (r *locationRepository) Create(ctx context.Context, location *models.Location) error {
	return r.DB().WithContext(ctx).Create(location).Error
}

// GetByID 通过ID获取地点
func (r *locationRepository) GetByID(ctx context.Context, id int64) (*models.Location, error) {
	var location models.Location
	err := r.DB().WithContext(ctx).First(&location, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return &location, nil
}

// ListByCoupleID 获取情侣关系下的所有地点
func (r *locationRepository) ListByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.Location, int64, error) {
	var locations []*models.Location
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.Location{}).Where("couple_id = ?", coupleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&locations).Error; err != nil {
		return nil, 0, err
	}

	return locations, total, nil
}

// Update 更新地点
func (r *locationRepository) Update(ctx context.Context, location *models.Location) error {
	result := r.DB().WithContext(ctx).Save(location)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLocationNotFound
	}
	return nil
}

// Delete 删除地点
func (r *locationRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB().WithContext(ctx).Delete(&models.Location{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLocationNotFound
	}
	return nil
}

// FindNearby 查找附近地点
func (r *locationRepository) FindNearby(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*models.Location, error) {
	var locations []*models.Location

	// 使用PostGIS的ST_DWithin函数查询指定半径内的地点
	// radiusMeters需要除以111111将米转换为度（1度约等于111.111公里）
	radius := radiusMeters / 111111.0

	query := r.DB().WithContext(ctx).Raw(`
		SELECT * FROM locations
		WHERE ST_DWithin(
			coordinates,
			ST_SetSRID(ST_MakePoint(?, ?), 4326),
			?
		)
		ORDER BY ST_Distance(
			coordinates,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)
		)
		LIMIT ?
	`, lng, lat, radius, lng, lat, limit)

	if err := query.Scan(&locations).Error; err != nil {
		return nil, err
	}

	return locations, nil
}

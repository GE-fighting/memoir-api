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
	FindByID(ctx context.Context, id int64) (*models.Location, error)
	FindByIDs(ctx context.Context, ids []int64) ([]models.Location, error)
	FindByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]models.Location, int64, error)
}

// locationRepository 地点仓库实现
type locationRepository struct {
	*BaseRepository
}

func (r *locationRepository) FindByID(ctx context.Context, id int64) (*models.Location, error) {
	//TODO implement me
	panic("implement me")
}

func (r *locationRepository) FindByIDs(ctx context.Context, ids []int64) ([]models.Location, error) {
	var locations []models.Location
	err := r.DB().WithContext(ctx).Where("id IN ?", ids).Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *locationRepository) FindByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]models.Location, int64, error) {
	//TODO implement me
	panic("implement me")
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



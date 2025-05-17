package repository

import (
	"context"
	"errors"

	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrCoupleNotFound = errors.New("情侣关系不存在")
)

// CoupleRepository 情侣关系仓库接口
type CoupleRepository interface {
	Repository
	Create(ctx context.Context, couple *models.Couple) error
	GetByID(ctx context.Context, id int64) (*models.Couple, error)
	GetByPairToken(ctx context.Context, pairToken string) (*models.Couple, error)
	List(ctx context.Context, offset, limit int) ([]*models.Couple, int64, error)
	Update(ctx context.Context, couple *models.Couple) error
	Delete(ctx context.Context, id int64) error
}

// coupleRepository 情侣关系仓库实现
type coupleRepository struct {
	*BaseRepository
}

// NewCoupleRepository 创建情侣关系仓库
func NewCoupleRepository(db *gorm.DB) CoupleRepository {
	return &coupleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建情侣关系
func (r *coupleRepository) Create(ctx context.Context, couple *models.Couple) error {
	return r.DB().WithContext(ctx).Create(couple).Error
}

// GetByID 通过ID获取情侣关系
func (r *coupleRepository) GetByID(ctx context.Context, id int64) (*models.Couple, error) {
	var couple models.Couple
	err := r.DB().WithContext(ctx).First(&couple, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCoupleNotFound
		}
		return nil, err
	}
	return &couple, nil
}

// GetByPairToken 通过配对令牌获取情侣关系
func (r *coupleRepository) GetByPairToken(ctx context.Context, pairToken string) (*models.Couple, error) {
	var couple models.Couple
	err := r.DB().WithContext(ctx).Where("pair_token = ?", pairToken).First(&couple).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCoupleNotFound
		}
		return nil, err
	}
	return &couple, nil
}

// List 获取情侣关系列表
func (r *coupleRepository) List(ctx context.Context, offset, limit int) ([]*models.Couple, int64, error) {
	var couples []*models.Couple
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.Couple{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&couples).Error; err != nil {
		return nil, 0, err
	}

	return couples, total, nil
}

// Update 更新情侣关系
func (r *coupleRepository) Update(ctx context.Context, couple *models.Couple) error {
	result := r.DB().WithContext(ctx).Save(couple)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCoupleNotFound
	}
	return nil
}

// Delete 删除情侣关系
func (r *coupleRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB().WithContext(ctx).Delete(&models.Couple{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCoupleNotFound
	}
	return nil
}

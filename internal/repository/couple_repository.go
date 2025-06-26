package repository

import (
	"context"
	"errors"

	"memoir-api/internal/logger"
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
	logger logger.Logger
}

// NewCoupleRepository 创建情侣关系仓库
func NewCoupleRepository(db *gorm.DB) CoupleRepository {
	return &coupleRepository{
		BaseRepository: NewBaseRepository(db),
		logger:         logger.GetLogger("couple_repository"),
	}
}

// Create 创建情侣关系
func (r *coupleRepository) Create(ctx context.Context, couple *models.Couple) error {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Creating couple relationship",
		"pair_token", couple.PairToken,
		"anniversary_date", couple.AnniversaryDate,
	)

	err := r.DB().WithContext(ctx).Create(couple).Error
	if err != nil {
		log.Error(err, "Failed to create couple relationship")
		return err
	}

	log.Info("Created couple relationship successfully", "couple_id", couple.ID)

	return nil
}

// GetByID 通过ID获取情侣关系
func (r *coupleRepository) GetByID(ctx context.Context, id int64) (*models.Couple, error) {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Getting couple by ID", "couple_id", id)

	var couple models.Couple
	err := r.DB().WithContext(ctx).First(&couple, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug("Couple not found", "couple_id", id)
			return nil, ErrCoupleNotFound
		}
		log.Error(err, "Failed to get couple by ID", "couple_id", id)
		return nil, err
	}

	return &couple, nil
}

// GetByPairToken 通过配对令牌获取情侣关系
func (r *coupleRepository) GetByPairToken(ctx context.Context, pairToken string) (*models.Couple, error) {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Getting couple by pair token", "pair_token", pairToken)

	var couple models.Couple
	err := r.DB().WithContext(ctx).Where("pair_token = ?", pairToken).First(&couple).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug("Couple not found by pair token")
			return nil, ErrCoupleNotFound
		}
		log.Error(err, "Failed to get couple by pair token")
		return nil, err
	}

	return &couple, nil
}

// List 获取情侣关系列表
func (r *coupleRepository) List(ctx context.Context, offset, limit int) ([]*models.Couple, int64, error) {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Listing couples", "offset", offset, "limit", limit)

	var couples []*models.Couple
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.Couple{}).Count(&total).Error; err != nil {
		log.Error(err, "Failed to count couples")
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx)
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&couples).Error; err != nil {
		log.Error(err, "Failed to list couples")
		return nil, 0, err
	}

	log.Debug("Listed couples successfully", "count", len(couples), "total", total)

	return couples, total, nil
}

// Update 更新情侣关系
func (r *coupleRepository) Update(ctx context.Context, couple *models.Couple) error {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Updating couple", "couple_id", couple.ID)

	result := r.DB().WithContext(ctx).Save(couple)
	if result.Error != nil {
		log.Error(result.Error, "Failed to update couple", "couple_id", couple.ID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Debug("No couple updated, couple not found", "couple_id", couple.ID)
		return ErrCoupleNotFound
	}

	log.Info("Updated couple successfully", "couple_id", couple.ID)

	return nil
}

// Delete 删除情侣关系
func (r *coupleRepository) Delete(ctx context.Context, id int64) error {
	log := logger.FromContext(ctx).WithComponent("couple_repository")

	log.Debug("Deleting couple", "couple_id", id)

	result := r.DB().WithContext(ctx).Delete(&models.Couple{}, id)
	if result.Error != nil {
		log.Error(result.Error, "Failed to delete couple", "couple_id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Debug("No couple deleted, couple not found", "couple_id", id)
		return ErrCoupleNotFound
	}

	log.Info("Deleted couple successfully", "couple_id", id)

	return nil
}

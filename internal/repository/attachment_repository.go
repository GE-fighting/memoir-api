package repository

import (
	"context"
	"errors"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrAttachmentNotFound = errors.New("附件不存在")
)

// AttachmentRepository 附件仓库接口
type AttachmentRepository interface {
	Repository
	Create(ctx context.Context, attachment *models.Attachment) error
	GetByID(ctx context.Context, id int64) (*models.Attachment, error)
	Query(ctx context.Context, params *dto.AttachmentQueryParams) ([]models.Attachment, int64, error)
	Update(ctx context.Context, attachment *models.Attachment) error
	Delete(ctx context.Context, id int64) error
	ListByUserID(ctx context.Context, userID int64, spaceType string) ([]models.Attachment, error)
	ListByCoupleID(ctx context.Context, coupleID int64) ([]models.Attachment, error)
}

// attachmentRepository 附件仓库实现
type attachmentRepository struct {
	*BaseRepository
}

// NewAttachmentRepository 创建附件仓库
func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建附件
func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	return r.DB().WithContext(ctx).Create(attachment).Error
}

// GetByID 通过ID获取附件
func (r *attachmentRepository) GetByID(ctx context.Context, id int64) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.DB().WithContext(ctx).First(&attachment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttachmentNotFound
		}
		return nil, err
	}
	return &attachment, nil
}

// Query 查询附件（支持分页和筛选）
func (r *attachmentRepository) Query(ctx context.Context, params *dto.AttachmentQueryParams) ([]models.Attachment, int64, error) {
	db := r.DB().WithContext(ctx).Model(&models.Attachment{})

	// 应用过滤条件
	if params.UserID != 0 {
		db = db.Where("user_id = ?", params.UserID)
	}

	if params.CoupleID != 0 {
		db = db.Where("couple_id = ?", params.CoupleID)
	}

	if params.SpaceType != "" {
		db = db.Where("space_type = ?", params.SpaceType)
	}

	if params.FileType != "" {
		db = db.Where("file_type = ?", params.FileType)
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页并获取数据
	var attachments []models.Attachment
	if err := db.Offset(params.Offset()).Limit(params.Limit()).Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, 0, err
	}

	return attachments, total, nil
}

// Update 更新附件
func (r *attachmentRepository) Update(ctx context.Context, attachment *models.Attachment) error {
	return r.DB().WithContext(ctx).Save(attachment).Error
}

// Delete 删除附件
func (r *attachmentRepository) Delete(ctx context.Context, id int64) error {
	return r.DB().WithContext(ctx).Delete(&models.Attachment{}, id).Error
}

// ListByUserID 根据用户ID列出附件
func (r *attachmentRepository) ListByUserID(ctx context.Context, userID int64, spaceType string) ([]models.Attachment, error) {
	var attachments []models.Attachment
	query := r.DB().WithContext(ctx).Where("user_id = ?", userID)

	if spaceType != "" {
		query = query.Where("space_type = ?", spaceType)
	}

	err := query.Order("created_at DESC").Find(&attachments).Error
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

// ListByCoupleID 根据情侣ID列出附件
func (r *attachmentRepository) ListByCoupleID(ctx context.Context, coupleID int64) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.DB().WithContext(ctx).
		Where("couple_id = ?", coupleID).
		Where("space_type = ?", "couple").
		Order("created_at DESC").
		Find(&attachments).Error
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

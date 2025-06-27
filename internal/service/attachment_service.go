package service

import (
	"context"
	"errors"
	"fmt"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

var (
	ErrAttachmentNotFound = errors.New("附件不存在")
)

// AttachmentService 附件服务接口
type AttachmentService interface {
	Service
	CreateAttachment(ctx context.Context, request *dto.CreateAttachmentRequest) (*models.Attachment, error)
	GetAttachmentByID(ctx context.Context, id int64) (*models.Attachment, error)
	QueryAttachments(ctx context.Context, params *dto.AttachmentQueryParams) (*dto.PageResult, error)
	DeleteAttachment(ctx context.Context, id int64) error
	ListByUserID(ctx context.Context, userID int64, spaceType string) ([]models.Attachment, error)
	ListByCoupleID(ctx context.Context, coupleID int64) ([]models.Attachment, error)
}

// attachmentService 附件服务实现
type attachmentService struct {
	*BaseService
	repo       repository.AttachmentRepository
	userRepo   repository.UserRepository
	coupleRepo repository.CoupleRepository
}

// NewAttachmentService 创建附件服务
func NewAttachmentService(
	repo repository.AttachmentRepository,
	userRepo repository.UserRepository,
	coupleRepo repository.CoupleRepository,
) AttachmentService {
	return &attachmentService{
		BaseService: NewBaseService(repo),
		repo:        repo,
		userRepo:    userRepo,
		coupleRepo:  coupleRepo,
	}
}

// CreateAttachment 创建附件
func (s *attachmentService) CreateAttachment(ctx context.Context, request *dto.CreateAttachmentRequest) (*models.Attachment, error) {
	// 校验用户是否存在
	_, err := s.userRepo.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 如果是情侣空间，校验情侣关系是否存在
	if request.SpaceType == "couple" && request.CoupleID > 0 {
		_, err := s.coupleRepo.GetByID(ctx, request.CoupleID)
		if err != nil {
			return nil, fmt.Errorf("查询情侣关系失败: %w", err)
		}
	}

	// 创建附件模型
	attachment := request.ToModel()

	// 保存附件
	if err := s.repo.Create(ctx, attachment); err != nil {
		return nil, fmt.Errorf("创建附件失败: %w", err)
	}

	return attachment, nil
}

// GetAttachmentByID 通过ID获取附件
func (s *attachmentService) GetAttachmentByID(ctx context.Context, id int64) (*models.Attachment, error) {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAttachmentNotFound) {
			return nil, ErrAttachmentNotFound
		}
		return nil, fmt.Errorf("获取附件失败: %w", err)
	}
	return attachment, nil
}

// QueryAttachments 查询附件
func (s *attachmentService) QueryAttachments(ctx context.Context, params *dto.AttachmentQueryParams) (*dto.PageResult, error) {
	attachments, total, err := s.repo.Query(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("查询附件失败: %w", err)
	}

	// 将模型转换为DTO
	attachmentDTOs := dto.AttachmentsFromModels(attachments)

	// 创建分页结果
	result := dto.NewPageResult(attachmentDTOs, total, params.Page, params.PageSize)
	return &result, nil
}

// DeleteAttachment 删除附件
func (s *attachmentService) DeleteAttachment(ctx context.Context, id int64) error {
	// 检查附件是否存在
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAttachmentNotFound) {
			return ErrAttachmentNotFound
		}
		return fmt.Errorf("查询附件失败: %w", err)
	}

	// 删除附件
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除附件失败: %w", err)
	}

	return nil
}

// ListByUserID 通过用户ID列出附件
func (s *attachmentService) ListByUserID(ctx context.Context, userID int64, spaceType string) ([]models.Attachment, error) {
	return s.repo.ListByUserID(ctx, userID, spaceType)
}

// ListByCoupleID 通过情侣ID列出附件
func (s *attachmentService) ListByCoupleID(ctx context.Context, coupleID int64) ([]models.Attachment, error) {
	return s.repo.ListByCoupleID(ctx, coupleID)
}

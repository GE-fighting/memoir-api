package service

import (
	"context"

	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

// PersonalMediaService 个人媒体服务接口
type PersonalMediaService interface {
	// 通过URL创建个人媒体（前端直接上传到OSS）
	CreateWithURL(ctx context.Context, request dto.CreatePersonalMediaWithURLRequest) (*models.PersonalMedia, error)
	// 分页查询个人媒体
	PageQuery(ctx context.Context, pageRequest dto.QueryPersonalMediaRequest) ([]models.PersonalMedia, int64, error)
}

// DefaultPersonalMediaService 个人媒体服务的默认实现
type DefaultPersonalMediaService struct {
	repo repository.PersonalMediaRepository
}

// NewPersonalMediaService 创建个人媒体服务实例
func NewPersonalMediaService(repo repository.PersonalMediaRepository) PersonalMediaService {
	return &DefaultPersonalMediaService{
		repo: repo,
	}
}

// CreateWithURL 通过URL创建个人媒体（前端直接上传到OSS）
func (s *DefaultPersonalMediaService) CreateWithURL(ctx context.Context, request dto.CreatePersonalMediaWithURLRequest) (*models.PersonalMedia, error) {

	// 创建个人媒体记录
	media := &models.PersonalMedia{
		UserID:       request.UserID,
		MediaURL:     request.MediaURL,
		MediaType:    request.MediaType,
		Category:     &request.Category,
		ThumbnailURL: &request.ThumbnailURL,
		Description:  &request.Description,
		Title:        &request.Title,
	}

	err := s.repo.Create(ctx, media)
	if err != nil {
		return nil, err
	}

	return media, nil
}

// 分页查询 查询个人媒体
func (s *DefaultPersonalMediaService) PageQuery(ctx context.Context, pageRequest dto.QueryPersonalMediaRequest) ([]models.PersonalMedia, int64, error) {
	return s.repo.Query(ctx, pageRequest)
}

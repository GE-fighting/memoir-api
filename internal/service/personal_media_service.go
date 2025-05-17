package service

import (
	"context"
	"errors"
	"mime/multipart"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

// PersonalMediaService 个人媒体服务接口
type PersonalMediaService interface {
	// 创建个人媒体
	Create(ctx context.Context, userID int64, mediaType, category, title string, description []byte, file *multipart.FileHeader, isPrivate bool, tags []string) (*models.PersonalMedia, error)
	// 获取用户所有个人媒体
	GetByUserID(ctx context.Context, userID int64, category string) ([]models.PersonalMedia, error)
	// 获取单个个人媒体
	GetByID(ctx context.Context, id, userID int64) (*models.PersonalMedia, error)
	// 更新个人媒体
	Update(ctx context.Context, id, userID int64, title string, description []byte, isPrivate bool, tags []string) (*models.PersonalMedia, error)
	// 删除个人媒体
	Delete(ctx context.Context, id, userID int64) error
	// 查询个人媒体（支持分页和筛选）
	Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error)
}

// DefaultPersonalMediaService 个人媒体服务的默认实现
type DefaultPersonalMediaService struct {
	repo          repository.PersonalMediaRepository
	mediaUploader MediaUploader // 使用现有的媒体上传服务
}

// NewPersonalMediaService 创建个人媒体服务实例
func NewPersonalMediaService(repo repository.PersonalMediaRepository, mediaUploader MediaUploader) PersonalMediaService {
	return &DefaultPersonalMediaService{
		repo:          repo,
		mediaUploader: mediaUploader,
	}
}

// Create 创建个人媒体
func (s *DefaultPersonalMediaService) Create(ctx context.Context, userID int64, mediaType, category, title string, description []byte, file *multipart.FileHeader, isPrivate bool, tags []string) (*models.PersonalMedia, error) {
	// 上传文件到存储服务
	mediaURL, thumbnailURL, err := s.mediaUploader.UploadMedia(ctx, file, mediaType)
	if err != nil {
		return nil, err
	}

	// 创建个人媒体记录
	media := &models.PersonalMedia{
		UserID:       userID,
		MediaURL:     mediaURL,
		MediaType:    mediaType,
		Category:     category,
		ThumbnailURL: &thumbnailURL,
		Description:  description,
		Title:        title,
		IsPrivate:    isPrivate,
		Tags:         tags,
	}

	err = s.repo.Create(ctx, media)
	if err != nil {
		return nil, err
	}

	return media, nil
}

// GetByUserID 获取用户所有个人媒体
func (s *DefaultPersonalMediaService) GetByUserID(ctx context.Context, userID int64, category string) ([]models.PersonalMedia, error) {
	return s.repo.FindByUserID(ctx, userID, category)
}

// GetByID 获取单个个人媒体，并检查是否属于该用户
func (s *DefaultPersonalMediaService) GetByID(ctx context.Context, id, userID int64) (*models.PersonalMedia, error) {
	media, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查是否属于该用户
	if media.UserID != userID {
		return nil, errors.New("无权访问此媒体")
	}

	return media, nil
}

// Update 更新个人媒体
func (s *DefaultPersonalMediaService) Update(ctx context.Context, id, userID int64, title string, description []byte, isPrivate bool, tags []string) (*models.PersonalMedia, error) {
	// 先获取媒体，确保存在并属于该用户
	media, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	media.Title = title
	media.Description = description
	media.IsPrivate = isPrivate
	media.Tags = tags

	// 保存更新
	err = s.repo.Update(ctx, media)
	if err != nil {
		return nil, err
	}

	return media, nil
}

// Delete 删除个人媒体
func (s *DefaultPersonalMediaService) Delete(ctx context.Context, id, userID int64) error {
	// 先获取媒体，确保存在并属于该用户
	media, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// 删除媒体
	return s.repo.Delete(ctx, media.ID)
}

// Query 查询个人媒体（支持分页和筛选）
func (s *DefaultPersonalMediaService) Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error) {
	return s.repo.Query(ctx, userID, category, mediaType, page, pageSize)
}

package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

// PersonalMediaService 个人媒体服务接口
type PersonalMediaService interface {
	// 通过URL创建个人媒体（前端直接上传到OSS）
	CreateWithURL(ctx context.Context, userID int64, mediaType, category, title string, description []byte, mediaURL, thumbnailURL string, isPrivate bool, tags []string) (*models.PersonalMedia, error)

	// 查询个人媒体（支持分页和筛选）
	Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error)

	// 获取单个个人媒体
	GetByID(ctx context.Context, id int64, userID int64) (*models.PersonalMedia, error)

	// 更新个人媒体
	Update(ctx context.Context, id int64, userID int64, title string, description []byte, isPrivate bool, tags []string) (*models.PersonalMedia, error)

	// 删除个人媒体
	Delete(ctx context.Context, id int64, userID int64) error

	// 迁移现有媒体数据，添加路径信息
	MigrateExistingMedia(ctx context.Context) error
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

// 从URL提取路径信息
func extractPathFromURL(mediaURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return "", err
	}

	// 获取路径部分（去除开头的斜杠）
	path := strings.TrimPrefix(parsedURL.Path, "/")
	if path == "" {
		return "", errors.New("无法从URL提取路径")
	}

	return path, nil
}

// CreateWithURL 通过URL创建个人媒体（前端直接上传到OSS）
func (s *DefaultPersonalMediaService) CreateWithURL(ctx context.Context, userID int64, mediaType, category, title string, description []byte, mediaURL, thumbnailURL string, isPrivate bool, tags []string) (*models.PersonalMedia, error) {
	// 从URL提取完整路径
	path, err := extractPathFromURL(mediaURL)
	if err != nil {
		// 如果提取失败，使用空字符串
		path = ""
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
		Path:         path, // 保存路径信息
	}

	err = s.repo.Create(ctx, media)
	if err != nil {
		return nil, err
	}

	return media, nil
}

// Query 查询个人媒体
func (s *DefaultPersonalMediaService) Query(ctx context.Context, userID int64, category, mediaType string, page, pageSize int) ([]models.PersonalMedia, int64, error) {
	return s.repo.Query(ctx, userID, category, mediaType, page, pageSize)
}

// GetByID 获取单个个人媒体
func (s *DefaultPersonalMediaService) GetByID(ctx context.Context, id int64, userID int64) (*models.PersonalMedia, error) {
	media, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("媒体不存在")
	}

	// 检查是否属于当前用户
	if media.UserID != userID {
		return nil, errors.New("无权访问此媒体")
	}

	return media, nil
}

// Update 更新个人媒体
func (s *DefaultPersonalMediaService) Update(ctx context.Context, id int64, userID int64, title string, description []byte, isPrivate bool, tags []string) (*models.PersonalMedia, error) {
	// 获取媒体
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
func (s *DefaultPersonalMediaService) Delete(ctx context.Context, id int64, userID int64) error {
	// 获取媒体（检查权限）
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// 删除记录
	return s.repo.Delete(ctx, id)
}

// MigrateExistingMedia 迁移现有媒体数据，添加路径信息
func (s *DefaultPersonalMediaService) MigrateExistingMedia(ctx context.Context) error {
	return s.repo.MigrateExistingMedia(ctx)
}

package service

import (
	"context"
	"errors"
	"fmt"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

var (
	ErrPhotoVideoNotFound = errors.New("照片/视频不存在")
)

// PhotoVideoService 照片和视频服务接口
type PhotoVideoService interface {
	Service
	CreatePhotoVideo(ctx context.Context, photoVideo *models.PhotoVideo) (*models.PhotoVideo, error)
	GetPhotoVideoByID(ctx context.Context, id int64) (*models.PhotoVideo, error)
	ListPhotoVideosByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.PhotoVideo, int64, error)
	ListPhotoVideosByCategory(ctx context.Context, coupleID int64, category string, offset, limit int) ([]*models.PhotoVideo, int64, error)
	ListPhotoVideosByMediaType(ctx context.Context, coupleID int64, mediaType string, offset, limit int) ([]*models.PhotoVideo, int64, error)
	ListPhotoVideosByEventID(ctx context.Context, eventID int64) ([]*models.PhotoVideo, error)
	ListPhotoVideosByLocationID(ctx context.Context, locationID int64) ([]*models.PhotoVideo, error)
	UpdatePhotoVideo(ctx context.Context, photoVideo *models.PhotoVideo) error
	DeletePhotoVideo(ctx context.Context, id int64) error
}

// photoVideoService 照片和视频服务实现
type photoVideoService struct {
	*BaseService
	photoVideoRepo repository.PhotoVideoRepository
}

// NewPhotoVideoService 创建照片和视频服务
func NewPhotoVideoService(photoVideoRepo repository.PhotoVideoRepository) PhotoVideoService {
	return &photoVideoService{
		BaseService:    NewBaseService(photoVideoRepo),
		photoVideoRepo: photoVideoRepo,
	}
}

// CreatePhotoVideo 创建照片/视频
func (s *photoVideoService) CreatePhotoVideo(ctx context.Context, photoVideo *models.PhotoVideo) (*models.PhotoVideo, error) {
	if err := s.photoVideoRepo.Create(ctx, photoVideo); err != nil {
		return nil, fmt.Errorf("创建照片/视频失败: %w", err)
	}
	return photoVideo, nil
}

// GetPhotoVideoByID 通过ID获取照片/视频
func (s *photoVideoService) GetPhotoVideoByID(ctx context.Context, id int64) (*models.PhotoVideo, error) {
	photoVideo, err := s.photoVideoRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPhotoVideoNotFound) {
			return nil, ErrPhotoVideoNotFound
		}
		return nil, fmt.Errorf("获取照片/视频失败: %w", err)
	}
	return photoVideo, nil
}

// ListPhotoVideosByCoupleID 获取情侣关系下的所有照片/视频
func (s *photoVideoService) ListPhotoVideosByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	return s.photoVideoRepo.ListByCoupleID(ctx, coupleID, offset, limit)
}

// ListPhotoVideosByCategory 按分类获取照片/视频
func (s *photoVideoService) ListPhotoVideosByCategory(ctx context.Context, coupleID int64, category string, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	return s.photoVideoRepo.ListByCategory(ctx, coupleID, category, offset, limit)
}

// ListPhotoVideosByMediaType 按媒体类型获取照片/视频
func (s *photoVideoService) ListPhotoVideosByMediaType(ctx context.Context, coupleID int64, mediaType string, offset, limit int) ([]*models.PhotoVideo, int64, error) {
	return s.photoVideoRepo.ListByMediaType(ctx, coupleID, mediaType, offset, limit)
}

// ListPhotoVideosByEventID 按事件ID获取照片/视频
func (s *photoVideoService) ListPhotoVideosByEventID(ctx context.Context, eventID int64) ([]*models.PhotoVideo, error) {
	return s.photoVideoRepo.ListByEventID(ctx, eventID)
}

// ListPhotoVideosByLocationID 按地点ID获取照片/视频
func (s *photoVideoService) ListPhotoVideosByLocationID(ctx context.Context, locationID int64) ([]*models.PhotoVideo, error) {
	return s.photoVideoRepo.ListByLocationID(ctx, locationID)
}

// UpdatePhotoVideo 更新照片/视频
func (s *photoVideoService) UpdatePhotoVideo(ctx context.Context, photoVideo *models.PhotoVideo) error {
	// 检查照片/视频是否存在
	_, err := s.photoVideoRepo.GetByID(ctx, photoVideo.ID)
	if err != nil {
		if errors.Is(err, repository.ErrPhotoVideoNotFound) {
			return ErrPhotoVideoNotFound
		}
		return fmt.Errorf("查询照片/视频失败: %w", err)
	}

	// 更新照片/视频
	if err := s.photoVideoRepo.Update(ctx, photoVideo); err != nil {
		return fmt.Errorf("更新照片/视频失败: %w", err)
	}
	return nil
}

// DeletePhotoVideo 删除照片/视频
func (s *photoVideoService) DeletePhotoVideo(ctx context.Context, id int64) error {
	// 检查照片/视频是否存在
	_, err := s.photoVideoRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPhotoVideoNotFound) {
			return ErrPhotoVideoNotFound
		}
		return fmt.Errorf("查询照片/视频失败: %w", err)
	}

	// 删除照片/视频
	if err := s.photoVideoRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除照片/视频失败: %w", err)
	}
	return nil
}

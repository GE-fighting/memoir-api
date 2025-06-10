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
	ErrPhotoVideoNotFound = errors.New("照片/视频不存在")
)

// PhotoVideoService 照片和视频服务接口
type PhotoVideoService interface {
	Service
	CreatePhotoVideo(ctx context.Context, dto *dto.CreatePhotoVideoRequest) (*models.PhotoVideo, error)
	GetPhotoVideoByID(ctx context.Context, id int64) (*models.PhotoVideo, error)
	Query(ctx context.Context, params *dto.PhotoVideoQueryParams) (*dto.PageResult, error)
	UpdatePhotoVideo(ctx context.Context, photoVideo *models.PhotoVideo) error
	DeletePhotoVideo(ctx context.Context, id int64) error
}

// photoVideoService 照片和视频服务实现
type photoVideoService struct {
	*BaseService
	photoVideoRepo repository.PhotoVideoRepository
	userRepo       repository.UserRepository
	ablumRepo      repository.CoupleAlbumRepository
}

func (s *photoVideoService) Query(ctx context.Context, params *dto.PhotoVideoQueryParams) (*dto.PageResult, error) {
	data, total, err := s.photoVideoRepo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	result := dto.NewPageResult(data, total, params.Page, params.PageSize)
	return &result, nil
}

// NewPhotoVideoService 创建照片和视频服务
func NewPhotoVideoService(photoVideoRepo repository.PhotoVideoRepository, userRepo repository.UserRepository, ablumRepo repository.CoupleAlbumRepository) PhotoVideoService {
	return &photoVideoService{
		BaseService:    NewBaseService(photoVideoRepo),
		photoVideoRepo: photoVideoRepo,
		userRepo:       userRepo,
		ablumRepo:      ablumRepo,
	}
}

// CreatePhotoVideo 创建照片/视频
func (s *photoVideoService) CreatePhotoVideo(ctx context.Context, dto *dto.CreatePhotoVideoRequest) (*models.PhotoVideo, error) {
	photoVideo := dto.ToModel()
	user, err := s.userRepo.GetByID(ctx, dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user.CoupleID == 0 {
		return nil, fmt.Errorf("用户没有情侣关系")
	}
	photoVideo.CoupleID = user.CoupleID
	if err := s.photoVideoRepo.Create(ctx, photoVideo); err != nil {
		return nil, fmt.Errorf("创建照片/视频失败: %w", err)
	}
	//对应相册 照片数量+1
	ablum, err := s.ablumRepo.GetByID(ctx, photoVideo.AlbumID)
	if err != nil {
		return nil, fmt.Errorf("查询相册失败：%w", err)
	}
	ablum.Count = ablum.Count + 1
	err = s.ablumRepo.Update(ctx, ablum)
	if err != nil {
		return nil, fmt.Errorf("更新相册中媒体数量失败：%w", err)
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

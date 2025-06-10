package service

import (
	"context"
	"errors"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

// CoupleAlbumService 情侣相册服务接口
type CoupleAlbumService interface {
	Service
	Create(ctx context.Context, req *dto.CreateCoupleAlbumRequest) (*models.CoupleAlbum, error)
	GetByID(ctx context.Context, id int64) (*models.CoupleAlbum, error)
	GetByCoupleID(ctx context.Context, coupleID int64) ([]*models.CoupleAlbum, error)
	Update(ctx context.Context, id int64, req *dto.UpdateCoupleAlbumRequest) (*models.CoupleAlbum, error)
	Delete(ctx context.Context, id int64) error
	GetWithPhotos(ctx context.Context, id int64) (*models.CoupleAlbum, error)
}

// coupleAlbumService 情侣相册服务实现
type coupleAlbumService struct {
	*BaseService
	coupleAlbumRepo repository.CoupleAlbumRepository
	userService     UserService
}

// NewCoupleAlbumService 创建情侣相册服务
func NewCoupleAlbumService(
	coupleAlbumRepo repository.CoupleAlbumRepository,
	userService UserService,
) CoupleAlbumService {
	return &coupleAlbumService{
		BaseService:     NewBaseService(coupleAlbumRepo),
		coupleAlbumRepo: coupleAlbumRepo,
		userService:     userService,
	}
}

// Create 创建情侣相册
func (s *coupleAlbumService) Create(ctx context.Context, req *dto.CreateCoupleAlbumRequest) (*models.CoupleAlbum, error) {
	// 获取用户信息
	user, err := s.userService.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 确保用户属于一个情侣关系
	if user.CoupleID == 0 {
		return nil, errors.New("用户不属于任何情侣关系")
	}

	// 创建相册
	var description string
	if req.Description != "" {
		description = req.Description
	}

	var coverURL *string
	if req.CoverURL != nil {
		coverURL = req.CoverURL
	}

	album := &models.CoupleAlbum{
		CoupleID:    user.CoupleID,
		Title:       req.Title,
		Description: description,
		CoverURL:    coverURL,
	}

	if err := s.coupleAlbumRepo.Create(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

// GetByID 通过ID获取情侣相册
func (s *coupleAlbumService) GetByID(ctx context.Context, id int64) (*models.CoupleAlbum, error) {
	return s.coupleAlbumRepo.GetByID(ctx, id)
}

// GetByCoupleID 通过情侣ID获取所有相册
func (s *coupleAlbumService) GetByCoupleID(ctx context.Context, coupleID int64) ([]*models.CoupleAlbum, error) {
	return s.coupleAlbumRepo.GetByCoupleID(ctx, coupleID)
}

// Update 更新情侣相册
func (s *coupleAlbumService) Update(ctx context.Context, id int64, req *dto.UpdateCoupleAlbumRequest) (*models.CoupleAlbum, error) {
	album, err := s.coupleAlbumRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新相册信息
	if req.Title != "" {
		album.Title = req.Title
	}

	if req.Description != "" {
		album.Description = req.Description
	}

	if err := s.coupleAlbumRepo.Update(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

// Delete 删除情侣相册
func (s *coupleAlbumService) Delete(ctx context.Context, id int64) error {
	return s.coupleAlbumRepo.Delete(ctx, id)
}

// GetWithPhotos 获取相册及其包含的照片和视频
func (s *coupleAlbumService) GetWithPhotos(ctx context.Context, id int64) (*models.CoupleAlbum, error) {
	return s.coupleAlbumRepo.GetWithPhotos(ctx, id)
}

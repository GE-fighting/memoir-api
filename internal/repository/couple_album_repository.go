package repository

import (
	"context"
	"memoir-api/internal/models"
)

type CoupleAlbumRepository interface {
	Repository
	Create(ctx context.Context, album *models.CoupleAlbum) error
	GetByID(ctx context.Context, id int64) (*models.CoupleAlbum, error)
	GetByCoupleID(ctx context.Context, coupleID int64) ([]*models.CoupleAlbum, error)
	Update(ctx context.Context, album *models.CoupleAlbum) error
	Delete(ctx context.Context, id int64) error
	GetWithPhotos(ctx context.Context, id int64) (*models.CoupleAlbum, error)
}

type coupleAlbumRepository struct {
	*BaseRepository
}

func NewCoupleAlbumRepository(db *BaseRepository) CoupleAlbumRepository {
	return &coupleAlbumRepository{db}
}

func (r *coupleAlbumRepository) Create(ctx context.Context, album *models.CoupleAlbum) error {
	return r.DB().WithContext(ctx).Create(album).Error
}

func (r *coupleAlbumRepository) GetByID(ctx context.Context, id int64) (*models.CoupleAlbum, error) {
	var album models.CoupleAlbum
	err := r.DB().WithContext(ctx).First(&album, id).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (r *coupleAlbumRepository) GetByCoupleID(ctx context.Context, coupleID int64) ([]*models.CoupleAlbum, error) {
	var albums []*models.CoupleAlbum
	err := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (r *coupleAlbumRepository) Update(ctx context.Context, album *models.CoupleAlbum) error {
	return r.DB().WithContext(ctx).Save(album).Error
}

func (r *coupleAlbumRepository) Delete(ctx context.Context, id int64) error {
	return r.DB().WithContext(ctx).Delete(&models.CoupleAlbum{}, id).Error
}

func (r *coupleAlbumRepository) GetWithPhotos(ctx context.Context, id int64) (*models.CoupleAlbum, error) {
	var album models.CoupleAlbum
	err := r.DB().WithContext(ctx).First(&album, id).Error
	if err != nil {
		return nil, err
	}

	// 查询相关的照片和视频
	var photos []models.PhotoVideo
	err = r.DB().WithContext(ctx).Where("album_id = ?", id).Find(&photos).Error
	if err != nil {
		return nil, err
	}

	album.PhotosVideos = photos
	return &album, nil
}

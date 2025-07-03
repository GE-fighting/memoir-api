package dto

import (
	"memoir-api/internal/models"
	"time"
)

type CreateCoupleAlbumRequest struct {
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	CoverURL    *string `json:"cover_url,omitempty"`
}

type UpdateCoupleAlbumRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CoverURL    *string `json:"cover_url,omitempty"`
}

type DeleteCoupleAlbumPhotosRequest struct {
	AlbumID       int64      `json:"album_id,string" binding:"required"`
	PhotoVideoIDs Int64Array `json:"photo_video_ids" binding:"required"`
}

type CoupleAlbumDTO struct {
	ID          int64           `json:"id"`
	CoupleID    int64           `json:"couple_id,string"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CoverURL    *string         `json:"cover_url,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	PhotoCount  int             `json:"photo_size"`
	Photos      []PhotoVideoDTO `json:"photos_videos,omitempty"`
}

type CoupleAlbumQueryParams struct {
	PaginationRequest
	CoupleID  int64  `form:"couple_id,string"`
	Title     string `form:"title"`
	MediaType string `form:"media_type"`
}

// ToModel 将DTO转换为模型对象
func (r *CreateCoupleAlbumRequest) ToModel(coupleID int64) *models.CoupleAlbum {
	return &models.CoupleAlbum{
		CoupleID:    coupleID,
		Title:       r.Title,
		Description: r.Description,
		CoverURL:    r.CoverURL,
	}
}

// FromModel 从模型创建DTO
func CoupleAlbumFromModel(album *models.CoupleAlbum) CoupleAlbumDTO {
	dto := CoupleAlbumDTO{
		ID:          album.ID,
		CoupleID:    album.CoupleID,
		Title:       album.Title,
		Description: album.Description,
		CoverURL:    album.CoverURL,
		CreatedAt:   album.CreatedAt,
		UpdatedAt:   album.UpdatedAt,
		PhotoCount:  album.Count,
	}

	// 转换照片视频列表
	if len(album.PhotosVideos) > 0 {
		dto.Photos = make([]PhotoVideoDTO, 0, len(album.PhotosVideos))
		for _, pv := range album.PhotosVideos {
			photoDTO := PhotoVideoDTO{
				ID:           pv.ID,
				MediaType:    pv.MediaType,
				MediaURL:     pv.MediaURL,
				ThumbnailURL: pv.ThumbnailURL,
				Title:        pv.Title,
				Description:  pv.Description,
				CreatedAt:    pv.CreatedAt,
			}
			dto.Photos = append(dto.Photos, photoDTO)
		}
	}

	return dto
}

// PhotoVideoDTO 简化的照片视频DTO，用于相册列表
type PhotoVideoDTO struct {
	ID           int64     `json:"id"`
	MediaType    string    `json:"media_type"`
	MediaURL     string    `json:"media_url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

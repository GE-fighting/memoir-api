package dto

import "memoir-api/internal/models"

type CreatePhotoVideoRequest struct {
	UserID       int64  `json:"user_id"`
	MediaType    string `json:"media_type" binding:"required,oneof=photo video"`
	Title        string `json:"title"`
	MediaURL     string `json:"media_url" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	Description  string `json:"description"`
	CoupleID     int64  `json:"couple_id"`
	AlbumID      int64  `json:"album_id,string" binding:"required"`
	EventID      *int64 `json:"event_id,string,omitempty"`
	LocationID   *int64 `json:"location_id,string,omitempty"`
}

type PhotoVideoQueryParams struct {
	PaginationRequest
	CoupleID   int64  `form:"couple_id,string"`
	AlbumID    int64  `form:"album_id,string"`
	MediaType  string `form:"media_type" binding:"omitempty,oneof=photo video"`
	EventID    int64  `form:"event_id,string"`
	LocationID int64  `form:"location_id,string"`
}

func (r *CreatePhotoVideoRequest) ToModel() *models.PhotoVideo {
	return &models.PhotoVideo{
		MediaType:    r.MediaType,
		Title:        r.Title,
		MediaURL:     r.MediaURL,
		ThumbnailURL: r.ThumbnailURL,
		Description:  r.Description,
		CoupleID:     r.CoupleID,
		AlbumID:      r.AlbumID,
		EventID:      r.EventID,
		LocationID:   r.LocationID,
	}
}

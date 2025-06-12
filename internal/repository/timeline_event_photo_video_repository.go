package repository

import (
	"context"
	"errors"
	"memoir-api/internal/models"

	"gorm.io/gorm"
)

// TimelineEventPhotoVideoRepository 时间线事件与照片/视频关联仓库接口
type TimelineEventPhotoVideoRepository interface {
	Repository
	Create(ctx context.Context, eventPhotoVideo *models.TimelineEventPhotoVideo) error
	FindByEventID(ctx context.Context, eventID int64) ([]models.TimelineEventPhotoVideo, error)
	FindByPhotoVideoID(ctx context.Context, photoVideoID int64) ([]models.TimelineEventPhotoVideo, error)
	DeleteByEventID(ctx context.Context, eventID int64) error
	DeleteByPhotoVideoID(ctx context.Context, photoVideoID int64) error
}

// ErrTimelineEventPhotoVideoNotFound 时间线事件与照片/视频关联不存在错误
var ErrTimelineEventPhotoVideoNotFound = errors.New("timeline event photo video not found")

// timelineEventPhotoVideoRepository 时间线事件与照片/视频关联仓库实现
type timelineEventPhotoVideoRepository struct {
	*BaseRepository
}

// NewTimelineEventPhotoVideoRepository 创建时间线事件与照片/视频关联仓库
func NewTimelineEventPhotoVideoRepository(db *gorm.DB) TimelineEventPhotoVideoRepository {
	return &timelineEventPhotoVideoRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建时间线事件与照片/视频关联
func (r *timelineEventPhotoVideoRepository) Create(ctx context.Context, eventPhotoVideo *models.TimelineEventPhotoVideo) error {
	return r.DB().WithContext(ctx).Create(eventPhotoVideo).Error
}

// FindByEventID 根据事件ID查询关联
func (r *timelineEventPhotoVideoRepository) FindByEventID(ctx context.Context, eventID int64) ([]models.TimelineEventPhotoVideo, error) {
	var eventPhotosVideos []models.TimelineEventPhotoVideo
	err := r.DB().WithContext(ctx).Where("timeline_event_id = ?", eventID).Find(&eventPhotosVideos).Error
	if err != nil {
		return nil, err
	}
	return eventPhotosVideos, nil
}

// FindByPhotoVideoID 根据照片/视频ID查询关联
func (r *timelineEventPhotoVideoRepository) FindByPhotoVideoID(ctx context.Context, photoVideoID int64) ([]models.TimelineEventPhotoVideo, error) {
	var eventPhotosVideos []models.TimelineEventPhotoVideo
	err := r.DB().WithContext(ctx).Where("photo_video_id = ?", photoVideoID).Find(&eventPhotosVideos).Error
	if err != nil {
		return nil, err
	}
	return eventPhotosVideos, nil
}

// DeleteByEventID 根据事件ID删除关联
func (r *timelineEventPhotoVideoRepository) DeleteByEventID(ctx context.Context, eventID int64) error {
	return r.DB().WithContext(ctx).Where("timeline_event_id = ?", eventID).Delete(&models.TimelineEventPhotoVideo{}).Error
}

// DeleteByPhotoVideoID 根据照片/视频ID删除关联
func (r *timelineEventPhotoVideoRepository) DeleteByPhotoVideoID(ctx context.Context, photoVideoID int64) error {
	return r.DB().WithContext(ctx).Where("photo_video_id = ?", photoVideoID).Delete(&models.TimelineEventPhotoVideo{}).Error
}

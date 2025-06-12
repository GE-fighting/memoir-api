package repository

import (
	"context"
	"errors"
	"memoir-api/internal/models"

	"gorm.io/gorm"
)

// TimelineEventLocationRepository 时间线事件与位置关联仓库接口
type TimelineEventLocationRepository interface {
	Repository
	Create(ctx context.Context, eventLocation *models.TimelineEventLocation) error
	FindByEventID(ctx context.Context, eventID int64) ([]models.TimelineEventLocation, error)
	FindByLocationID(ctx context.Context, locationID int64) ([]models.TimelineEventLocation, error)
	DeleteByEventID(ctx context.Context, eventID int64) error
	DeleteByLocationID(ctx context.Context, locationID int64) error
}

// ErrTimelineEventLocationNotFound 时间线事件与位置关联不存在错误
var ErrTimelineEventLocationNotFound = errors.New("timeline event location not found")

// timelineEventLocationRepository 时间线事件与位置关联仓库实现
type timelineEventLocationRepository struct {
	*BaseRepository
}

// NewTimelineEventLocationRepository 创建时间线事件与位置关联仓库
func NewTimelineEventLocationRepository(db *gorm.DB) TimelineEventLocationRepository {
	return &timelineEventLocationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建时间线事件与位置关联
func (r *timelineEventLocationRepository) Create(ctx context.Context, eventLocation *models.TimelineEventLocation) error {
	return r.DB().WithContext(ctx).Create(eventLocation).Error
}

// FindByEventID 根据事件ID查询关联
func (r *timelineEventLocationRepository) FindByEventID(ctx context.Context, eventID int64) ([]models.TimelineEventLocation, error) {
	var eventLocations []models.TimelineEventLocation
	err := r.DB().WithContext(ctx).Where("timeline_event_id = ?", eventID).Find(&eventLocations).Error
	if err != nil {
		return nil, err
	}
	return eventLocations, nil
}

// FindByLocationID 根据位置ID查询关联
func (r *timelineEventLocationRepository) FindByLocationID(ctx context.Context, locationID int64) ([]models.TimelineEventLocation, error) {
	var eventLocations []models.TimelineEventLocation
	err := r.DB().WithContext(ctx).Where("location_id = ?", locationID).Find(&eventLocations).Error
	if err != nil {
		return nil, err
	}
	return eventLocations, nil
}

// DeleteByEventID 根据事件ID删除关联
func (r *timelineEventLocationRepository) DeleteByEventID(ctx context.Context, eventID int64) error {
	return r.DB().WithContext(ctx).Where("timeline_event_id = ?", eventID).Delete(&models.TimelineEventLocation{}).Error
}

// DeleteByLocationID 根据位置ID删除关联
func (r *timelineEventLocationRepository) DeleteByLocationID(ctx context.Context, locationID int64) error {
	return r.DB().WithContext(ctx).Where("location_id = ?", locationID).Delete(&models.TimelineEventLocation{}).Error
}

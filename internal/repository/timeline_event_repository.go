package repository

import (
	"context"
	"errors"

	"memoir-api/internal/models"

	"gorm.io/gorm"
)

var (
	ErrTimelineEventNotFound = errors.New("时间轴事件不存在")
)

// TimelineEventRepository 时间轴事件仓库接口
type TimelineEventRepository interface {
	Repository
	Create(ctx context.Context, event *models.TimelineEvent) error
	FindByID(ctx context.Context, id int64) (*models.TimelineEvent, error)
	FindWithPagination(ctx context.Context, conditions map[string]interface{}, offset, limit int) ([]models.TimelineEvent, int64, error)
	Update(ctx context.Context, event *models.TimelineEvent) error
	Delete(ctx context.Context, id int64) error
	FindByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error)

	// 关联查询方法
	FindWithLocationsAndPhotos(ctx context.Context, id int64) (*models.TimelineEvent, error)
	FindWithLocationsAndPhotosByJoins(ctx context.Context, id int64) (*models.TimelineEvent, error)
	FindEventsByLocationID(ctx context.Context, locationID int64) ([]models.TimelineEvent, error)
	FindEventsByPhotoVideoID(ctx context.Context, photoVideoID int64) ([]models.TimelineEvent, error)
	FindEventsWithFilters(ctx context.Context, coupleID int64, startDate, endDate *string, title *string, locationID *int64, offset, limit int) ([]models.TimelineEvent, int64, error)
}

// timelineEventRepository 时间轴事件仓库实现
type timelineEventRepository struct {
	*BaseRepository
}

// NewTimelineEventRepository 创建时间轴事件仓库
func NewTimelineEventRepository(db *gorm.DB) TimelineEventRepository {
	return &timelineEventRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建时间轴事件
func (r *timelineEventRepository) Create(ctx context.Context, event *models.TimelineEvent) error {
	return r.DB().WithContext(ctx).Create(event).Error
}

// FindByID 根据ID查询时间轴事件
func (r *timelineEventRepository) FindByID(ctx context.Context, id int64) (*models.TimelineEvent, error) {
	var event models.TimelineEvent
	err := r.DB().WithContext(ctx).First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimelineEventNotFound
		}
		return nil, err
	}
	return &event, nil
}

// FindWithPagination 分页查询时间轴事件
func (r *timelineEventRepository) FindWithPagination(ctx context.Context, conditions map[string]interface{}, offset, limit int) ([]models.TimelineEvent, int64, error) {
	var events []models.TimelineEvent
	var total int64

	query := r.DB().WithContext(ctx).Model(&models.TimelineEvent{})

	// 应用查询条件
	for key, value := range conditions {
		if key == "id IN ?" || key == "event_date >= ?" || key == "event_date <= ?" || key == "title LIKE ?" {
			query = query.Where(key, value)
		} else {
			query = query.Where(key, value)
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// Update 更新时间轴事件
func (r *timelineEventRepository) Update(ctx context.Context, event *models.TimelineEvent) error {
	return r.DB().WithContext(ctx).Save(event).Error
}

// Delete 删除时间轴事件
func (r *timelineEventRepository) Delete(ctx context.Context, id int64) error {
	return r.DB().WithContext(ctx).Delete(&models.TimelineEvent{}, id).Error
}

// FindByCoupleID 根据CoupleID查询时间轴事件
func (r *timelineEventRepository) FindByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error) {
	var events []*models.TimelineEvent
	var total int64

	if err := r.DB().WithContext(ctx).Model(&models.TimelineEvent{}).Where("couple_id = ?", coupleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID).Offset(offset).Limit(limit).Order("end_date DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// FindWithLocationsAndPhotos 查询时间轴事件并加载关联的地点和照片
func (r *timelineEventRepository) FindWithLocationsAndPhotos(ctx context.Context, id int64) (*models.TimelineEvent, error) {
	var event models.TimelineEvent

	// 方法1：使用Preload预加载关联数据（会产生多条SQL查询）
	err := r.DB().WithContext(ctx).
		Preload("TimelineEventLocations").
		Preload("TimelineEventPhotosVideos").
		First(&event, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimelineEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

// FindWithLocationsAndPhotosByJoins 使用Joins方式查询时间轴事件及其关联的地点和照片
func (r *timelineEventRepository) FindWithLocationsAndPhotosByJoins(ctx context.Context, id int64) (*models.TimelineEvent, error) {
	var event models.TimelineEvent

	// 方法2：使用Joins进行关联查询（生成单条SQL查询，但对于一对多关系需要手动处理）
	// 注意：这种方式只适合一对一关系，对于一对多关系可能会导致数据重复
	err := r.DB().WithContext(ctx).
		Joins("LEFT JOIN timeline_event_locations ON timeline_events.id = timeline_event_locations.timeline_event_id").
		Joins("LEFT JOIN locations ON timeline_event_locations.location_id = locations.id").
		Joins("LEFT JOIN timeline_event_photo_videos ON timeline_events.id = timeline_event_photo_videos.timeline_event_id").
		Joins("LEFT JOIN photo_videos ON timeline_event_photo_videos.photo_video_id = photo_videos.id").
		Where("timeline_events.id = ?", id).
		First(&event).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimelineEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

// FindEventsByLocationID 根据地点ID查询关联的时间轴事件
func (r *timelineEventRepository) FindEventsByLocationID(ctx context.Context, locationID int64) ([]models.TimelineEvent, error) {
	var events []models.TimelineEvent

	// 使用子查询方式，先查询关联表获取事件ID，再查询事件
	subQuery := r.DB().WithContext(ctx).
		Table("timeline_event_locations").
		Select("timeline_event_id").
		Where("location_id = ?", locationID)

	err := r.DB().WithContext(ctx).
		Where("id IN (?)", subQuery).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

// FindEventsByPhotoVideoID 根据照片/视频ID查询关联的时间轴事件
func (r *timelineEventRepository) FindEventsByPhotoVideoID(ctx context.Context, photoVideoID int64) ([]models.TimelineEvent, error) {
	var events []models.TimelineEvent

	// 使用子查询方式，先查询关联表获取事件ID，再查询事件
	subQuery := r.DB().WithContext(ctx).
		Table("timeline_event_photo_videos").
		Select("timeline_event_id").
		Where("photo_video_id = ?", photoVideoID)

	err := r.DB().WithContext(ctx).
		Where("id IN (?)", subQuery).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

// FindEventsWithFilters 根据多条件查询时间轴事件
func (r *timelineEventRepository) FindEventsWithFilters(ctx context.Context, coupleID int64, startDate, endDate *string, title *string, locationID *int64, offset, limit int) ([]models.TimelineEvent, int64, error) {
	var events []models.TimelineEvent
	var total int64

	query := r.DB().WithContext(ctx).Model(&models.TimelineEvent{}).Where("couple_id = ?", coupleID)

	// 添加日期范围条件
	if startDate != nil {
		query = query.Where("event_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("event_date <= ?", *endDate)
	}

	// 添加标题模糊查询条件
	if title != nil && *title != "" {
		query = query.Where("title LIKE ?", "%"+*title+"%")
	}

	// 添加地点过滤条件
	if locationID != nil && *locationID > 0 {
		// 使用子查询方式，先查询关联表获取事件ID，再过滤事件
		subQuery := r.DB().WithContext(ctx).
			Table("timeline_event_locations").
			Select("timeline_event_id").
			Where("location_id = ?", *locationID)

		query = query.Where("id IN (?)", subQuery)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

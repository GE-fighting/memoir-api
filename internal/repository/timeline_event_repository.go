package repository

import (
	"context"
	"errors"
	"time"

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
	GetByID(ctx context.Context, id int64) (*models.TimelineEvent, error)
	ListByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error)
	ListByDateRange(ctx context.Context, coupleID int64, startDate, endDate time.Time) ([]*models.TimelineEvent, error)
	ListByYear(ctx context.Context, coupleID int64, year int) ([]*models.TimelineEvent, error)
	ListByMonth(ctx context.Context, coupleID int64, year, month int) ([]*models.TimelineEvent, error)
	ListByLocationID(ctx context.Context, locationID int64) ([]*models.TimelineEvent, error)
	Update(ctx context.Context, event *models.TimelineEvent) error
	Delete(ctx context.Context, id int64) error
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

// GetByID 通过ID获取时间轴事件
func (r *timelineEventRepository) GetByID(ctx context.Context, id int64) (*models.TimelineEvent, error) {
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

// ListByCoupleID 获取情侣关系下的所有时间轴事件，按事件日期倒序排列
func (r *timelineEventRepository) ListByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error) {
	var events []*models.TimelineEvent
	var total int64

	// 获取总数
	if err := r.DB().WithContext(ctx).Model(&models.TimelineEvent{}).Where("couple_id = ?", coupleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.DB().WithContext(ctx).Where("couple_id = ?", coupleID).Order("event_date DESC")
	if offset >= 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// ListByDateRange 按日期范围获取时间轴事件
func (r *timelineEventRepository) ListByDateRange(ctx context.Context, coupleID int64, startDate, endDate time.Time) ([]*models.TimelineEvent, error) {
	var events []*models.TimelineEvent
	err := r.DB().WithContext(ctx).
		Where("couple_id = ? AND event_date BETWEEN ? AND ?", coupleID, startDate, endDate).
		Order("event_date DESC").
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListByYear 按年份获取时间轴事件
func (r *timelineEventRepository) ListByYear(ctx context.Context, coupleID int64, year int) ([]*models.TimelineEvent, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)
	return r.ListByDateRange(ctx, coupleID, startDate, endDate)
}

// ListByMonth 按年月获取时间轴事件
func (r *timelineEventRepository) ListByMonth(ctx context.Context, coupleID int64, year, month int) ([]*models.TimelineEvent, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // 下个月第一天减一天
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
	return r.ListByDateRange(ctx, coupleID, startDate, endDate)
}

// ListByLocationID 按地点ID获取时间轴事件
func (r *timelineEventRepository) ListByLocationID(ctx context.Context, locationID int64) ([]*models.TimelineEvent, error) {
	var events []*models.TimelineEvent
	err := r.DB().WithContext(ctx).
		Where("location_id = ?", locationID).
		Order("event_date DESC").
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

// Update 更新时间轴事件
func (r *timelineEventRepository) Update(ctx context.Context, event *models.TimelineEvent) error {
	result := r.DB().WithContext(ctx).Save(event)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTimelineEventNotFound
	}
	return nil
}

// Delete 删除时间轴事件
func (r *timelineEventRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB().WithContext(ctx).Delete(&models.TimelineEvent{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTimelineEventNotFound
	}
	return nil
}

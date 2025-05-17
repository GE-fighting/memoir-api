package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

var (
	ErrTimelineEventNotFound = errors.New("时间轴事件不存在")
)

// TimelineEventService 时间轴事件服务接口
type TimelineEventService interface {
	Service
	CreateTimelineEvent(ctx context.Context, event *models.TimelineEvent) (*models.TimelineEvent, error)
	GetTimelineEventByID(ctx context.Context, id int64) (*models.TimelineEvent, error)
	ListTimelineEventsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error)
	ListTimelineEventsByDateRange(ctx context.Context, coupleID int64, startDate, endDate time.Time) ([]*models.TimelineEvent, error)
	ListTimelineEventsByYear(ctx context.Context, coupleID int64, year int) ([]*models.TimelineEvent, error)
	ListTimelineEventsByMonth(ctx context.Context, coupleID int64, year, month int) ([]*models.TimelineEvent, error)
	ListTimelineEventsByLocationID(ctx context.Context, locationID int64) ([]*models.TimelineEvent, error)
	UpdateTimelineEvent(ctx context.Context, event *models.TimelineEvent) error
	DeleteTimelineEvent(ctx context.Context, id int64) error
}

// timelineEventService 时间轴事件服务实现
type timelineEventService struct {
	*BaseService
	timelineEventRepo repository.TimelineEventRepository
}

// NewTimelineEventService 创建时间轴事件服务
func NewTimelineEventService(timelineEventRepo repository.TimelineEventRepository) TimelineEventService {
	return &timelineEventService{
		BaseService:       NewBaseService(timelineEventRepo),
		timelineEventRepo: timelineEventRepo,
	}
}

// CreateTimelineEvent 创建时间轴事件
func (s *timelineEventService) CreateTimelineEvent(ctx context.Context, event *models.TimelineEvent) (*models.TimelineEvent, error) {
	if err := s.timelineEventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("创建时间轴事件失败: %w", err)
	}
	return event, nil
}

// GetTimelineEventByID 通过ID获取时间轴事件
func (s *timelineEventService) GetTimelineEventByID(ctx context.Context, id int64) (*models.TimelineEvent, error) {
	event, err := s.timelineEventRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTimelineEventNotFound) {
			return nil, ErrTimelineEventNotFound
		}
		return nil, fmt.Errorf("获取时间轴事件失败: %w", err)
	}
	return event, nil
}

// ListTimelineEventsByCoupleID 获取情侣关系下的所有时间轴事件
func (s *timelineEventService) ListTimelineEventsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error) {
	return s.timelineEventRepo.ListByCoupleID(ctx, coupleID, offset, limit)
}

// ListTimelineEventsByDateRange 按日期范围获取时间轴事件
func (s *timelineEventService) ListTimelineEventsByDateRange(ctx context.Context, coupleID int64, startDate, endDate time.Time) ([]*models.TimelineEvent, error) {
	return s.timelineEventRepo.ListByDateRange(ctx, coupleID, startDate, endDate)
}

// ListTimelineEventsByYear 按年份获取时间轴事件
func (s *timelineEventService) ListTimelineEventsByYear(ctx context.Context, coupleID int64, year int) ([]*models.TimelineEvent, error) {
	return s.timelineEventRepo.ListByYear(ctx, coupleID, year)
}

// ListTimelineEventsByMonth 按年月获取时间轴事件
func (s *timelineEventService) ListTimelineEventsByMonth(ctx context.Context, coupleID int64, year, month int) ([]*models.TimelineEvent, error) {
	return s.timelineEventRepo.ListByMonth(ctx, coupleID, year, month)
}

// ListTimelineEventsByLocationID 按地点ID获取时间轴事件
func (s *timelineEventService) ListTimelineEventsByLocationID(ctx context.Context, locationID int64) ([]*models.TimelineEvent, error) {
	return s.timelineEventRepo.ListByLocationID(ctx, locationID)
}

// UpdateTimelineEvent 更新时间轴事件
func (s *timelineEventService) UpdateTimelineEvent(ctx context.Context, event *models.TimelineEvent) error {
	// 检查时间轴事件是否存在
	_, err := s.timelineEventRepo.GetByID(ctx, event.ID)
	if err != nil {
		if errors.Is(err, repository.ErrTimelineEventNotFound) {
			return ErrTimelineEventNotFound
		}
		return fmt.Errorf("查询时间轴事件失败: %w", err)
	}

	// 更新时间轴事件
	if err := s.timelineEventRepo.Update(ctx, event); err != nil {
		return fmt.Errorf("更新时间轴事件失败: %w", err)
	}
	return nil
}

// DeleteTimelineEvent 删除时间轴事件
func (s *timelineEventService) DeleteTimelineEvent(ctx context.Context, id int64) error {
	// 检查时间轴事件是否存在
	_, err := s.timelineEventRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTimelineEventNotFound) {
			return ErrTimelineEventNotFound
		}
		return fmt.Errorf("查询时间轴事件失败: %w", err)
	}

	// 删除时间轴事件
	if err := s.timelineEventRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除时间轴事件失败: %w", err)
	}
	return nil
}

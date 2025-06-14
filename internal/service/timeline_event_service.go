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
	ErrTimelineEventNotFound = errors.New("时间轴事件不存在")
)

// TimelineEventService 时间轴事件服务接口
type TimelineEventService interface {
	Service
	CreateTimelineEvent(ctx context.Context, createReq *dto.CreateTimelineEventRequest) (bool, error)
	GetTimelineEventByID(ctx context.Context, id int64) (*models.TimelineEvent, error)
	ListTimelineEventsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error)
	UpdateTimelineEvent(ctx context.Context, event *models.TimelineEvent, locationIDs, photoVideoIDs []int64) (*models.TimelineEvent, error)
	DeleteTimelineEvent(ctx context.Context, id int64) error
}

// timelineEventService 时间轴事件服务实现
type timelineEventService struct {
	*BaseService
	timelineEventRepo   repository.TimelineEventRepository
	locationRepo        repository.LocationRepository
	photoVideoRepo      repository.PhotoVideoRepository
	eventLocationRepo   repository.TimelineEventLocationRepository
	eventPhotoVideoRepo repository.TimelineEventPhotoVideoRepository
}

// NewTimelineEventService 创建时间轴事件服务
func NewTimelineEventService(
	timelineEventRepo repository.TimelineEventRepository,
	locationRepo repository.LocationRepository,
	photoVideoRepo repository.PhotoVideoRepository,
	eventLocationRepo repository.TimelineEventLocationRepository,
	eventPhotoVideoRepo repository.TimelineEventPhotoVideoRepository,
) TimelineEventService {
	return &timelineEventService{
		BaseService:         NewBaseService(timelineEventRepo),
		timelineEventRepo:   timelineEventRepo,
		locationRepo:        locationRepo,
		photoVideoRepo:      photoVideoRepo,
		eventLocationRepo:   eventLocationRepo,
		eventPhotoVideoRepo: eventPhotoVideoRepo,
	}
}

// CreateTimelineEvent 创建时间轴事件
func (s *timelineEventService) CreateTimelineEvent(ctx context.Context, createReq *dto.CreateTimelineEventRequest) (bool, error) {
	model, err := createReq.ToModel()
	if err != nil {
		return false, fmt.Errorf("换成实体对象失败：%w", err)
	}
	if err := s.timelineEventRepo.Create(ctx, model); err != nil {
		return false, fmt.Errorf("创建时间轴事件失败: %w", err)
	}

	if err := s.associateLocations(ctx, model.ID, createReq.LocationIDs); err != nil {
		return false, fmt.Errorf("关联地点失败: %w", err)
	}

	if err := s.associatePhotosVideos(ctx, model.ID, createReq.PhotoVideoIDs); err != nil {
		return false, fmt.Errorf("关联照片/视频失败: %w", err)
	}

	return true, nil
}

// GetTimelineEventByID 通过ID获取时间轴事件
func (s *timelineEventService) GetTimelineEventByID(ctx context.Context, id int64) (*models.TimelineEvent, error) {
	event, err := s.timelineEventRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTimelineEventNotFound) {
			return nil, ErrTimelineEventNotFound
		}
		return nil, fmt.Errorf("获取时间轴事件失败: %w", err)
	}

	if err := s.loadTimelineEventAssociations(ctx, event); err != nil {
		return nil, fmt.Errorf("加载时间轴事件关联数据失败: %w", err)
	}

	return event, nil
}

// ListTimelineEventsByCoupleID 获取情侣关系下的所有时间轴事件
func (s *timelineEventService) ListTimelineEventsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.TimelineEvent, int64, error) {
	return s.timelineEventRepo.FindByCoupleID(ctx, coupleID, offset, limit)
}

// UpdateTimelineEvent 更新时间轴事件
func (s *timelineEventService) UpdateTimelineEvent(ctx context.Context, event *models.TimelineEvent, locationIDs, photoVideoIDs []int64) (*models.TimelineEvent, error) {
	if err := s.timelineEventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("更新时间轴事件失败: %w", err)
	}

	if locationIDs != nil {
		if err := s.eventLocationRepo.DeleteByEventID(ctx, event.ID); err != nil {
			return nil, fmt.Errorf("删除现有关联失败: %w", err)
		}

		if err := s.associateLocations(ctx, event.ID, locationIDs); err != nil {
			return nil, fmt.Errorf("关联地点失败: %w", err)
		}
	}

	if photoVideoIDs != nil {
		if err := s.eventPhotoVideoRepo.DeleteByEventID(ctx, event.ID); err != nil {
			return nil, fmt.Errorf("删除现有关联失败: %w", err)
		}

		if err := s.associatePhotosVideos(ctx, event.ID, photoVideoIDs); err != nil {
			return nil, fmt.Errorf("关联照片/视频失败: %w", err)
		}
	}

	return s.GetTimelineEventByID(ctx, event.ID)
}

// DeleteTimelineEvent 删除时间轴事件
func (s *timelineEventService) DeleteTimelineEvent(ctx context.Context, id int64) error {
	if err := s.timelineEventRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除时间轴事件失败: %w", err)
	}

	if err := s.eventLocationRepo.DeleteByEventID(ctx, id); err != nil {
		return fmt.Errorf("删除地点关联失败: %w", err)
	}

	if err := s.eventPhotoVideoRepo.DeleteByEventID(ctx, id); err != nil {
		return fmt.Errorf("删除照片/视频关联失败: %w", err)
	}

	return nil
}

// 私有辅助方法

func (s *timelineEventService) loadTimelineEventAssociations(ctx context.Context, event *models.TimelineEvent) error {
	eventLocations, err := s.eventLocationRepo.FindByEventID(ctx, event.ID)
	if err != nil {
		return err
	}

	var locationIDs []int64
	for _, el := range eventLocations {
		locationIDs = append(locationIDs, el.LocationID)
	}

	if len(locationIDs) > 0 {
		locations, err := s.locationRepo.FindByIDs(ctx, locationIDs)
		if err != nil {
			return err
		}
		event.Locations = locations
	}

	eventPhotosVideos, err := s.eventPhotoVideoRepo.FindByEventID(ctx, event.ID)
	if err != nil {
		return err
	}

	var photoVideoIDs []int64
	for _, epv := range eventPhotosVideos {
		photoVideoIDs = append(photoVideoIDs, epv.PhotoVideoID)
	}

	if len(photoVideoIDs) > 0 {
		photosVideos, err := s.photoVideoRepo.FindByIDs(ctx, photoVideoIDs)
		if err != nil {
			return err
		}
		event.PhotosVideos = photosVideos
	}

	return nil
}

func (s *timelineEventService) associateLocations(ctx context.Context, eventID int64, locationIDs []int64) error {
	if len(locationIDs) == 0 {
		return nil
	}

	locations, err := s.locationRepo.FindByIDs(ctx, locationIDs)
	if err != nil {
		return err
	}

	if len(locations) != len(locationIDs) {
		return errors.New("some locations do not exist")
	}

	for _, locationID := range locationIDs {
		eventLocation := &models.TimelineEventLocation{
			TimelineEventID: eventID,
			LocationID:      locationID,
		}
		if err := s.eventLocationRepo.Create(ctx, eventLocation); err != nil {
			return err
		}
	}

	return nil
}

func (s *timelineEventService) associatePhotosVideos(ctx context.Context, eventID int64, photoVideoIDs []int64) error {
	if len(photoVideoIDs) == 0 {
		return nil
	}

	photosVideos, err := s.photoVideoRepo.FindByIDs(ctx, photoVideoIDs)
	if err != nil {
		return err
	}

	if len(photosVideos) != len(photoVideoIDs) {
		return errors.New("some photos/videos do not exist")
	}

	for _, photoVideoID := range photoVideoIDs {
		eventPhotoVideo := &models.TimelineEventPhotoVideo{
			TimelineEventID: eventID,
			PhotoVideoID:    photoVideoID,
		}
		if s.eventPhotoVideoRepo == nil {
			fmt.Println("eventPhotoVideoRepo is nil")
		}
		if err := s.eventPhotoVideoRepo.Create(ctx, eventPhotoVideo); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

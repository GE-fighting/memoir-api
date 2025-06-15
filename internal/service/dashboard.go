package service

import (
	"context"
	"memoir-api/internal/api/dto"
	"time"
)

type DashboardService interface {
	GetDashboardData(ctx context.Context, userId int64) (*dto.DashboardDTO, error)
}

type dashboardService struct {
	*BaseService
	userService          UserService
	albumService         CoupleAlbumService
	coupleService        CoupleService
	photoVideoService    PhotoVideoService
	timelineEventService TimelineEventService
	locationService      LocationService
}

func (d dashboardService) GetDashboardData(ctx context.Context, userId int64) (*dto.DashboardDTO, error) {
	// 获取情侣ID
	coupleID, err := d.userService.GetCoupleID(ctx, userId)
	if err != nil {
		return nil, err
	}
	// 创建仪表盘DTO
	dashboard := &dto.DashboardDTO{}

	// 获取故事(时间线事件)数量
	total, err := d.timelineEventService.CountByCoupleID(ctx, coupleID)
	if err != nil {
		return nil, err
	}
	dashboard.StoryCount = int(total)

	// 获取媒体数量
	total, err = d.photoVideoService.CountByCoupleID(ctx, coupleID)
	if err != nil {
		return nil, err
	}
	dashboard.MediaCount = int(total)

	// 获取相册数量
	total, err = d.albumService.CountByCoupleID(ctx, coupleID)
	if err != nil {
		return nil, err
	}
	dashboard.AlbumCount = int(total)

	// 获取情侣天数
	couple, err := d.coupleService.GetCoupleByID(ctx, coupleID)
	if err != nil {
		return nil, err
	}

	// 计算情侣天数
	// AnniversaryDate是time.Time类型，不是指针
	if !couple.AnniversaryDate.IsZero() {
		days := int(time.Since(couple.AnniversaryDate).Hours() / 24)
		if days < 0 {
			days = 0
		}
		dashboard.CoupleDays = days
	}

	locations, _, err := d.locationService.ListLocationsByCoupleID(ctx, coupleID, -1, -1)
	if err != nil {
		return nil, err
	}
	dashboard.Locations = locations
	return dashboard, nil
}

func NewDashboardService(service UserService,
	albumService CoupleAlbumService, coupleService CoupleService,
	photoVideoService PhotoVideoService,
	timelineEventService TimelineEventService, locationService LocationService) DashboardService {
	return &dashboardService{
		BaseService:          NewBaseService(nil),
		userService:          service,
		albumService:         albumService,
		coupleService:        coupleService,
		photoVideoService:    photoVideoService,
		timelineEventService: timelineEventService,
		locationService:      locationService,
	}
}

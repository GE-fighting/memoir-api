package service

import (
	"memoir-api/internal/repository"
)

// Factory 服务工厂接口
type Factory interface {
	User() UserService
	Couple() CoupleService
	JWT() JWTService
	Location() LocationService
	TimelineEvent() TimelineEventService
	PhotoVideo() PhotoVideoService
	Wishlist() WishlistService
	PersonalMedia() PersonalMediaService
	CoupleAlbum() CoupleAlbumService
	Dashboard() DashboardService
	Attachment() AttachmentService
}

// factory 服务工厂实现
type factory struct {
	userService          UserService
	coupleService        CoupleService
	jwtService           JWTService
	locationService      LocationService
	timelineEventService TimelineEventService
	photoVideoService    PhotoVideoService
	wishlistService      WishlistService
	personalMediaService PersonalMediaService
	coupleAlbumService   CoupleAlbumService
	dashboardService     DashboardService
	attachmentService    AttachmentService
}

// NewFactory 创建服务工厂
func NewFactory(repoFactory repository.Factory) Factory {
	jwtService := NewJWTService()
	userService := NewUserService(repoFactory.User(), repoFactory.Couple())
	coupleService := NewCoupleService(
		repoFactory.Couple(),
		repoFactory.User(),
	)
	locationService := NewLocationService(repoFactory.Location())
	timelineEventService := NewTimelineEventService(repoFactory.TimelineEvent(), repoFactory.Location(), repoFactory.PhotoVideo(), repoFactory.TimelineEventLocation(), repoFactory.TimelineEventPhotoVideo())
	photoVideoService := NewPhotoVideoService(repoFactory.PhotoVideo(), repoFactory.User(), repoFactory.CoupleAlbum())
	wishlistService := NewWishlistService(
		repoFactory.Wishlist(),
		repoFactory.WishlistAttachment(),
		repoFactory.Attachment(),
	)
	personalMediaService := NewPersonalMediaService(repoFactory.PersonalMedia())
	coupleAlbumService := NewCoupleAlbumService(repoFactory.CoupleAlbum(), userService, photoVideoService)
	dashboardService := NewDashboardService(userService, coupleAlbumService, coupleService, photoVideoService, timelineEventService, locationService)
	attachmentService := NewAttachmentService(repoFactory.Attachment(), repoFactory.User(), repoFactory.Couple())

	return &factory{
		userService:          userService,
		coupleService:        coupleService,
		jwtService:           jwtService,
		locationService:      locationService,
		timelineEventService: timelineEventService,
		photoVideoService:    photoVideoService,
		wishlistService:      wishlistService,
		personalMediaService: personalMediaService,
		coupleAlbumService:   coupleAlbumService,
		dashboardService:     dashboardService,
		attachmentService:    attachmentService,
	}
}

// User 获取用户服务
func (f *factory) User() UserService {
	return f.userService
}

// Couple 获取情侣关系服务
func (f *factory) Couple() CoupleService {
	return f.coupleService
}

// JWT 获取JWT服务
func (f *factory) JWT() JWTService {
	return f.jwtService
}

// Location 获取地点服务
func (f *factory) Location() LocationService {
	return f.locationService
}

// TimelineEvent 获取时间轴事件服务
func (f *factory) TimelineEvent() TimelineEventService {
	return f.timelineEventService
}

// PhotoVideo 获取照片和视频服务
func (f *factory) PhotoVideo() PhotoVideoService {
	return f.photoVideoService
}

// Wishlist 获取心愿清单服务
func (f *factory) Wishlist() WishlistService {
	return f.wishlistService
}

// PersonalMedia 获取个人媒体服务
func (f *factory) PersonalMedia() PersonalMediaService {
	return f.personalMediaService
}

// CoupleAlbum 获取情侣相册服务
func (f *factory) CoupleAlbum() CoupleAlbumService {
	return f.coupleAlbumService
}

// Dashboard 获取仪表盘服务
func (f *factory) Dashboard() DashboardService { return f.dashboardService }

// Attachment 获取附件服务
func (f *factory) Attachment() AttachmentService {
	return f.attachmentService
}

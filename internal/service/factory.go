package service

import (
	"memoir-api/internal/cache"
	"memoir-api/internal/config"
	"memoir-api/internal/email"
	"memoir-api/internal/logger"
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
	Email() EmailService
	CoupleReminder() CoupleReminderService
}

// factory 服务工厂实现
type factory struct {
	userService           UserService
	coupleService         CoupleService
	jwtService            JWTService
	locationService       LocationService
	timelineEventService  TimelineEventService
	photoVideoService     PhotoVideoService
	wishlistService       WishlistService
	personalMediaService  PersonalMediaService
	coupleAlbumService    CoupleAlbumService
	dashboardService      DashboardService
	attachmentService     AttachmentService
	emailService          EmailService
	coupleReminderService CoupleReminderService
}

// NewFactory 创建服务工厂
func NewFactory(repoFactory repository.Factory) Factory {
	// 创建服务实例
	userRepo := repoFactory.User()
	coupleRepo := repoFactory.Couple()

	// 获取Redis客户端
	redisClient := cache.GetRedisClient().Client

	// 创建邮件服务
	cfg := config.New()
	emailService, err := email.NewEmailService(cfg, redisClient)
	if err != nil {
		logger.Fatal(err, "Failed to create email service")
	}

	// 创建用户服务
	userService := NewUserService(userRepo, coupleRepo, emailService)

	// 创建JWT服务
	jwtService := NewJWTService()

	// 创建情侣服务
	coupleService := NewCoupleService(coupleRepo, userRepo)

	// 创建位置服务
	locationService := NewLocationService(repoFactory.Location())

	// 创建时间线事件服务
	timelineEventService := NewTimelineEventService(
		repoFactory.TimelineEvent(),
		repoFactory.Location(),
		repoFactory.PhotoVideo(),
		repoFactory.TimelineEventLocation(),
		repoFactory.TimelineEventPhotoVideo(),
	)

	// 创建照片视频服务
	photoVideoService := NewPhotoVideoService(
		repoFactory.PhotoVideo(),
		userRepo,
		repoFactory.CoupleAlbum(),
	)

	// 创建心愿单服务
	wishlistService := NewWishlistService(
		repoFactory.Wishlist(),
		repoFactory.WishlistAttachment(),
		repoFactory.Attachment(),
	)

	// 创建个人媒体服务
	personalMediaService := NewPersonalMediaService(repoFactory.PersonalMedia())

	// 创建情侣相册服务
	coupleAlbumService := NewCoupleAlbumService(
		repoFactory.CoupleAlbum(),
		userService,
		photoVideoService,
	)

	// 创建附件服务
	attachmentService := NewAttachmentService(
		repoFactory.Attachment(),
		userRepo,
		coupleRepo,
	)

	// 创建仪表盘服务
	dashboardService := NewDashboardService(
		userService,
		coupleAlbumService,
		coupleService,
		photoVideoService,
		timelineEventService,
		locationService,
	)

	// 创建情侣纪念日提醒服务
	coupleReminderService := NewCoupleReminderService(
		coupleRepo,
		userRepo,
		emailService,
	)

	return &factory{
		userService:           userService,
		coupleService:         coupleService,
		jwtService:            jwtService,
		locationService:       locationService,
		timelineEventService:  timelineEventService,
		photoVideoService:     photoVideoService,
		wishlistService:       wishlistService,
		personalMediaService:  personalMediaService,
		coupleAlbumService:    coupleAlbumService,
		dashboardService:      dashboardService,
		attachmentService:     attachmentService,
		emailService:          emailService,
		coupleReminderService: coupleReminderService,
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

func (f *factory) Email() EmailService {
	return f.emailService
}

func (f *factory) CoupleReminder() CoupleReminderService {
	return f.coupleReminderService
}

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
	timelineEventService := NewTimelineEventService(repoFactory.TimelineEvent())
	photoVideoService := NewPhotoVideoService(repoFactory.PhotoVideo())
	wishlistService := NewWishlistService(repoFactory.Wishlist())

	return &factory{
		userService:          userService,
		coupleService:        coupleService,
		jwtService:           jwtService,
		locationService:      locationService,
		timelineEventService: timelineEventService,
		photoVideoService:    photoVideoService,
		wishlistService:      wishlistService,
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
